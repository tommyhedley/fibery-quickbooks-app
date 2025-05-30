package main

import (
	"fmt"
	"log/slog"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type Type2 interface {
	Id() string
	Name() string
	Schema() map[string]fibery.Field
	SelectStatement(start, pageSize int) string
	GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error)
}

type DeltaDepType interface {
	Type
	SourceId() string
	SourceKey(sourceEntity any) (IdKey, error)
	MapSource(sourceEntity any) (map[string]struct{}, error)
}

type UnionType2 interface {
	Type
	SourceIds() []string
	UnionTypes() []Type
}

type CDCType2 interface {
	Type
	processCDC(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error)
}

type WebhookType2 interface {
	Type
	ProcessWebhookUpdate(batchData *quickbooks.BatchItemResponse, resp *fibery.WebhookTransformResponse) error
	ProcessWebhookDeletions(ids []string, resp *fibery.WebhookTransformResponse)
	GetRelatedTypes() []CDCType2
}

type CDCDepType interface {
	DeltaDepType
	processCDC(cache *IdCache, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error)
}

type WebhookDepType interface {
	DeltaDepType
	ProcessWebhookUpdate(batchData *quickbooks.BatchItemResponse, resp *fibery.WebhookTransformResponse, idCache *IdCache) error
	ProcessWebhookDeletions(sourceIds []string, resp *fibery.WebhookTransformResponse, idCache *IdCache)
}

type schemaGenFunc[T any] func(T) (map[string]any, error)
type pageQueryFunc[T any] func(client *quickbooks.Client, requestParams quickbooks.RequestParameters, startPosition, pageSize int) ([]T, error)
type entityFieldValueFunc[T, V any] func(T) V

type depSchemaGenFunc[ST any] func(ST) ([]map[string]any, error)

type QuickBooksType[T any] struct {
	fibery.Type
	FieldDef[T]
	schemaGen schemaGenFunc[T]
	pageQuery pageQueryFunc[T]
}

type QuickBooksCDCType[T any] struct {
	QuickBooksType[T]
	entityId     entityFieldValueFunc[T, string]
	entityStatus entityFieldValueFunc[T, string]
}

type QuickBooksWHType[T any] struct {
	QuickBooksType[T]
	entityId     entityFieldValueFunc[T, string]
	relatedTypes []CDCType
}

type QuickBooksDualType[T any] struct {
	QuickBooksType[T]
	entityId     entityFieldValueFunc[T, string]
	entityStatus entityFieldValueFunc[T, string]
	relatedTypes []CDCType
}

type dependentBaseType[ST, T any] struct {
	fibery.Type
	DependentFieldDef[ST, T]
	schemaGen depSchemaGenFunc[ST]
}

type DependentDataType[ST any] struct {
	dependentBaseType[ST]
	sourceType *QuickBooksType[ST]
}

type DependentCDCType2[ST any] struct {
	dependentBaseType[ST]
	sourceType   *QuickBooksCDCType[ST]
	sourceId     entityFieldValueFunc[ST, string]
	sourceStatus entityFieldValueFunc[ST, string]
	sourceMapper sourceMapperFunc[ST]
}

type DependentWHType[ST any] struct {
	dependentBaseType[ST]
	sourceType   *QuickBooksWHType[ST]
	sourceId     entityFieldValueFunc[ST, string]
	sourceMapper sourceMapperFunc[ST]
}

type DependentDualType2[ST any] struct {
	dependentBaseType[ST]
	sourceType   *QuickBooksDualType[ST]
	sourceId     entityFieldValueFunc[ST, string]
	sourceStatus entityFieldValueFunc[ST, string]
	sourceMapper sourceMapperFunc[ST]
}

type UnionDataType struct {
	fibery.BaseType
	unionTypes []Type
	schemaGen  func(typeId string, input []map[string]any) ([]map[string]any, error)
}

func NewQuickBooksType[T any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg schemaGenFunc[T],
	pq pageQueryFunc[T],
) *QuickBooksType[T] {
	return &QuickBooksType[T]{
		BaseType: fibery.BaseType{
			TypeId:     id,
			TypeName:   name,
			TypeSchema: schema,
		},
		schemaGen: sg,
		pageQuery: pq,
	}
}

func NewQuickBooksCDCType[T any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg schemaGenFunc[T],
	pq pageQueryFunc[T],
	eId entityFieldValueFunc[T, string],
	eStat entityFieldValueFunc[T, string],
) *QuickBooksCDCType[T] {
	return &QuickBooksCDCType[T]{
		QuickBooksType: QuickBooksType[T]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
			pageQuery: pq,
		},
		entityId:     eId,
		entityStatus: eStat,
	}
}

