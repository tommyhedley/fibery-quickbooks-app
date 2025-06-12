package integration

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type SyncRequest struct {
	RequestedType     string                             `json:"requestedType"`
	OperationId       string                             `json:"operationId"`
	Types             []string                           `json:"types"`
	Schema            map[string]map[string]fibery.Field `json:"schema"`
	Filter            map[string]any                     `json:"filter"`
	Account           QuickBooksAccountInfo              `json:"account"`
	LastSyncronizedAt time.Time                          `json:"lastSynchronizedAt"`
	Pagination        fibery.NextPageConfig              `json:"pagination"`
}

type RequestType int

const (
	Normal RequestType = iota
	ChangeDataCapture
)

type OperationDataHandlerResponse struct {
	Error error
	fibery.DataHandlerResponse
}

type SourceGroup struct {
	getAttachable bool
	expectedUses  int
	request       RequestType
	batchPages    map[int]*quickbooks.BatchItemResponse
	attachables   map[string][]quickbooks.Attachable
}

type Operation struct {
	sync.Mutex
	id                string
	wg                sync.WaitGroup
	startOnce         sync.Once
	cleanupOnce       sync.Once
	integration       *Integration
	existingCache     bool
	idCache           *IdCache
	account           QuickBooksAccountInfo
	lastSynced        time.Time
	lastRequest       time.Time
	requestTypes      map[string]fibery.Type
	sourceGroups      map[string]*SourceGroup
	changeDataCapture *quickbooks.ChangeDataCapture
	chans             map[string]chan OperationDataHandlerResponse
	ctx               context.Context
	cancel            context.CancelFunc
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

func (om *OperationManager) GetOrAddOperation(req SyncRequest, i *Integration) (*Operation, error) {
	om.Lock()
	defer om.Unlock()
	if op, exists := om.Operations[req.OperationId]; exists {
		slog.Debug(fmt.Sprintf("requested operation: %s exists: %t", req.OperationId, exists))
		return op, nil
	}

	op := buildOperation(req, i)

	om.Operations[req.OperationId] = op

	return op, nil
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
		if !op.lastRequest.IsZero() && now.Sub(op.lastRequest) > om.ttl {
			if op.cancel != nil {
				op.cancel()
			}
			delete(om.Operations, opId)
		}
	}
}

func buildOperation(req SyncRequest, i *Integration) *Operation {
	idCache, cacheExisted := i.idStore.GetOrCreateIdCache(req.Account.RealmId)
	ctx, cancel := context.WithCancel(i.ctx)
	op := &Operation{
		id:            req.OperationId,
		integration:   i,
		existingCache: cacheExisted,
		idCache:       idCache,
		account:       req.Account,
		requestTypes:  make(map[string]fibery.Type, len(req.Types)),
		sourceGroups:  make(map[string]*SourceGroup, len(req.Types)),
		chans:         make(map[string]chan OperationDataHandlerResponse, len(req.Types)),
		ctx:           ctx,
		cancel:        cancel,
	}
	op.wg.Add(len(req.Types))
	return op
}

func (op *Operation) propagateError(err error) {
	op.cancel()

	op.cleanupOnce.Do(func() {
		for key, ch := range op.chans {
			ch <- OperationDataHandlerResponse{
				Error: err,
			}
			delete(op.chans, key)
		}
	})
}

func (op *Operation) addSourceGroup(
	source string,
	reqType RequestType,
	getAttachable bool,
) {
	grp, ok := op.sourceGroups[source]
	if !ok {
		grp = &SourceGroup{
			expectedUses:  0,
			request:       reqType,
			getAttachable: getAttachable,
			batchPages:    make(map[int]*quickbooks.BatchItemResponse),
			attachables:   make(map[string][]quickbooks.Attachable),
		}
	}

	grp.expectedUses++

	if getAttachable {
		grp.getAttachable = true
	}

	if reqType == Normal || grp.request == Normal {
		grp.request = Normal
	}

	op.sourceGroups[source] = grp
}

