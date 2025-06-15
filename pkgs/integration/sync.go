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
		slog.Debug(fmt.Sprintf("requested operation: %s exists", req.OperationId))
		return op, nil
	}

	slog.Debug(fmt.Sprintf("requested operation: %s does not exist", req.OperationId))

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
			slog.Warn(fmt.Sprintf("operation %s exceeded timeout of %v, cancelling and cleaning up", opId, om.ttl))

			// Cancel the operation context to stop any ongoing work
			if op.cancel != nil {
				op.cancel()
			}

			// Propagate timeout error to all waiting channels
			timeoutErr := fmt.Errorf("operation %s timed out after %v without completion or new requests", opId, om.ttl)
			op.propagateError(timeoutErr)

			// Remove the operation from the manager
			delete(om.Operations, opId)

			slog.Debug(fmt.Sprintf("operation %s cleaned up due to timeout", opId))
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
		lastRequest:   time.Now(),
		chans:         make(map[string]chan OperationDataHandlerResponse, len(req.Types)),
		ctx:           ctx,
		cancel:        cancel,
	}
	op.wg.Add(len(req.Types))
	slog.Debug(fmt.Sprintf("operation: %s built with %d waitgroup", op.id, len(req.Types)))
	return op
}

func (op *Operation) propagateError(err error) {
	op.cancel()

	op.cleanupOnce.Do(func() {
		op.Lock()
		keys := make([]string, 0, len(op.chans))
		for k := range op.chans {
			keys = append(keys, k)
		}
		op.Unlock()

		for _, k := range keys {
			op.completeChannel(k, OperationDataHandlerResponse{Error: err})
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
	defer op.wg.Done()

	select {
	case <-op.ctx.Done():
		return fmt.Errorf("operation %s has been cancelled", op.id)
	default:
	}

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

			return nil

		case UnionType:
			schema, ok := req.Schema[req.RequestedType]
			if !ok {
				return fmt.Errorf("no schema for %s was provided on the request", req.RequestedType)
			}

			getAttach := attachableField(schema, attachableFieldId)

			for _, sourceType := range t.Types() {
				innerReq := req
				if !t.CDC() {
					innerReq.LastSyncronizedAt = time.Time{}
				}

				if err := op.processTypeEntry(sourceType, innerReq, getAttach); err != nil {
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
		slog.Debug(fmt.Sprintf("type: %s submitted, waitgroup decremented", req.RequestedType))
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
	select {
	case <-op.ctx.Done():
		return nil, fmt.Errorf("operation %s has been cancelled", op.id)
	default:
	}

	op.Lock()
	defer op.Unlock()

	ch, ok := op.chans[key]
	if !ok {
		return nil, fmt.Errorf("no channel found for key: %s", key)
	}

	return ch, nil
}

func (op *Operation) completeChannel(key string, resp OperationDataHandlerResponse) {
	op.Lock()
	ch, ok := op.chans[key]
	if ok {
		delete(op.chans, key)
	}
	op.Unlock()

	if !ok {
		return
	}
	ch <- resp
	close(ch)
}

func (op *Operation) TryCleanup() {
	op.Lock()
	shouldDelete := len(op.requestTypes) == 0
	opId := op.id
	op.Unlock()

	if shouldDelete {
		op.integration.opManager.DeleteOperation(opId)
		slog.Debug(fmt.Sprintf("operation %s deleted", opId))
	}
}

func (op *Operation) indexBatchQuery(batch []quickbooks.BatchItemResponse) (map[string]struct{}, error) {
	moreAttachables := map[string]struct{}{}
	op.Lock()
	for _, resp := range batch {
		faults := resp.Fault.Faults
		if len(faults) > 0 {
			return nil, fmt.Errorf("fault for %s: %w", resp.BID, quickbooks.BatchError{Faults: faults})
		}

		entityType, startPosition, isAttachable, err := DecodeQueryBID(resp.BID)
		if err != nil {
			return nil, fmt.Errorf("error decoding query BID: %w", err)
		}

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
	op.Unlock()

	return moreAttachables, nil
}

func (op *Operation) doBatch(req []quickbooks.BatchItemRequest, params quickbooks.RequestParameters) {
	client := op.integration.client

	page := 1

	slog.Debug("starting inital batch loop")

	for {
		slog.Debug("in loop")
		batch, err := client.BatchRequest(params, req)
		if err != nil {
			slog.Error(fmt.Sprintf("error fetching inital batch: %s", err.Error()))
			op.propagateError(fmt.Errorf("error fetching inital batch: %w", err))
			return
		}

		slog.Debug("batch request complete")

		nextAttachEntities, err := op.indexBatchQuery(batch)
		if err != nil {
			slog.Error(fmt.Sprintf("error indexing inital batch: %s", err.Error()))
			op.propagateError(fmt.Errorf("error indexing inital batch: %w", err))
			return
		}

		slog.Debug(fmt.Sprintf("inital batch page query: %d complete", page))

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
	slog.Debug("inital batch complete")
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
	slog.Debug("inital cdc complete")
}

func (op *Operation) fetchAll() error {
	slog.Debug("fetch started")

	select {
	case <-op.ctx.Done():
		return fmt.Errorf("operation %s has been cancelled", op.id)
	default:
	}

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
		slog.Debug("making cdc request")
		initalFetch.Add(1)
		go func(req []string) {
			defer initalFetch.Done()
			op.doCDC(req, params)
		}(cdcReq)
	}

	if len(batchReq) > 0 {
		slog.Debug("making batch request")
		initalFetch.Add(1)
		go func(req []quickbooks.BatchItemRequest) {
			defer initalFetch.Done()
			op.doBatch(req, params)
		}(batchReq)
	}

	initalFetch.Wait()

	slog.Debug("inital fetch complete")

	time.Sleep(30 * time.Second)

	op.dispatchPages()

	return nil
}

func (op *Operation) dispatchPages() error {
	select {
	case <-op.ctx.Done():
		return fmt.Errorf("operation %s has been cancelled", op.id)
	default:
	}

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
			_, ok := op.chans[key]
			if !ok {
				op.propagateError(fmt.Errorf("missing channel %s", key))
				break
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

				for _, sourceType := range t.Types() {
					sg, ok := op.sourceGroups[sourceType.Type()]
					if !ok {
						op.propagateError(fmt.Errorf("no sourceGroup found for %s", sourceType.Type()))
						break
					}

					if sg.request == Normal {
						request = Normal

						batchData, ok := sg.batchPages[page]
						if !ok {
							op.propagateError(fmt.Errorf("no batch data for sourceGroup %s, page %d", sourceType.Type(), page))
							break
						}
						batchResponses[sourceType.Type()] = batchData
					}
				}

				switch request {
				case ChangeDataCapture:
					if op.changeDataCapture == nil {
						op.propagateError(fmt.Errorf("nil reference for op.changeDataCapture"))
						break
					}

					items, err := t.ProcessCDCQuery(op.changeDataCapture, pageSize)
					if err != nil {
						op.propagateError(fmt.Errorf("error processing changeDataCapture: %w", err))
						break
					}

					resp := OperationDataHandlerResponse{
						DataHandlerResponse: fibery.DataHandlerResponse{
							Items:               items,
							SynchronizationType: fibery.Delta,
						},
					}

					resp.DataHandlerResponse.Pagination.HasNext = false
					op.completeChannel(key, resp)

					for _, sourceType := range t.Types() {
						sg := op.sourceGroups[sourceType.Type()]
						sg.expectedUses--
						if sg.expectedUses == 0 {
							delete(op.sourceGroups, sourceType.Type())
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
						op.completeChannel(key, resp)
					} else {
						resp.DataHandlerResponse.Pagination.HasNext = false
						op.completeChannel(key, resp)

						for _, sourceType := range t.Types() {
							sg := op.sourceGroups[sourceType.Type()]
							sg.expectedUses--
							if sg.expectedUses == 0 {
								delete(op.sourceGroups, sourceType.Type())
							}
						}

						delete(op.requestTypes, t.Id())
					}
				}
				continue
			default:
				op.propagateError(fmt.Errorf("unsupported type %T", regType))
				break
			}

			grp, exists := op.sourceGroups[src]
			if !exists {
				op.propagateError(fmt.Errorf("no sourceGroup for %s", src))
				break
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
					break
				}
			case Normal:
				batchResp, ok := grp.batchPages[page]
				if !ok {
					op.propagateError(fmt.Errorf("no batch data for %s page %d", src, page))
					break
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
					break
				}

				if more {
					resp.Pagination.NextPageConfig.Page = page + 1
					nextKey := ResponseChannelKey(typeId, page+1)
					op.chans[nextKey] = make(chan OperationDataHandlerResponse, 1)
					nextBatchSources[src] = struct{}{}
				}
			}

			op.completeChannel(key, OperationDataHandlerResponse{Error: err, DataHandlerResponse: resp})

			grp.expectedUses--
			if grp.expectedUses == 0 {
				delete(op.sourceGroups, src)
			}

			if grp.request == ChangeDataCapture || !more {
				delete(op.requestTypes, typeId)
			}
		}

		slog.Debug("inital dispatch")

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
			break
		}
		if _, err := op.indexBatchQuery(resps); err != nil {
			op.propagateError(fmt.Errorf("indexing page %d: %w", page, err))
			break
		}
	}

	op.TryCleanup()
	return nil
}
