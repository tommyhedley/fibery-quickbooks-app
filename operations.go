package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
	"golang.org/x/sync/singleflight"
)

type RequestType struct {
	SourceId    string
	Sync        fibery.SyncType
	GroupSize   int
	Done        bool
	Attachables map[string][]string
}

type DataKey struct {
	DataType      string
	StartPosition int
}

type CacheEntry struct {
	Value   any
	Pending int
}

type Operation struct {
	sync.Mutex
	Group      singleflight.Group
	Types      map[string]RequestType
	Account    QuickBooksAccountInfo
	DataCache  map[DataKey]*CacheEntry
	IdCache    *IdCache
	LastSynced time.Time
	lastReturn time.Time
	ctx        context.Context
	cancel     context.CancelFunc
}

type OperationManager struct {
	sync.Mutex
	Group      singleflight.Group
	Operations map[string]*Operation
	ttl        time.Duration
}

func NewOperationManager(ttl time.Duration) *OperationManager {
	return &OperationManager{
		Operations: make(map[string]*Operation),
		ttl:        ttl,
	}
}

func (om *OperationManager) DeleteOperation(operationId string) {
	om.Lock()
	defer om.Unlock()
	if op, exists := om.Operations[operationId]; exists {
		if op.cancel != nil {
			op.cancel()
		}
	}
	delete(om.Operations, operationId)
}

func (om *OperationManager) CleanupExpired() {
	om.Lock()
	defer om.Unlock()
	now := time.Now()
	for opId, op := range om.Operations {
		if !op.lastReturn.IsZero() && now.Sub(op.lastReturn) > om.ttl {
			if op.cancel != nil {
				op.cancel()
			}
			delete(om.Operations, opId)
		}
	}
}