func NewQuickBooksWHType[T any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg schemaGenFunc[T],
	pq pageQueryFunc[T],
	eId entityFieldValueFunc[T, string],
) *QuickBooksWHType[T] {
	return &QuickBooksWHType[T]{
		QuickBooksType: QuickBooksType[T]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
			pageQuery: pq,
		},
		entityId: eId,
	}
}

func NewQuickBooksDualType[T any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg schemaGenFunc[T],
	pq pageQueryFunc[T],
	eId entityFieldValueFunc[T, string],
	eStat entityFieldValueFunc[T, string],
) *QuickBooksDualType[T] {
	return &QuickBooksDualType[T]{
		QuickBooksType: QuickBooksType[T]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
			pageQuery: pq,
		},
		entityId:     eId,
		entityStatus: eStat,
	}
}

func NewDependentDataType[ST any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg depSchemaGenFunc[ST],
	st *QuickBooksType[ST],
) *DependentDataType[ST] {
	return &DependentDataType[ST]{
		dependentBaseType: dependentBaseType[ST]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
		},
		sourceType: st,
	}
}

func NewDependentCDCType[ST any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg depSchemaGenFunc[ST],
	st *QuickBooksCDCType[ST],
	sId entityFieldValueFunc[ST, string],
	sStat entityFieldValueFunc[ST, string],
	sMap sourceMapperFunc[ST],
) *DependentCDCType[ST] {
	return &DependentCDCType[ST]{
		dependentBaseType: dependentBaseType[ST]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
		},
		sourceType:   st,
		sourceId:     sId,
		sourceStatus: sStat,
		sourceMapper: sMap,
	}
}

func NewDependentWHType[ST any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg depSchemaGenFunc[ST],
	st *QuickBooksWHType[ST],
	sId entityFieldValueFunc[ST, string],
	sMap sourceMapperFunc[ST],
) *DependentWHType[ST] {
	return &DependentWHType[ST]{
		dependentBaseType: dependentBaseType[ST]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
		},
		sourceType:   st,
		sourceId:     sId,
		sourceMapper: sMap,
	}
}

func NewDependentDualType[ST any](
	id string,
	name string,
	schema map[string]fibery.Field,
	sg depSchemaGenFunc[ST],
	st *QuickBooksDualType[ST],
	sId entityFieldValueFunc[ST, string],
	sStat entityFieldValueFunc[ST, string],
	sMap sourceMapperFunc[ST],
) *DependentDualType[ST] {
	return &DependentDualType[ST]{
		dependentBaseType: dependentBaseType[ST]{
			BaseType: fibery.BaseType{
				TypeId:     id,
				TypeName:   name,
				TypeSchema: schema,
			},
			schemaGen: sg,
		},
		sourceType:   st,
		sourceId:     sId,
		sourceStatus: sStat,
		sourceMapper: sMap,
	}
}

func NewUnionDataType(
	id string,
	name string,
	schema map[string]fibery.Field,
	ut []Type,
	sg func(typeId string, input []map[string]any) ([]map[string]any, error),
) *UnionDataType {
	return &UnionDataType{
		BaseType: fibery.BaseType{
			TypeId:     id,
			TypeName:   name,
			TypeSchema: schema,
		},
		unionTypes: ut,
		schemaGen:  sg,
	}
}