func (op *Operation) processTypeEntry(
	regType fibery.Type,
	req SyncRequest,
	getAttach bool,
) error {
	var (
		source  string
		reqMode RequestType
	)

	switch t := regType.(type) {
	case CDCType:
		source = t.Type()
		if op.existingCache && !req.LastSyncronizedAt.IsZero() {
			reqMode = ChangeDataCapture
		} else {
			reqMode = Normal
		}

	case StandardType:
		source = t.Type()
		reqMode = Normal

	case CDCDependentType:
		source = t.SourceType()
		if op.existingCache && !req.LastSyncronizedAt.IsZero() {
			reqMode = ChangeDataCapture
		} else {
			reqMode = Normal
		}

	case StandardDependentType:
		source = t.SourceType()
		reqMode = Normal

	default:
		return fmt.Errorf(
			"registered type %s not a supported interface",
			regType.Id(),
		)
	}

	op.addSourceGroup(source, reqMode, getAttach)
	return nil
}

func ResponseChannelKey(id string, page int) string {
	return fmt.Sprintf("%s:%d", id, page)
}

func (op *Operation) SubmitRequest(req SyncRequest) error {
	op.Lock()
	op.lastRequest = time.Now()
	op.Unlock()

	if req.Pagination.Page == 0 || req.Pagination.Page == 1 {
		regType, exists := op.integration.types.Get(req.RequestedType)
		if !exists {
			return fmt.Errorf("requestedType: %s not found", req.RequestedType)
		}

		op.Lock()

		op.requestTypes[req.RequestedType] = regType

		channelKey := ResponseChannelKey(req.RequestedType, 1)
		if _, ok := op.chans[channelKey]; !ok {
			op.chans[channelKey] = make(chan OperationDataHandlerResponse, 1)
		}

		attachableFieldId := op.integration.config.AttachableFieldId

		switch t := regType.(type) {
		case StaticType:
			resp := OperationDataHandlerResponse{
				DataHandlerResponse: fibery.DataHandlerResponse{
					Items:               t.GetData(),
					SynchronizationType: fibery.Full,
				},
			}

			respChan := op.chans[channelKey]

			respChan <- resp

			op.wg.Done()
			return nil

		case UnionType:
			schema, ok := req.Schema[req.RequestedType]
			if !ok {
				return fmt.Errorf("no schema for %s was provided on the request", req.RequestedType)
			}

			getAttach := attachableField(schema, attachableFieldId)

			for _, sourceTypeId := range t.Types() {
				inner, ok := op.integration.types.Get(sourceTypeId)
				if !ok {
					return fmt.Errorf("inner type %s not found", sourceTypeId)
				}
				if err := op.processTypeEntry(inner, req, getAttach); err != nil {
					return err
				}
			}

		default:
			schema, ok := req.Schema[req.RequestedType]
			if !ok {
				return fmt.Errorf("no schema for %s was provided on the request", req.RequestedType)
			}

			getAttach := attachableField(schema, attachableFieldId)

			if err := op.processTypeEntry(regType, req, getAttach); err != nil {
				return err
			}
		}

		if op.lastSynced.After(req.LastSyncronizedAt) {
			op.lastSynced = req.LastSyncronizedAt
		}

		op.Unlock()
		op.wg.Done()
		op.startOnce.Do(func() {
			go func() {
				op.wg.Wait()
				op.fetchAll()
			}()
		})
	}
	return nil
}

func (op *Operation) GetChannel(key string) (<-chan OperationDataHandlerResponse, error) {
	op.Lock()
	defer op.Unlock()

	ch, ok := op.chans[key]
	if !ok {
		return nil, fmt.Errorf("no channel found for key: %s", key)
	}

	return ch, nil
}

func (op *Operation) TryCleanup() {
	op.Lock()
	shouldDelete := len(op.requestTypes) == 0
	opId := op.id
	op.Unlock()

	if shouldDelete {
		op.integration.opManager.DeleteOperation(opId)
	}
}