func buildOperation(req SyncDataRequest, tr TypeRegistry, s *IdStore, client *quickbooks.Client, pageSize int) (*Operation, error) {
	var op *Operation
	types := make(map[string]RequestType, len(req.Types))
	idCache, cacheExists := s.GetOrCreateIdCache(req.Account.RealmID)
	slog.Debug(fmt.Sprintf("a cache for realmId: %s exists: %t", req.Account.RealmID, cacheExists))
	slog.Debug(fmt.Sprintf("lastSynced is zero: %t", req.LastSyncronizedAt.IsZero()))
	groupCounts := make(map[string]int)
	for _, requestedType := range req.Types {
		storedType, exists := tr.GetType(requestedType)
		if !exists {
			return nil, fmt.Errorf("requested type: %s does not exist in the TypeRegistry", requestedType)
		}
		src := storedType.SourceId() // call to the SourceId method
		groupCounts[src]++
	}

	cdc := !req.LastSyncronizedAt.IsZero() && cacheExists

	for _, requestedType := range req.Types {
		storedType, _ := tr.GetType(requestedType)
		groupSize := groupCounts[storedType.SourceId()]
		if cdc {
			switch storedType.(type) {
			case CDCType:
				slog.Debug(fmt.Sprintf("%s is cdc\n", storedType.Id()))
				types[requestedType] = RequestType{
					SourceId:  storedType.SourceId(),
					Sync:      fibery.Delta,
					GroupSize: groupSize,
					Done:      false,
				}
			case CDCDepType:
				slog.Debug(fmt.Sprintf("%s is cdc\n", storedType.Id()))
				types[requestedType] = RequestType{
					SourceId:  storedType.SourceId(),
					Sync:      fibery.Delta,
					GroupSize: groupSize,
					Done:      false,
				}
			default:
				slog.Debug(fmt.Sprintf("%s is not cdc\n", storedType.Id()))
				types[requestedType] = RequestType{
					SourceId:  storedType.SourceId(),
					Sync:      fibery.Full,
					GroupSize: groupSize,
					Done:      false,
				}
			}
		} else {
			types[requestedType] = RequestType{
				SourceId:  storedType.SourceId(),
				Sync:      fibery.Full,
				GroupSize: groupSize,
				Done:      false,
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		if op == nil {
			cancel()
		}
	}()

	attachmentSources := GetAttachmentSources(req.Schema, "Attachables")
	if len(attachmentSources) > 0 {
		requestParams := quickbooks.RequestParameters{
			Ctx:             ctx,
			RealmId:         req.Account.RealmID,
			Token:           &req.Account.BearerToken,
			WaitOnRateLimit: true,
		}
		if cdc {
			cdc, err := client.ChangeDataCapture(requestParams, []string{"Attachable"}, req.LastSyncronizedAt)
			if err != nil {
				return nil, fmt.Errorf("unable to make ChangeDataCapture for Attachable entities: %w", err)
			}
			slog.Debug("attachables cdc request completed")

			attachables := quickbooks.CDCQueryExtractor[quickbooks.Attachable](&cdc)

			for typeId, reqType := range types {
				if attachmentSources[typeId] {
					if reqType.Attachables == nil {
						reqType.Attachables = make(map[string][]string)
						types[typeId] = reqType
					}
				}
			}

			for _, attachable := range attachables {
				if len(attachable.AttachableRef) == 0 {
					continue
				}

				attachableURL := GenerateAttachablesURL(attachable)

				for _, ref := range attachable.AttachableRef {
					if !attachmentSources[ref.EntityRef.Type] {
						continue
					}

					for typeId, reqType := range types {
						if typeId == ref.EntityRef.Type {
							if reqType.Attachables == nil {
								reqType.Attachables = make(map[string][]string)
							}

							reqType.Attachables[ref.EntityRef.Value] = append(reqType.Attachables[ref.EntityRef.Value], attachableURL)

							types[typeId] = reqType
						}
					}
				}
			}
		} else {
			pendingSources := make(map[string]int, len(attachmentSources))
			for source := range attachmentSources {
				pendingSources[source] = 1
			}

			for len(pendingSources) > 0 {
				batchRequest := make([]quickbooks.BatchItemRequest, 0, len(pendingSources))
				for source, startPos := range pendingSources {
					req := quickbooks.BatchItemRequest{
						BID:   source,
						Query: fmt.Sprintf("Select Id, AttachableRef From Attachable Where AttachableRef.EntityRef.Type = '%s' STARTPOSITION %d MAXRESULTS %d", source, startPos, pageSize),
					}
					batchRequest = append(batchRequest, req)
				}

				batchResponse, err := client.BatchRequest(requestParams, batchRequest)
				if err != nil {
					return nil, fmt.Errorf("unable to make batch request: %w", err)
				}
				slog.Debug("attachables batch request completed")

				for _, itemResponse := range batchResponse {
					if faultType := itemResponse.Fault.FaultType; faultType != "" {
						return nil, fmt.Errorf("batch request error: %s", faultType)
					}

					attachables := quickbooks.BatchQueryExtractor[quickbooks.Attachable](&itemResponse)

					if len(attachables) == pageSize {
						pendingSources[itemResponse.BID] += pageSize
					} else {
						delete(pendingSources, itemResponse.BID)
					}

					for _, attachable := range attachables {
						if len(attachable.AttachableRef) == 0 {
							continue
						}

						attachableString := fmt.Sprintf("app://resource?type=%s&id=%s", "attachable", attachable.Id)

						for _, ref := range attachable.AttachableRef {
							if itemResponse.BID != ref.EntityRef.Type {
								continue
							}

							reqType := types[itemResponse.BID]
							if reqType.Attachables == nil {
								reqType.Attachables = make(map[string][]string)
							}

							reqType.Attachables[ref.EntityRef.Value] = append(reqType.Attachables[ref.EntityRef.Value], attachableString)

							types[itemResponse.BID] = reqType

						}
					}
				}
			}
		}
	}

	op = &Operation{
		Types:      types,
		Account:    req.Account,
		DataCache:  make(map[DataKey]*CacheEntry),
		IdCache:    idCache,
		LastSynced: req.LastSyncronizedAt,
		lastReturn: time.Now(),
		ctx:        ctx,
		cancel:     cancel,
	}
	return op, nil
}

func (om *OperationManager) GetOrAddOperation(req SyncDataRequest, tr TypeRegistry, s *IdStore, client *quickbooks.Client, pageSize int) (*Operation, error) {
	om.Lock()
	if op, exists := om.Operations[req.OperationID]; exists {
		slog.Debug(fmt.Sprintf("requested operation: %s exists: %t", req.OperationID, exists))
		om.Unlock()
		return op, nil
	}
	om.Unlock()

	result, err, _ := om.Group.Do(req.OperationID, func() (interface{}, error) {
		return buildOperation(req, tr, s, client, pageSize)
	})
	if err != nil {
		return nil, fmt.Errorf("unable to build operation: %w", err)
	}

	op := result.(*Operation)

	om.Lock()
	om.Operations[req.OperationID] = op
	om.Unlock()

	return op, nil
}

func (op *Operation) RefreshLastReturn() {
	op.Lock()
	op.lastReturn = time.Now()
	op.Unlock()
}

func (op *Operation) MarkTypeFulfilled(requestedType string) {
	op.Lock()
	if reqType, ok := op.Types[requestedType]; ok {
		reqType.Done = true
		op.Types[requestedType] = reqType
	}
	op.Unlock()

}

func (op *Operation) IsComplete() bool {
	op.Lock()
	defer op.Unlock()
	for _, rt := range op.Types {
		if !rt.Done {
			return false
		}
	}
	return true
}

func (op *Operation) GetOrFetchData(key DataKey, expected int, fetchFn func() (any, error)) (any, error) {
	if expected == 1 {
		slog.Debug(fmt.Sprintf("%s: only 1 fulfillment expected, no caching necessary", key.DataType))
		return fetchFn()
	}

	op.Lock()
	if entry, ok := op.DataCache[key]; ok {
		slog.Debug(fmt.Sprintf("datacache for key: %s:%d exists", key.DataType, key.StartPosition))
		res := entry.Value
		entry.Pending--
		slog.Debug(fmt.Sprintf("%s: %d datatypes remaining to be fulfilled", key.DataType, entry.Pending))
		if entry.Pending == 0 {
			delete(op.DataCache, key)
			slog.Debug(fmt.Sprintf("datacache for key: %s:%d was deleted", key.DataType, key.StartPosition))
		}
		op.Unlock()
		return res, nil
	}
	op.Unlock()

	slog.Debug(fmt.Sprintf("datacache for key: %s:%d does not exist", key.DataType, key.StartPosition))

	keyStr := fmt.Sprintf("%s:%d", key.DataType, key.StartPosition)

	res, err, _ := op.Group.Do(keyStr, func() (interface{}, error) {
		v, err := fetchFn()
		if err != nil {
			return nil, err
		}
		op.Lock()
		op.DataCache[key] = &CacheEntry{Value: v, Pending: expected}
		slog.Debug(fmt.Sprintf("%s cache stored (pending=%d)", key.DataType, expected))
		op.Unlock()
		return v, nil
	})
	if err != nil {
		return nil, err
	}

	op.Lock()
	e := op.DataCache[key]
	e.Pending--
	slog.Debug(fmt.Sprintf("%s: %d remaining", key.DataType, e.Pending))
	if e.Pending == 0 {
		delete(op.DataCache, key)
		slog.Debug(fmt.Sprintf("%s cache deleted", key.DataType))
	}
	op.Unlock()
	return res, nil
}
