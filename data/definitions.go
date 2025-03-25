package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type Request struct {
	Filter        map[string]any
	LastSynced    time.Time
	Ctx           context.Context
	Types         []string
	RealmId       string
	StartPosition int
	PageSize      int
	OpCache       *OperationCache
	IdCache       *IdCache
	Client        *quickbooks.Client
	Token         *quickbooks.BearerToken
}

type Type interface {
	Id() string
	Name() string
	Schema() map[string]fibery.Field
	GetData(req Request) (fibery.DataHandlerResponse, error)
}

type CDCQueryable interface {
	Type
	processCDC(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error)
}

type DepCDCQueryable interface {
	Type
	processCDC(req Request, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error)
}

type WHReceivable interface {
	Type
	ProcessWH(req Request, batchResponse *quickbooks.BatchItemResponse, resp *fibery.WebhookData) error
}

type DepWHReceivable interface {
	Type
	sourceTypeId() string
	processWH(req Request, sourceData any) ([]map[string]any, error)
}

type schemaGenFunc[T any] func(T) (map[string]any, error)
type pageQueryFunc[T any] func(req Request) ([]T, error)
type entityFieldValueFunc[T, V any] func(T) V

type depSchemaGenFunc[ST any] func(ST) ([]map[string]any, error)
type sourceMapperFunc[ST any] func(ST) map[string]struct{}

func createResponse(items []map[string]any, req Request, syncType fibery.SyncType) fibery.DataHandlerResponse {
	return fibery.DataHandlerResponse{
		Items: items,
		Pagination: fibery.Pagination{
			HasNext: len(items) == req.PageSize,
			NextPageConfig: fibery.NextPageConfig{
				StartPosition: req.StartPosition + req.PageSize,
			},
		},
		SynchronizationType: syncType,
	}
}

func getFullData[T any](req Request, typeId string, pageQuery pageQueryFunc[T], process func(data []T) ([]map[string]any, error)) (fibery.DataHandlerResponse, error) {
	key := fmt.Sprintf("%s:%d", typeId, req.StartPosition)
	opCache := req.OpCache
	var data []T

	opCache.Lock()
	entry, exists := opCache.Results[key]
	if exists {
		slog.Debug(fmt.Sprintf("opCache exists for: %s", key))

		opCache.RefreshAccessTime()
		opCache.Unlock()
		var cacheData any
		select {
		case cacheData = <-entry.Data:
			var ok bool
			data, ok = cacheData.([]T)
			if !ok {
				return fibery.DataHandlerResponse{}, fmt.Errorf("cacheData type mismatch")
			}
		case <-opCache.ctx.Done():
			return fibery.DataHandlerResponse{}, nil
		}
	} else {
		slog.Debug(fmt.Sprintf("opCache does not exist for: %s", key))

		entry = &CacheEntry{Data: make(chan any, 1)}
		opCache.Results[key] = entry
		opCache.Unlock()

		slog.Debug("querying page")

		d, err := pageQuery(req)
		if err != nil {
			opCache.Lock()
			delete(opCache.Results, key)
			opCache.Unlock()
			return fibery.DataHandlerResponse{}, err
		}

		slog.Debug("data returned")

		entry.Data <- d
		data = d

		slog.Debug("data set")
	}

	slog.Debug("processing data")

	items, err := process(data)
	if err != nil {
		return fibery.DataHandlerResponse{}, fmt.Errorf("error processing data: %w", err)
	}

	slog.Debug("data processed")

	resp := createResponse(items, req, fibery.Full)

	slog.Debug("response created")

	if !resp.Pagination.HasNext {
		opCache.MarkTypeComplete(typeId)
	}

	slog.Debug("response ready")

	return resp, nil
}