func getData[T any](
	client *quickbooks.Client,
	op *Operation,
	storedType Type,
	startPosition int,
	pageSize int,
	queryByPage pageQueryFunc[T],
	processQuery func([]T) ([]map[string]any, error),
) (fibery.DataHandlerResponse, error) {
	requestType := op.Types[storedType.Id()]
	requestParams := quickbooks.RequestParameters{
		Ctx:             op.ctx,
		RealmId:         op.Account.RealmID,
		Token:           &op.Account.BearerToken,
		WaitOnRateLimit: true,
	}

	switch requestType.Sync {
	case fibery.Delta:
		dataRequestKey := DataKey{DataType: "CDC"}
		var expectedGroupSize int
		var cdcTypes []string
		cdcTypeMap := map[string]struct{}{}

		for _, rt := range op.Types {
			if rt.Sync == fibery.Delta {
				expectedGroupSize++
				cdcTypeMap[rt.SourceId] = struct{}{}
			}
		}

		for id := range cdcTypeMap {
			cdcTypes = append(cdcTypes, id)
			slog.Debug(fmt.Sprintf("requesting cdc data for: %s", id))
		}

		slog.Debug(fmt.Sprintf("cdc request time: %s", op.LastSynced.String()))

		result, err := op.GetOrFetchData(dataRequestKey, expectedGroupSize, func() (any, error) {
			return client.ChangeDataCapture(requestParams, cdcTypes, op.LastSynced)
		})
		if err != nil {
			return fibery.DataHandlerResponse{}, fmt.Errorf("error requesting changeDataCapture: %w", err)
		}

		cdc := result.(quickbooks.ChangeDataCapture)

		var items []map[string]any
		switch cdcType := storedType.(type) {
		case CDCType:
			items, err = cdcType.processCDC(&cdc)
			if err != nil {
				return fibery.DataHandlerResponse{}, fmt.Errorf("error processing cdc: %w", err)
			}
		case CDCDepType:
			items, err = cdcType.processCDC(op.IdCache, &cdc)
			if err != nil {
				return fibery.DataHandlerResponse{}, fmt.Errorf("error processing cdc: %w", err)
			}
		default:
			return fibery.DataHandlerResponse{}, fmt.Errorf("invalid type was passed into getChangeDataCapture: %s", cdcType.Id())
		}

		if requestType.Attachables != nil && len(requestType.Attachables) > 0 {
			for i, item := range items {
				if id, ok := item["id"].(string); ok {
					if attachments, exists := requestType.Attachables[id]; exists && len(attachments) > 0 {
						items[i]["Files"] = attachments
						slog.Debug(fmt.Sprintf("attachables linked to %s:%s", storedType.Id(), id))
						fmt.Println(items[i])
					}
				}
			}
		}

		op.MarkTypeFulfilled(storedType.Id())

		slog.Debug(fmt.Sprintf("items for %s: %s", storedType.Id(), items))

		return fibery.DataHandlerResponse{
			Items:               items,
			SynchronizationType: fibery.Delta,
		}, nil

	case fibery.Full:
		var dataRequestKey DataKey
		switch typ := storedType.(type) {
		case DeltaDepType:
			dataRequestKey = DataKey{DataType: typ.SourceId(), StartPosition: startPosition}
		default:
			dataRequestKey = DataKey{DataType: typ.Id(), StartPosition: startPosition}
		}
		expectedGroupSize := requestType.GroupSize

		result, err := op.GetOrFetchData(dataRequestKey, expectedGroupSize, func() (any, error) {
			return queryByPage(client, requestParams, startPosition, pageSize)
		})
		if err != nil {
			return fibery.DataHandlerResponse{}, fmt.Errorf("error get/fetching data: %w", err)
		}

		data := result.([]T)

		items, err := processQuery(data)
		if err != nil {
			return fibery.DataHandlerResponse{}, fmt.Errorf("error processing data: %w", err)
		}

		if deltaType, ok := storedType.(DeltaDepType); ok {
			for _, sourceEntity := range data {
				sk, err := deltaType.SourceKey(sourceEntity)
				if err != nil {
					return fibery.DataHandlerResponse{}, err
				}
				sourceMap, err := deltaType.MapSource(sourceEntity)
				if err != nil {
					return fibery.DataHandlerResponse{}, err
				}
				op.IdCache.SetIds(sk, deltaType.Id(), sourceMap)
			}
		}

		if requestType.Attachables != nil && len(requestType.Attachables) > 0 {
			for i, item := range items {
				if id, ok := item["id"].(string); ok {
					if attachments, exists := requestType.Attachables[id]; exists && len(attachments) > 0 {
						items[i]["Files"] = attachments
						slog.Debug(fmt.Sprintf("attachables linked to %s:%s", storedType.Id(), id))
						fmt.Println(items[i])
					}
				}
			}
		}

		more := len(items) == pageSize

		if !more {
			op.MarkTypeFulfilled(storedType.Id())
		}

		return fibery.DataHandlerResponse{
			Items: items,
			Pagination: fibery.Pagination{
				HasNext: more,
				NextPageConfig: fibery.NextPageConfig{
					StartPosition: startPosition + pageSize,
				},
			},
			SynchronizationType: fibery.Full,
		}, nil
	default:
		return fibery.DataHandlerResponse{}, fmt.Errorf("unsupported sync type")
	}
}