func (op *Operation) indexBatchQuery(batch []quickbooks.BatchItemResponse) (map[string]struct{}, error) {
	moreAttachables := map[string]struct{}{}
	for _, resp := range batch {
		faults := resp.Fault.Faults
		if len(faults) > 0 {
			return nil, fmt.Errorf("fault for %s: %w", resp.BID, quickbooks.BatchError{Faults: faults})
		}

		entityType, startPosition, isAttachable, err := DecodeQueryBID(resp.BID)
		if err != nil {
			return nil, fmt.Errorf("error decoding query BID: %w", err)
		}

		op.Lock()
		defer op.Unlock()

		sourceGroup, exists := op.sourceGroups[entityType]
		if !exists {
			return nil, fmt.Errorf("no sourceGroup found for: %s", entityType)
		}

		if isAttachable {
			attachables := resp.QueryResponse.Attachable
			if len(attachables) > 0 {
				existing := sourceGroup.attachables
				updated, more := indexAttachables(
					entityType, attachables, existing, op.integration.config.QuickBooks.PageSize,
				)
				sourceGroup.attachables = updated
				if more {
					moreAttachables[entityType] = struct{}{}
				}
			}
			continue
		} else {
			if _, exists := sourceGroup.batchPages[startPosition]; exists {
				return nil, fmt.Errorf("a batch response entry already exists for %s:%d", entityType, startPosition)
			}

			sourceGroup.batchPages[startPosition] = &resp
		}
	}

	return moreAttachables, nil
}

func (op *Operation) doBatch(req []quickbooks.BatchItemRequest, params quickbooks.RequestParameters) {
	client := op.integration.client

	page := 1

	for {

		batch, err := client.BatchRequest(params, req)
		if err != nil {
			op.propagateError(fmt.Errorf("error fetching inital batch: %w", err))
		}

		nextAttachEntities, err := op.indexBatchQuery(batch)
		if err != nil {
			op.propagateError(fmt.Errorf("error indexing inital batch: %w", err))
		}

		if len(nextAttachEntities) > 0 {
			page++
			pageSize := op.integration.config.QuickBooks.PageSize
			nextReq := make([]quickbooks.BatchItemRequest, 0, len(nextAttachEntities))
			for entityType := range nextAttachEntities {
				r := batchQueryRequest(entityType, nil, page, pageSize, true)
				nextReq = append(nextReq, r)
			}
			req = nextReq
		} else {
			break
		}
	}

}

func (op *Operation) doCDC(req []string, params quickbooks.RequestParameters) {
	client := op.integration.client

	cdc, err := client.ChangeDataCapture(params, req, op.lastSynced)
	if err != nil {
		op.propagateError(fmt.Errorf("error fetching cdc: %w", err))
	}

	op.Lock()
	defer op.Unlock()

	op.changeDataCapture = &cdc
}

func (op *Operation) fetchAll() {
	var (
		initalFetch sync.WaitGroup
		batchReq    []quickbooks.BatchItemRequest
		cdcReq      []string
	)

	page := 1
	pageSize := op.integration.config.QuickBooks.PageSize

	for sourceType, group := range op.sourceGroups {
		if group.getAttachable {
			req := batchQueryRequest(sourceType, nil, page, pageSize, true)
			batchReq = append(batchReq, req)
		}
		switch group.request {
		case ChangeDataCapture:
			cdcReq = append(cdcReq, sourceType)
		case Normal:
			req := batchQueryRequest(sourceType, nil, page, pageSize, false)
			batchReq = append(batchReq, req)
		}
	}

	params := quickbooks.RequestParameters{
		Ctx:             op.ctx,
		RealmId:         op.account.RealmId,
		Token:           &op.account.BearerToken,
		WaitOnRateLimit: true,
	}

	if len(cdcReq) > 0 {
		initalFetch.Add(1)
		go func(req []string) {
			defer initalFetch.Done()
			op.doCDC(req, params)
		}(cdcReq)
	}

	if len(batchReq) > 0 {
		initalFetch.Add(1)
		go func(req []quickbooks.BatchItemRequest) {
			defer initalFetch.Done()
			op.doBatch(req, params)
		}(batchReq)
	}

	initalFetch.Wait()
	op.dispatchPages()
}