func getFullDataDep[ST any](req Request, typeId string, pageQuery func(req Request) ([]ST, error), process func(data []ST) ([]map[string]any, error), update func(data []ST)) (fibery.DataHandlerResponse, error) {
	key := fmt.Sprintf("%s:%d", typeId, req.StartPosition)
	opCache := req.OpCache
	var data []ST

	opCache.Lock()
	entry, exists := opCache.Results[key]
	if exists {
		opCache.RefreshAccessTime()
		opCache.Unlock()
		var cacheData any
		select {
		case cacheData = <-entry.Data:
			var ok bool
			data, ok = cacheData.([]ST)
			if !ok {
				return fibery.DataHandlerResponse{}, fmt.Errorf("cacheData type mismatch")
			}
		case <-opCache.ctx.Done():
			return fibery.DataHandlerResponse{}, nil
		}
	} else {
		entry = &CacheEntry{Data: make(chan any, 1)}
		opCache.Results[key] = entry
		opCache.Unlock()

		d, err := pageQuery(req)
		if err != nil {
			opCache.Lock()
			delete(opCache.Results, key)
			opCache.Unlock()
			return fibery.DataHandlerResponse{}, err
		}
		entry.Data <- d
		data = d
	}

	if update != nil {
		update(data)
	}

	items, err := process(data)
	if err != nil {
		return fibery.DataHandlerResponse{}, fmt.Errorf("error processing data: %w", err)
	}

	resp := createResponse(items, req, fibery.Full)
	if !resp.Pagination.HasNext {
		opCache.MarkTypeComplete(typeId)
	}
	return resp, nil
}

func getDeltaData(req Request, key string, process func(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error)) (fibery.DataHandlerResponse, error) {
	opCache := req.OpCache
	var cdc quickbooks.ChangeDataCapture
	opCache.Lock()
	entry, exists := opCache.Results[key]
	if exists {
		opCache.RefreshAccessTime()
		opCache.Unlock()
		var cacheData any
		select {
		case cacheData = <-entry.Data:
			var ok bool
			cdc, ok = cacheData.(quickbooks.ChangeDataCapture)
			if !ok {
				return fibery.DataHandlerResponse{}, fmt.Errorf("cacheData is not quickbooks.ChangeDataCapture type")
			}
		case <-opCache.ctx.Done():
			return fibery.DataHandlerResponse{}, nil
		}
	} else {
		opCache.Unlock()
		return fibery.DataHandlerResponse{}, fmt.Errorf("delta key not found in cache")
	}

	items, err := process(&cdc)
	if err != nil {
		return fibery.DataHandlerResponse{}, fmt.Errorf("error processing delta data: %w", err)
	}

	return fibery.DataHandlerResponse{
		Items:               items,
		Pagination:          fibery.Pagination{HasNext: false},
		SynchronizationType: fibery.Delta,
	}, nil
}

func getDeltaDataDep(req Request, key string, process func(req Request, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error)) (fibery.DataHandlerResponse, error) {
	opCache := req.OpCache
	var cdc quickbooks.ChangeDataCapture
	opCache.Lock()
	entry, exists := opCache.Results[key]
	if exists {
		opCache.RefreshAccessTime()
		opCache.Unlock()
		var cacheData any
		select {
		case cacheData = <-entry.Data:
			var ok bool
			cdc, ok = cacheData.(quickbooks.ChangeDataCapture)
			if !ok {
				return fibery.DataHandlerResponse{}, fmt.Errorf("cacheData is not quickbooks.ChangeDataCapture type")
			}
		case <-opCache.ctx.Done():
			return fibery.DataHandlerResponse{}, nil
		}
	} else {
		opCache.Unlock()
		return fibery.DataHandlerResponse{}, fmt.Errorf("delta key not found in cache")
	}

	items, err := process(req, &cdc)
	if err != nil {
		return fibery.DataHandlerResponse{}, fmt.Errorf("error processing delta data: %w", err)
	}

	return fibery.DataHandlerResponse{
		Items:               items,
		Pagination:          fibery.Pagination{HasNext: false},
		SynchronizationType: fibery.Delta,
	}, nil
}

type QuickBooksType[T any] struct {
	fibery.BaseType
	schemaGen schemaGenFunc[T]
	pageQuery pageQueryFunc[T]
}

func (t QuickBooksType[T]) processQuery(entities []T) ([]map[string]any, error) {
	items := []map[string]any{}
	for _, entity := range entities {
		item, err := t.schemaGen(entity)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
		}
		items = append(items, item)
	}
	return items, nil
}

func (t QuickBooksType[T]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Full:
		return getFullData(req, t.Id(), t.pageQuery, t.processQuery)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}

}