func processCDC[T any](
	cdc *quickbooks.ChangeDataCapture,
	id entityFieldValueFunc[T, string],
	status entityFieldValueFunc[T, string],
	sg schemaGenFunc[T],
) ([]map[string]any, error) {
	entities := quickbooks.CDCQueryExtractor[T](cdc)
	items := []map[string]any{}
	for _, entity := range entities {
		if status(entity) == "Deleted" {
			item := map[string]any{
				"id":           id(entity),
				"__syncAction": fibery.REMOVE,
			}
			items = append(items, item)
		} else {
			item, err := sg(entity)
			if err != nil {
				return nil, fmt.Errorf("error generating fibery schema")
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func processCDCDep[ST any](
	idCache *IdCache,
	cdc *quickbooks.ChangeDataCapture,
	sourceId entityFieldValueFunc[ST, string],
	sourceStatus entityFieldValueFunc[ST, string],
	sg depSchemaGenFunc[ST],
	sm sourceMapperFunc[ST],
	typeId, sourceTypeId string,
) ([]map[string]any, error) {
	sourceEntities := quickbooks.CDCQueryExtractor[ST](cdc)
	items := []map[string]any{}

	for _, source := range sourceEntities {
		sk := IdKey{EntityType: sourceTypeId, EntityId: sourceId(source)}
		cachedIds, exists := idCache.GetIdsByType(sk, typeId)
		if !exists {
			cachedIds = make(map[string]struct{})
		}

		newIds := sm(source)

		if sourceStatus(source) == "Deleted" {
			for cachedId := range cachedIds {
				items = append(items, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
			idCache.RemoveEntityType(sk, typeId)
			continue
		}

		itemSlice, err := sg(source)
		if err != nil {
			return nil, fmt.Errorf("error generating fibery schema")
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

		idCache.AddIds(sk, typeId, newIds)
	}

	return items, nil
}

func processWebhookUpdates[T any](
	batchData *quickbooks.BatchItemResponse,
	resp *fibery.WebhookTransformResponse,
	processQuery func([]T) ([]map[string]any, error),
	typeId string,
) error {
	entities := quickbooks.BatchQueryExtractor[T](batchData)
	items, err := processQuery(entities)
	if err != nil {
		return fmt.Errorf("error processing data: %w", err)
	}
	resp.Data[typeId] = items
	return nil
}

func processWebhookUpdatesDep[ST any](
	batchData *quickbooks.BatchItemResponse,
	resp *fibery.WebhookTransformResponse,
	typeId, sourceTypeId string,
	idCache *IdCache,
	sourceId entityFieldValueFunc[ST, string],
	sg depSchemaGenFunc[ST],
	sm sourceMapperFunc[ST],
) error {
	sourceEntities := quickbooks.BatchQueryExtractor[ST](batchData)
	items := []map[string]any{}
	for _, source := range sourceEntities {
		sk := IdKey{EntityType: sourceTypeId, EntityId: sourceId(source)}
		cachedIds, exists := idCache.GetIdsByType(sk, typeId)
		if !exists {
			cachedIds = make(map[string]struct{})
		}

		newIds := sm(source)
		itemSlice, err := sg(source)
		if err != nil {
			return fmt.Errorf("error generating fibery schema")
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

		idCache.AddIds(sk, typeId, newIds)

	}

	resp.Data[typeId] = items
	return nil
}

func processWebhookDeletions(ids []string, resp *fibery.WebhookTransformResponse, typeId string) {
	for _, id := range ids {
		resp.Data[typeId] = append(resp.Data[typeId], map[string]any{
			"id":           id,
			"__syncAction": fibery.REMOVE,
		})
	}
}

func processWebhookDeletionsDep(sourceIds []string, resp *fibery.WebhookTransformResponse, typeId, sourceTypeId string, idCache *IdCache) {
	for _, sourceId := range sourceIds {
		sk := IdKey{EntityType: sourceTypeId, EntityId: sourceId}
		cachedIds, exists := idCache.GetIdsByType(sk, typeId)
		if !exists {
			cachedIds = make(map[string]struct{})
		}

		for cachedId := range cachedIds {
			resp.Data[typeId] = append(resp.Data[typeId], map[string]any{
				"id":           cachedId,
				"__syncAction": fibery.REMOVE,
			})
		}

		idCache.RemoveSource(sk)
	}
}

func (t *QuickBooksType[T]) processQuery(entities []T) ([]map[string]any, error) {
	items := []map[string]any{}
	for _, entity := range entities {
		item, err := t.schemaGen(entity)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema: %w", t.Id(), err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (t *dependentBaseType[ST]) processQuery(sourceEntities []ST) ([]map[string]any, error) {
	items := []map[string]any{}
	for _, source := range sourceEntities {
		itemSlice, err := t.schemaGen(source)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to fibery schema: %w", t.Id(), err)
		}
		items = append(items, itemSlice...)
	}
	return items, nil
}

func (t *QuickBooksCDCType[T]) processCDC(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	return processCDC(cdc, t.entityId, t.entityStatus, t.schemaGen)
}

func (t *QuickBooksDualType[T]) processCDC(cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	return processCDC(cdc, t.entityId, t.entityStatus, t.schemaGen)
}

func (t *DependentCDCType[ST]) processCDC(cache *IdCache, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	return processCDCDep(cache, cdc, t.sourceId, t.sourceStatus, t.schemaGen, t.sourceMapper, t.Id(), t.SourceId())
}

func (t *DependentDualType[ST]) processCDC(cache *IdCache, cdc *quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	return processCDCDep(cache, cdc, t.sourceId, t.sourceStatus, t.schemaGen, t.sourceMapper, t.Id(), t.SourceId())
}

func (t *QuickBooksWHType[T]) ProcessWebhookUpdate(batchData *quickbooks.BatchItemResponse, resp *fibery.WebhookTransformResponse) error {
	return processWebhookUpdates(batchData, resp, t.processQuery, t.Id())
}

func (t *QuickBooksDualType[T]) ProcessWebhookUpdate(batchData *quickbooks.BatchItemResponse, resp *fibery.WebhookTransformResponse) error {
	return processWebhookUpdates(batchData, resp, t.processQuery, t.Id())
}

func (t *DependentWHType[ST]) ProcessWebhookUpdate(batchData *quickbooks.BatchItemResponse, resp *fibery.WebhookTransformResponse, idCache *IdCache) error {
	return processWebhookUpdatesDep(batchData, resp, t.Id(), t.SourceId(), idCache, t.sourceId, t.schemaGen, t.sourceMapper)
}

func (t *DependentDualType[ST]) ProcessWebhookUpdate(batchData *quickbooks.BatchItemResponse, resp *fibery.WebhookTransformResponse, idCache *IdCache) error {
	return processWebhookUpdatesDep(batchData, resp, t.Id(), t.SourceId(), idCache, t.sourceId, t.schemaGen, t.sourceMapper)
}

func (t *QuickBooksWHType[T]) ProcessWebhookDeletions(ids []string, resp *fibery.WebhookTransformResponse) {
	processWebhookDeletions(ids, resp, t.Id())
}

func (t *QuickBooksDualType[T]) ProcessWebhookDeletions(ids []string, resp *fibery.WebhookTransformResponse) {
	processWebhookDeletions(ids, resp, t.Id())
}

func (t *DependentWHType[ST]) ProcessWebhookDeletions(sourceIds []string, resp *fibery.WebhookTransformResponse, idCache *IdCache) {
	processWebhookDeletionsDep(sourceIds, resp, t.Id(), t.SourceId(), idCache)
}

func (t *DependentDualType[ST]) ProcessWebhookDeletions(sourceIds []string, resp *fibery.WebhookTransformResponse, idCache *IdCache) {
	processWebhookDeletionsDep(sourceIds, resp, t.Id(), t.SourceId(), idCache)
}

func (t *QuickBooksWHType[T]) GetRelatedTypes() []CDCType {
	return t.relatedTypes
}

func (t *QuickBooksDualType[T]) GetRelatedTypes() []CDCType {
	return t.relatedTypes
}

func (t *DependentDataType[ST]) SourceId() string {
	return t.sourceType.Id()
}

func (t *DependentCDCType[ST]) SourceId() string {
	return t.sourceType.Id()
}

func (t *DependentWHType[ST]) SourceId() string {
	return t.sourceType.Id()
}

func (t *DependentDualType[ST]) SourceId() string {
	return t.sourceType.Id()
}

func (t *UnionDataType) SourceIds() []string {
	sourceIdMap := make(map[string]struct{})
	for _, unionType := range t.unionTypes {
		switch typ := unionType.(type) {
		case DeltaDepType:
			sourceIdMap[typ.SourceId()] = struct{}{}
		default:
			sourceIdMap[typ.Id()] = struct{}{}
		}
	}
	sourceIds := make([]string, 0, len(sourceIdMap))
	for id := range sourceIdMap {
		sourceIds = append(sourceIds, id)
	}
	return sourceIds
}

func (t *UnionDataType) UnionTypes() []Type {
	return t.unionTypes
}

func (t *DependentCDCType[ST]) SourceKey(data any) (IdKey, error) {
	sourceEntity, ok := data.(ST)
	if !ok {
		return IdKey{}, fmt.Errorf("unable to assert sourceType on data")
	}
	return IdKey{EntityType: t.sourceType.Id(), EntityId: t.sourceId(sourceEntity)}, nil
}

func (t *DependentWHType[ST]) SourceKey(data any) (IdKey, error) {
	sourceEntity, ok := data.(ST)
	if !ok {
		return IdKey{}, fmt.Errorf("unable to assert sourceType on data")
	}
	return IdKey{EntityType: t.sourceType.Id(), EntityId: t.sourceId(sourceEntity)}, nil
}

func (t *DependentDualType[ST]) SourceKey(data any) (IdKey, error) {
	sourceEntity, ok := data.(ST)
	if !ok {
		return IdKey{}, fmt.Errorf("unable to assert sourceType on data")
	}
	return IdKey{EntityType: t.sourceType.Id(), EntityId: t.sourceId(sourceEntity)}, nil
}

func (t *DependentCDCType[ST]) MapSource(data any) (map[string]struct{}, error) {
	sourceEntity, ok := data.(ST)
	if !ok {
		return nil, fmt.Errorf("unable to assert sourceType on data")
	}
	return t.sourceMapper(sourceEntity), nil
}

func (t *DependentWHType[ST]) MapSource(data any) (map[string]struct{}, error) {
	sourceEntity, ok := data.(ST)
	if !ok {
		return nil, fmt.Errorf("unable to assert sourceType on data")
	}
	return t.sourceMapper(sourceEntity), nil
}

func (t *DependentDualType[ST]) MapSource(data any) (map[string]struct{}, error) {
	sourceEntity, ok := data.(ST)
	if !ok {
		return nil, fmt.Errorf("unable to assert sourceType on data")
	}
	return t.sourceMapper(sourceEntity), nil
}

func (t *QuickBooksType[T]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.pageQuery, t.processQuery)
}

func (t *QuickBooksCDCType[T]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.pageQuery, t.processQuery)
}

func (t *QuickBooksWHType[T]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.pageQuery, t.processQuery)
}

func (t *QuickBooksDualType[T]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.pageQuery, t.processQuery)
}

func (t *DependentDataType[ST]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.sourceType.pageQuery, t.processQuery)
}

func (t *DependentCDCType[ST]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.sourceType.pageQuery, t.processQuery)
}

func (t *DependentWHType[ST]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.sourceType.pageQuery, t.processQuery)
}

func (t *DependentDualType[ST]) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	return getData(client, op, t, pagination.StartPosition, pageSize, t.sourceType.pageQuery, t.processQuery)
}

