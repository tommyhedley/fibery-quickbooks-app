package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"golang.org/x/sync/singleflight"
)

type RequestType struct {
	SourceId  string
	Sync      fibery.SyncType
	GroupSize int
	Done      bool
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
	Group       singleflight.Group
	GroupCounts map[DataKey]int
	Types       map[string]RequestType
	Account     QuickBooksAccountInfo
	DataCache   map[DataKey]*CacheEntry
	IdCache     *IdCache
	LastSynced  time.Time
	lastReturn  time.Time
	ctx         context.Context
	cancel      context.CancelFunc
}

type OperationManager struct {
	sync.Mutex
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

func (om *OperationManager) GetOrAddOperation(operationId string, lastSynced time.Time, requestedTypes []string, acct QuickBooksAccountInfo, tr TypeRegistry, s *IdStore) (*Operation, error) {
	om.Lock()
	op, exists := om.Operations[operationId]
	om.Unlock()
	if !exists {
		types := make(map[string]RequestType, len(requestedTypes))
		idCache, cacheExists := s.GetOrCreateIdCache(acct.RealmID)
		slog.Debug(fmt.Sprintf("an existing cache was available: %t, for realmId: %s\n", cacheExists, acct.RealmID))
		slog.Debug(fmt.Sprintf("lastSynced is not zero: %t\n", !lastSynced.IsZero()))
		groupCounts := make(map[string]int)
		for _, requestedType := range requestedTypes {
			storedType, exists := tr.GetType(requestedType)
			if !exists {
				return nil, fmt.Errorf("requested type: %s does not exist in the TypeRegistry", requestedType)
			}
			src := storedType.SourceId() // call to the SourceId method
			groupCounts[src]++
		}

		for _, requestedType := range requestedTypes {
			storedType, _ := tr.GetType(requestedType)
			groupSize := groupCounts[storedType.SourceId()]
			if !lastSynced.IsZero() && cacheExists {
				var cdc bool
				switch storedType.(type) {
				case CDCType:
					slog.Debug(fmt.Sprintf("%s is cdc\n", storedType.Id()))
					cdc = true
				case CDCDepType:
					slog.Debug(fmt.Sprintf("%s is cdc\n", storedType.Id()))
					cdc = true
				default:
					slog.Debug(fmt.Sprintf("%s is not cdc\n", storedType.Id()))
					cdc = false
				}
				if cdc {
					types[requestedType] = RequestType{
						SourceId:  storedType.SourceId(),
						Sync:      fibery.Delta,
						GroupSize: groupSize,
						Done:      false,
					}
				} else {
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

		op = &Operation{
			Types:      types,
			Account:    acct,
			DataCache:  make(map[DataKey]*CacheEntry),
			IdCache:    idCache,
			LastSynced: lastSynced,
			lastReturn: time.Now(),
			ctx:        ctx,
			cancel:     cancel,
		}

		om.Lock()
		defer om.Unlock()
		om.Operations[operationId] = op
	}

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
	op.Lock()
	if entry, ok := op.DataCache[key]; ok {
		res := entry.Value
		entry.Pending--
		if entry.Pending == 0 {
			delete(op.DataCache, key)
		}
		op.Unlock()
		return res, nil
	}
	op.Unlock()

	keyStr := fmt.Sprintf("%s-%d", key.DataType, key.StartPosition)

	op.Lock()
	if op.GroupCounts == nil {
		op.GroupCounts = make(map[DataKey]int)
	}
	op.GroupCounts[key]++
	op.Unlock()

	res, err, _ := op.Group.Do(keyStr, func() (interface{}, error) {
		return fetchFn()
	})
	if err != nil {
		op.Lock()
		op.GroupCounts[key]--
		if op.GroupCounts[key] <= 0 {
			delete(op.GroupCounts, key)
		}
		op.Unlock()
		return nil, err
	}

	op.Lock()
	count := op.GroupCounts[key]
	delete(op.GroupCounts, key)
	op.Unlock()

	if count < expected {
		op.Lock()
		op.DataCache[key] = &CacheEntry{
			Value:   res,
			Pending: expected - count,
		}
		op.Unlock()
	}
	return res, nil
}