type QuickBooksCDCType[T any] struct {
	QuickBooksType[T]
	entityId     entityFieldValueFunc[T, string]
	entityStatus entityFieldValueFunc[T, string]
}

func (t QuickBooksCDCType[T]) processCDC(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	entities := quickbooks.CDCQueryExtractor[T](cdc)
	items := []map[string]any{}
	for _, entity := range entities {
		if t.entityStatus(entity) == "Deleted" {
			item := map[string]any{
				"id":           t.entityId(entity),
				"__syncAction": fibery.REMOVE,
			}
			items = append(items, item)
		} else {
			item, err := t.schemaGen(entity)
			if err != nil {
				return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (t QuickBooksCDCType[T]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Delta:
		return getDeltaData(req, "cdc", t.processCDC)
	case fibery.Full:
		return getFullData(req, t.Id(), t.pageQuery, t.processQuery)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

// QuickBooksWHType established the additional function(s) required to process a webhook notifcation
type QuickBooksWHType[T any] struct {
	QuickBooksType[T]
	entityId entityFieldValueFunc[T, string]
}

func (t QuickBooksWHType[T]) ProcessWH(req Request, batchResponse *quickbooks.BatchItemResponse, resp *fibery.WebhookData) error {
	entities := quickbooks.BatchQueryExtractor[T](batchResponse)
	items, err := t.processQuery(entities)
	if err != nil {
		return fmt.Errorf("unable to process %s query", t.Id())
	}
	(*resp)[t.Id()] = append((*resp)[t.Id()], items...)
	if dependents, ok := Types.DepWHReceivable[t.Id()]; ok {
		for _, depPtr := range dependents {
			for _, selectedType := range req.Types {
				if depType := *depPtr; depType.Id() == selectedType {
					depItems, err := depType.processWH(req, entities)
					if err != nil {
						return fmt.Errorf("unable to process data for dependent: %s", depType.Id())
					}
					(*resp)[depType.Id()] = append((*resp)[depType.Id()], depItems...)
				}
				continue
			}

		}
	}
	return nil
}

// QuickBooksDualType requires the functions from both QuickBooksCDCType and QuickBooksWHType
type QuickBooksDualType[T any] struct {
	QuickBooksType[T]
	entityId     entityFieldValueFunc[T, string]
	entityStatus entityFieldValueFunc[T, string]
}

func (t QuickBooksDualType[T]) ProcessWH(req Request, batchResponse *quickbooks.BatchItemResponse, resp *fibery.WebhookData) error {
	entities := quickbooks.BatchQueryExtractor[T](batchResponse)
	items, err := t.processQuery(entities)
	if err != nil {
		return fmt.Errorf("unable to process %s query", t.Id())
	}
	(*resp)[t.Id()] = append((*resp)[t.Id()], items...)
	if dependents, ok := Types.DepWHReceivable[t.Id()]; ok {
		for _, depPtr := range dependents {
			for _, selectedType := range req.Types {
				if depType := *depPtr; depType.Id() == selectedType {
					depItems, err := depType.processWH(req, entities)
					if err != nil {
						return fmt.Errorf("unable to process data for dependent: %s", depType.Id())
					}
					(*resp)[depType.Id()] = append((*resp)[depType.Id()], depItems...)
				}
				continue
			}

		}
	}
	return nil
}

func (t QuickBooksDualType[T]) processCDC(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	entities := quickbooks.CDCQueryExtractor[T](cdc)
	items := []map[string]any{}
	for _, entity := range entities {
		if t.entityStatus(entity) == "Deleted" {
			item := map[string]any{
				"id":           t.entityId(entity),
				"__syncAction": fibery.REMOVE,
			}
			items = append(items, item)
		} else {
			item, err := t.schemaGen(entity)
			if err != nil {
				return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (t QuickBooksDualType[T]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Delta:
		return getDeltaData(req, "cdc", t.processCDC)
	case fibery.Full:
		return getFullData(req, t.Id(), t.pageQuery, t.processQuery)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

// DependentBaseType established the base functions required to process, extract, and convert dependent data from an array of source entities
type dependentBaseType[ST any] struct {
	fibery.BaseType
	schemaGen depSchemaGenFunc[ST]
}

func (t dependentBaseType[ST]) processQuery(sourceEntities []ST) ([]map[string]any, error) {
	items := []map[string]any{}
	for _, source := range sourceEntities {
		itemSlice, err := t.schemaGen(source)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
		}
		items = append(items, itemSlice...)
	}
	return items, nil
}

// DependentDataType corresponds to a QuickBooksType which can only be requested through a query or read operation
type DependentDataType[ST any] struct {
	dependentBaseType[ST]
	sourceType *QuickBooksType[ST]
}

func (t DependentDataType[ST]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Full:
		return getFullData(req, t.Id(), t.sourceType.pageQuery, t.processQuery)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

type DependentCDCType[ST any] struct {
	dependentBaseType[ST]
	sourceType   *QuickBooksCDCType[ST]
	sourceId     entityFieldValueFunc[ST, string]
	sourceStatus entityFieldValueFunc[ST, string]
	sourceMapper sourceMapperFunc[ST]
}

func (t DependentCDCType[ST]) processCDC(req Request, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	sourceEntities := quickbooks.CDCQueryExtractor[ST](cdc)
	items := []map[string]any{}
	idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
	if !exists {
		return nil, fmt.Errorf("no id cache entry found for %s:%s, perform full sync to populate cache", req.RealmId, t.Id())
	}

	for _, source := range sourceEntities {
		sourceId := t.sourceId(source)
		cachedIds := idSet[sourceId]
		newIds := t.sourceMapper(source)

		if t.sourceStatus(source) == "Deleted" {
			for cachedId := range cachedIds {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
			req.IdCache.RemoveSource(req.RealmId, t.Id(), sourceId)
			continue
		}

		itemSlice, err := t.schemaGen(source)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
		}
		items = append(items, itemSlice...)

		for cachedId := range cachedIds {
			if _, ok := newIds[cachedId]; !ok {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
		}

		idSet[sourceId] = newIds
	}

	req.IdCache.Set(req.RealmId, t.Id(), idSet)
	return items, nil
}

func (t DependentCDCType[ST]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Delta:
		return getDeltaDataDep(req, "cdc", t.processCDC)
	case fibery.Full:
		cacheUpdateFunc := func(sourceEntities []ST) {
			idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
			if !exists {
				idSet = make(IdSet)
			}
			for _, source := range sourceEntities {
				sId := t.sourceId(source)
				newIds := t.sourceMapper(source)
				idSet[sId] = newIds
			}
			req.IdCache.Set(req.RealmId, t.Id(), idSet)
		}
		return getFullDataDep(req, t.sourceType.Id(), t.sourceType.pageQuery, t.processQuery, cacheUpdateFunc)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

type DependentWHType[ST any] struct {
	dependentBaseType[ST]
	sourceType   *QuickBooksWHType[ST]
	sourceId     entityFieldValueFunc[ST, string]
	sourceMapper sourceMapperFunc[ST]
}

func (t DependentWHType[ST]) sourceTypeId() string {
	return t.sourceType.Id()
}

func (t DependentWHType[ST]) processWH(req Request, sourceEntities []ST) ([]map[string]any, error) {
	items := []map[string]any{}
	idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
	if !exists {
		return nil, fmt.Errorf("no id cache entry found for %s:%s, perform full sync to populate cache", req.RealmId, t.Id())
	}

	for _, source := range sourceEntities {
		sourceId := t.sourceId(source)
		cachedIds := idSet[sourceId]
		newIds := t.sourceMapper(source)

		itemSlice, err := t.schemaGen(source)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
		}
		items = append(items, itemSlice...)

		for cachedId := range cachedIds {
			if _, ok := newIds[cachedId]; !ok {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
		}

		idSet[sourceId] = newIds

	}

	req.IdCache.Set(req.RealmId, t.Id(), idSet)
	return items, nil
}

func (t DependentWHType[ST]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Full:
		cacheUpdateFunc := func(sourceEntities []ST) {
			idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
			if !exists {
				idSet = make(IdSet)
			}
			for _, source := range sourceEntities {
				sId := t.sourceId(source)
				newIds := t.sourceMapper(source)
				idSet[sId] = newIds
			}
			req.IdCache.Set(req.RealmId, t.Id(), idSet)
		}
		return getFullDataDep(req, t.sourceType.Id(), t.sourceType.pageQuery, t.processQuery, cacheUpdateFunc)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

type DependentDualType[ST any] struct {
	dependentBaseType[ST]
	sourceType   *QuickBooksDualType[ST]
	sourceId     entityFieldValueFunc[ST, string]
	sourceStatus entityFieldValueFunc[ST, string]
	sourceMapper sourceMapperFunc[ST]
}

func (t DependentDualType[ST]) sourceTypeId() string {
	return t.sourceType.Id()
}

func (t DependentDualType[ST]) processWH(req Request, data any) ([]map[string]any, error) {
	sourceEntities, ok := data.([]ST)
	if !ok {
		return nil, fmt.Errorf("unable to assert sourceType on data")
	}

	items := []map[string]any{}
	idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
	if !exists {
		return nil, fmt.Errorf("no id cache entry found for %s:%s, perform full sync to populate cache", req.RealmId, t.Id())
	}

	for _, source := range sourceEntities {
		sourceId := t.sourceId(source)
		cachedIds := idSet[sourceId]
		newIds := t.sourceMapper(source)

		itemSlice, err := t.schemaGen(source)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
		}
		items = append(items, itemSlice...)

		for cachedId := range cachedIds {
			if _, ok := newIds[cachedId]; !ok {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
		}

		idSet[sourceId] = newIds

	}

	req.IdCache.Set(req.RealmId, t.Id(), idSet)
	return items, nil
}

func (t DependentDualType[ST]) processCDC(req Request, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	sourceEntities := quickbooks.CDCQueryExtractor[ST](cdc)
	items := []map[string]any{}
	idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
	if !exists {
		return nil, fmt.Errorf("no id cache entry found for %s:%s, perform full sync to populate cache", req.RealmId, t.Id())
	}

	for _, source := range sourceEntities {
		sourceId := t.sourceId(source)
		cachedIds := idSet[sourceId]
		newIds := t.sourceMapper(source)

		if t.sourceStatus(source) == "Deleted" {
			for cachedId := range cachedIds {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
			req.IdCache.RemoveSource(req.RealmId, t.Id(), sourceId)
			continue
		}

		itemSlice, err := t.schemaGen(source)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema", t.Id())
		}
		items = append(items, itemSlice...)

		for cachedId := range cachedIds {
			if _, ok := newIds[cachedId]; !ok {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
		}

		idSet[sourceId] = newIds
	}

	req.IdCache.Set(req.RealmId, t.Id(), idSet)
	return items, nil
}

func (t DependentDualType[ST]) GetData(req Request) (fibery.DataHandlerResponse, error) {
	syncType := req.OpCache.SyncTypes[t.Id()]

	switch syncType {
	case fibery.Delta:
		return getDeltaDataDep(req, "cdc", t.processCDC)
	case fibery.Full:
		cacheUpdateFunc := func(sourceEntities []ST) {
			idSet, exists := req.IdCache.Get(req.RealmId, t.Id())
			if !exists {
				idSet = make(IdSet)
			}
			for _, source := range sourceEntities {
				sId := t.sourceId(source)
				newIds := t.sourceMapper(source)
				idSet[sId] = newIds
			}
			req.IdCache.Set(req.RealmId, t.Id(), idSet)
		}
		return getFullDataDep(req, t.sourceType.Id(), t.sourceType.pageQuery, t.processQuery, cacheUpdateFunc)
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

type TypeRegistry struct {
	All             map[string]*Type
	DepWHReceivable map[string][]*DepWHReceivable
}

var Types = TypeRegistry{
	All:             make(map[string]*Type),
	DepWHReceivable: make(map[string][]*DepWHReceivable),
}

func registerType(t Type) {
	Types.All[t.Id()] = &t
	depType, ok := t.(DepWHReceivable)
	if ok {
		Types.DepWHReceivable[depType.sourceTypeId()] = append(Types.DepWHReceivable[depType.sourceTypeId()], &depType)
	}
}