func (t *UnionDataType) GetData(client *quickbooks.Client, op *Operation, pagination fibery.NextPageConfig, pageSize int) (fibery.DataHandlerResponse, error) {
	response := fibery.DataHandlerResponse{
		SynchronizationType: fibery.Full,
	}
	var queryTypes []Type
	if len(pagination.Types) > 0 && pagination.StartPosition > pageSize {
		for _, pageType := range pagination.Types {
			for _, unionType := range t.unionTypes {
				if unionType.Id() == pageType {
					queryTypes = append(queryTypes, unionType)
				}
			}
		}
	} else {
		queryTypes = t.unionTypes
	}
	for _, unionType := range queryTypes {
		typeResponse, err := unionType.GetData(client, op, pagination, pageSize)
		if err != nil {
			return fibery.DataHandlerResponse{}, fmt.Errorf("unable to fetch data for union type %s:%s: %w", t.Id(), unionType.Id(), err)
		}
		if typeResponse.Pagination.HasNext {
			response.Pagination.HasNext = true
			response.Pagination.NextPageConfig.StartPosition = typeResponse.Pagination.NextPageConfig.StartPosition
			response.Pagination.NextPageConfig.Types = append(response.Pagination.NextPageConfig.Types, unionType.Id())
		}
		if typeResponse.SynchronizationType == fibery.Delta {
			response.SynchronizationType = fibery.Delta
		}
		typeItems, err := t.schemaGen(unionType.Id(), typeResponse.Items)
		if err != nil {
			return fibery.DataHandlerResponse{}, fmt.Errorf("unable to generate schema for union type %s:%s: %w", t.Id(), unionType.Id(), err)
		}
		response.Items = append(response.Items, typeItems...)
	}
	return response, nil
}