func (op *Operation) dispatchPages() {
	params := quickbooks.RequestParameters{
		Ctx:             op.ctx,
		RealmId:         op.account.RealmId,
		Token:           &op.account.BearerToken,
		WaitOnRateLimit: true,
	}
	page := 1
	pageSize := op.integration.config.QuickBooks.PageSize

	for len(op.requestTypes) > 0 {
		nextBatchSources := map[string]struct{}{}

		for typeId, regType := range op.requestTypes {
			key := ResponseChannelKey(typeId, page)
			ch, ok := op.chans[key]
			if !ok {
				op.propagateError(fmt.Errorf("missing channel %s", key))
				delete(op.requestTypes, typeId)
				continue
			}

			var src string
			switch t := regType.(type) {
			case CDCType:
				src = t.Type()
			case CDCDependentType:
				src = t.SourceType()
			case StandardType:
				src = t.Type()
			case StandardDependentType:
				src = t.SourceType()
			case UnionType:
				request := ChangeDataCapture
				batchResponses := make(map[string]*quickbooks.BatchItemResponse)

				for _, typeId := range t.Types() {
					sg, ok := op.sourceGroups[typeId]
					if !ok {
						op.propagateError(fmt.Errorf("no sourceGroup found for %s", typeId))
					}

					if sg.request == Normal {
						request = Normal

						batchData, ok := sg.batchPages[page]
						if !ok {
							op.propagateError(fmt.Errorf("no batch data for sourceGroup %s, page %d", typeId, page))
						}
						batchResponses[typeId] = batchData
					}
				}

				switch request {
				case ChangeDataCapture:
					if op.changeDataCapture == nil {
						op.propagateError(fmt.Errorf("nil reference for op.changeDataCapture"))
					}

					items, err := t.ProcessCDCQuery(op.changeDataCapture, pageSize)
					if err != nil {
						op.propagateError(fmt.Errorf("error processing changeDataCapture: %w", err))
					}

					resp := OperationDataHandlerResponse{
						DataHandlerResponse: fibery.DataHandlerResponse{
							Items:               items,
							SynchronizationType: fibery.Delta,
						},
					}

					resp.DataHandlerResponse.Pagination.HasNext = false
					ch <- resp
					close(ch)

					for _, typeId := range t.Types() {
						sg := op.sourceGroups[typeId]
						sg.expectedUses--
						if sg.expectedUses == 0 {
							delete(op.sourceGroups, typeId)
						}

					}

					delete(op.requestTypes, t.Id())
				case Normal:
					items, moreSource, err := t.ProcessBatchQuery(batchResponses, pageSize)
					if err != nil {
						op.propagateError(fmt.Errorf("error processing batch query: %w", err))
					}

					resp := OperationDataHandlerResponse{
						DataHandlerResponse: fibery.DataHandlerResponse{
							Items:               items,
							SynchronizationType: fibery.Full,
						},
					}

					if len(moreSource) > 0 {
						nextKey := ResponseChannelKey(t.Id(), page+1)
						op.chans[nextKey] = make(chan OperationDataHandlerResponse, 1)
						for sourceType := range moreSource {
							nextBatchSources[sourceType] = struct{}{}
						}

						resp.DataHandlerResponse.Pagination.HasNext = true
						resp.DataHandlerResponse.Pagination.NextPageConfig.Page = page + 1
						ch <- resp
						close(ch)
					} else {
						resp.DataHandlerResponse.Pagination.HasNext = false
						ch <- resp
						close(ch)

						for _, typeId := range t.Types() {
							sg := op.sourceGroups[typeId]
							sg.expectedUses--
							if sg.expectedUses == 0 {
								delete(op.sourceGroups, typeId)
							}

						}
						delete(op.requestTypes, t.Id())
					}
				}
				continue
			default:
				op.propagateError(fmt.Errorf("unsupported type %T", regType))
				close(ch)
				delete(op.requestTypes, typeId)
				continue
			}

			grp, exists := op.sourceGroups[src]
			if !exists {
				op.propagateError(fmt.Errorf("no sourceGroup for %s", src))
				close(ch)
				delete(op.requestTypes, typeId)
				continue
			}

			var (
				resp fibery.DataHandlerResponse
				err  error
				more bool
			)

			switch grp.request {
			case ChangeDataCapture:
				switch t := regType.(type) {
				case CDCType:
					var items []map[string]any
					items, err = t.ProcessCDCQuery(
						op.changeDataCapture,
						grp.attachables,
						pageSize,
					)
					resp = fibery.DataHandlerResponse{
						Items:               items,
						SynchronizationType: fibery.Delta,
					}

				case CDCDependentType:
					var items []map[string]any
					items, err = t.ProcessCDCQuery(
						op.changeDataCapture,
						op.idCache,
						pageSize,
					)
					resp = fibery.DataHandlerResponse{
						Items:               items,
						SynchronizationType: fibery.Delta,
					}

				default:
					op.propagateError(fmt.Errorf("type %T not CDC-capable", regType))
					close(ch)
					delete(op.requestTypes, typeId)
					continue
				}
			case Normal:
				batchResp, ok := grp.batchPages[page]
				if !ok {
					op.propagateError(fmt.Errorf("no batch data for %s page %d", src, page))
					close(ch)
					delete(op.requestTypes, typeId)
					continue
				}

				switch t := regType.(type) {
				case StandardType:
					var items []map[string]any
					items, more, err = t.ProcessBatchQuery(
						batchResp,
						grp.attachables,
						pageSize,
					)
					resp = fibery.DataHandlerResponse{
						Items:               items,
						SynchronizationType: fibery.Full,
						Pagination:          fibery.Pagination{HasNext: more},
					}

				case StandardDependentType:
					var items []map[string]any
					items, more, err = t.ProcessBatchQuery(
						batchResp,
						op.idCache,
						pageSize,
					)
					resp = fibery.DataHandlerResponse{
						Items:               items,
						SynchronizationType: fibery.Full,
						Pagination:          fibery.Pagination{HasNext: more},
					}

				default:
					op.propagateError(fmt.Errorf("type %T not Batch-capable", regType))
					close(ch)
					delete(op.requestTypes, typeId)
					continue
				}

				if more {
					resp.Pagination.NextPageConfig.Page = page + 1
					nextKey := ResponseChannelKey(typeId, page+1)
					op.chans[nextKey] = make(chan OperationDataHandlerResponse, 1)
					nextBatchSources[src] = struct{}{}
				}
			}

			ch <- OperationDataHandlerResponse{Error: err, DataHandlerResponse: resp}
			close(ch)

			grp.expectedUses--
			if grp.expectedUses == 0 {
				delete(op.sourceGroups, src)
			}

			if grp.request == ChangeDataCapture || !more {
				delete(op.requestTypes, typeId)
			}
		}

		if len(nextBatchSources) == 0 {
			break
		}

		page++
		batchReqs := make([]quickbooks.BatchItemRequest, 0, len(nextBatchSources))
		for src := range nextBatchSources {
			batchReqs = append(batchReqs,
				batchQueryRequest(src, nil, page, pageSize, false),
			)
		}
		resps, err := op.integration.client.BatchRequest(params, batchReqs)
		if err != nil {
			op.propagateError(fmt.Errorf("batch page %d failed: %w", page, err))
			return
		}
		if _, err := op.indexBatchQuery(resps); err != nil {
			op.propagateError(fmt.Errorf("indexing page %d: %w", page, err))
			return
		}
	}

	op.TryCleanup()
}
