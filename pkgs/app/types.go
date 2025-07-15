package app

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type TypeRegistry map[string]fibery.Type

func (tr TypeRegistry) Register(t fibery.Type) {
	tr[t.Id()] = t
}

func (tr TypeRegistry) Get(id string) (fibery.Type, bool) {
	if typ, exists := tr[id]; exists {
		return typ, true
	}
	return nil, false
}

func (tr TypeRegistry) GetAll() []fibery.SyncConfigTypes {
	types := make([]fibery.SyncConfigTypes, 0, len(tr))
	for _, typ := range tr {
		types = append(types, fibery.SyncConfigTypes{
			Id:   typ.Id(),
			Name: typ.Name(),
		})
	}
	return types
}

var Types = make(TypeRegistry)

type StandardType interface {
	fibery.Type
	Type() string
	Attachables(fieldId string) bool
	ProcessBatchQuery(batch *quickbooks.BatchItemResponse, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, bool, error)
}

type StandardDependentType interface {
	fibery.Type
	SourceType() string
	ProcessBatchQuery(batch *quickbooks.BatchItemResponse, idCache *IdCache, pageSize int) ([]map[string]any, bool, error)
}

type CDCType interface {
	StandardType
	ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, error)
}

type CDCDependentType interface {
	StandardDependentType
	ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, idCache *IdCache, pageSize int) ([]map[string]any, error)
}

type WebhookType interface {
	StandardType
	ProcessWebhookDeletions(ids []string) ([]map[string]any, error)
	RelatedTypes() map[string]CDCType
}

type WebhookDependentType interface {
	StandardDependentType
	ProcessWebhookDeletions(sourceIds []string, idCache *IdCache) ([]map[string]any, error)
}

type DualType interface {
	CDCType
	WebhookType
}

type UnionType interface {
	fibery.Type
	Types() []StandardType
	CDC() bool
	Webhook() bool
	ProcessBatchQuery(batches map[string]*quickbooks.BatchItemResponse, pageSize int) ([]map[string]any, map[string]struct{}, error)
	ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, pageSize int) ([]map[string]any, error)
	ProcessWebhookDeletions(deletedSources map[string][]string) ([]map[string]any, error)
}

type StaticType interface {
	fibery.Type
	GetData() []map[string]any
}

type StandardData[T any] struct {
	Item        T
	Attachables map[string][]quickbooks.Attachable
}

type DependentData[ST, T any] struct {
	SourceItem ST
	Item       T
}

type FieldDef[T any] struct {
	Params  fibery.Field
	Convert func(StandardData[T]) (any, error)
}

type DependentFieldDef[ST, T any] struct {
	Params  fibery.Field
	Convert func(DependentData[ST, T]) (any, error)
}

type UnionFieldDef struct {
	Params  fibery.Field
	Convert func(string, map[string]any) (any, error)
}

type StandardTypeDef[T any] struct {
	TypeId              string
	FiberyId            string
	FiberyName          string
	Fields              map[string]FieldDef[T]
	BatchItemExtractor  func(quickbooks.BatchItemResponse) T
	BatchQueryExtractor func(quickbooks.BatchQueryResponse) []T
}

type CDCTypeDef[T any] struct {
	StandardTypeDef[T]
	ItemId            func(T) string
	ItemStatus        func(T) string
	CDCQueryExtractor func(quickbooks.CDCQueryResponse) []T
}

type WebhookTypeDef[T any] struct {
	StandardTypeDef[T]
	ItemId      func(T) string
	ItemBuilder func(id string) T
	Related     []CDCType
}

type DualTypeDef[T any] struct {
	CDCTypeDef[T]
	ItemBuilder func(id string) T
	Related     []CDCType
}

type DependentTypeDef[ST, T any] struct {
	SourceTypeId        string
	FiberyId            string
	FiberyName          string
	ItemId              func(ST, T) string
	ItemCheck           func(ST, T) bool
	ItemExtractor       func(ST) []T
	Fields              map[string]DependentFieldDef[ST, T]
	BatchItemExtractor  func(quickbooks.BatchItemResponse) ST
	BatchQueryExtractor func(quickbooks.BatchQueryResponse) []ST
}

type DependentCDCTypeDef[ST, T any] struct {
	DependentTypeDef[ST, T]
	SourceId          func(ST) string
	SourceStatus      func(ST) string
	CDCQueryExtractor func(quickbooks.CDCQueryResponse) []ST
}

type DependentWebhookTypeDef[ST, T any] struct {
	DependentTypeDef[ST, T]
	SourceId      func(ST) string
	SourceBuilder func(id string) ST
}

type DependentDualTypeDef[ST, T any] struct {
	DependentCDCTypeDef[ST, T]
	SourceBuilder func(id string) ST
}

type UnionTypeDef struct {
	SourceTypes []StandardType
	FiberyId    string
	FiberyName  string
	Fields      map[string]UnionFieldDef
}

type StaticTypeDef struct {
	FiberyId   string
	FiberyName string
	Fields     map[string]fibery.Field
	Data       []map[string]any
}

func NewStandardType[T any](
	typeId, fiberyId, fiberyName string,
	itemId func(T) string,
	batchItemExtractor func(quickbooks.BatchItemResponse) T,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []T,
	addlFields map[string]FieldDef[T],
) *StandardTypeDef[T] {
	fields := make(map[string]FieldDef[T], len(addlFields)+1)
	for k, v := range addlFields {
		fields[k] = v
	}
	fields["id"] = FieldDef[T]{
		Params: fibery.Field{
			Name: "Id",
			Type: fibery.Id,
		},
		Convert: func(sd StandardData[T]) (any, error) {
			return itemId(sd.Item), nil
		},
	}

	return &StandardTypeDef[T]{
		TypeId:              typeId,
		FiberyId:            fiberyId,
		FiberyName:          fiberyName,
		Fields:              fields,
		BatchItemExtractor:  batchItemExtractor,
		BatchQueryExtractor: batchQueryExtractor,
	}
}

func NewCDCType[T any](
	typeId, fiberyId, fiberyName string,
	itemId func(T) string,
	itemStatus func(T) string,
	batchItemExtractor func(quickbooks.BatchItemResponse) T,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []T,
	cdcQueryExtractor func(quickbooks.CDCQueryResponse) []T,
	addlFields map[string]FieldDef[T],
) *CDCTypeDef[T] {
	return &CDCTypeDef[T]{
		StandardTypeDef: *NewStandardType(
			typeId,
			fiberyId,
			fiberyName,
			itemId,
			batchItemExtractor,
			batchQueryExtractor,
			addlFields,
		),
		ItemId:            itemId,
		ItemStatus:        itemStatus,
		CDCQueryExtractor: cdcQueryExtractor,
	}
}

func NewWebhookType[T any](
	typeId, fiberyId, fiberyName string,
	itemId func(T) string,
	itemStatus func(T) string,
	itemBuilder func(id string) T,
	batchItemExtractor func(quickbooks.BatchItemResponse) T,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []T,
	addlFields map[string]FieldDef[T],
	relatedTypes []CDCType,
) *WebhookTypeDef[T] {
	return &WebhookTypeDef[T]{
		StandardTypeDef: *NewStandardType(
			typeId,
			fiberyId,
			fiberyName,
			itemId,
			batchItemExtractor,
			batchQueryExtractor,
			addlFields,
		),
		ItemId:      itemId,
		ItemBuilder: itemBuilder,
		Related:     relatedTypes,
	}
}

func NewDualType[T any](
	typeId, fiberyId, fiberyName string,
	itemId func(T) string,
	itemStatus func(T) string,
	itemBuilder func(id string) T,
	batchItemExtractor func(quickbooks.BatchItemResponse) T,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []T,
	cdcQueryExtractor func(quickbooks.CDCQueryResponse) []T,
	addlFields map[string]FieldDef[T],
	relatedTypes []CDCType,
) *DualTypeDef[T] {
	return &DualTypeDef[T]{
		CDCTypeDef: CDCTypeDef[T]{
			StandardTypeDef: *NewStandardType(
				typeId,
				fiberyId,
				fiberyName,
				itemId,
				batchItemExtractor,
				batchQueryExtractor,
				addlFields,
			),
			ItemId:            itemId,
			ItemStatus:        itemStatus,
			CDCQueryExtractor: cdcQueryExtractor,
		},
		ItemBuilder: itemBuilder,
		Related:     relatedTypes,
	}
}

func NewDependentType[ST, T any](
	sourceTypeId, fiberyId, fiberyName string,
	itemId func(ST, T) string,
	itemCheck func(ST, T) bool,
	itemExtractor func(ST) []T,
	batchItemExtractor func(quickbooks.BatchItemResponse) ST,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []ST,
	addlFields map[string]DependentFieldDef[ST, T],
) *DependentTypeDef[ST, T] {
	fields := make(map[string]DependentFieldDef[ST, T], len(addlFields)+1)
	for k, v := range addlFields {
		fields[k] = v
	}
	fields["id"] = DependentFieldDef[ST, T]{
		Params: fibery.Field{
			Name: "Id",
			Type: fibery.Id,
		},
		Convert: func(dd DependentData[ST, T]) (any, error) {
			return itemId(dd.SourceItem, dd.Item), nil
		},
	}

	return &DependentTypeDef[ST, T]{
		SourceTypeId:        sourceTypeId,
		FiberyId:            fiberyId,
		FiberyName:          fiberyName,
		ItemId:              itemId,
		ItemCheck:           itemCheck,
		ItemExtractor:       itemExtractor,
		Fields:              fields,
		BatchItemExtractor:  batchItemExtractor,
		BatchQueryExtractor: batchQueryExtractor,
	}
}

func NewDependentCDCType[ST, T any](
	sourceTypeId, fiberyId, fiberyName string,
	itemId func(ST, T) string,
	itemCheck func(ST, T) bool,
	itemExtractor func(ST) []T,
	sourceId func(ST) string,
	sourceStatus func(ST) string,
	batchItemExtractor func(quickbooks.BatchItemResponse) ST,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []ST,
	cdcQueryExtractor func(quickbooks.CDCQueryResponse) []ST,
	addlFields map[string]DependentFieldDef[ST, T],
) *DependentCDCTypeDef[ST, T] {
	return &DependentCDCTypeDef[ST, T]{
		DependentTypeDef: *NewDependentType(
			sourceTypeId,
			fiberyId, fiberyName,
			itemId,
			itemCheck,
			itemExtractor,
			batchItemExtractor,
			batchQueryExtractor,
			addlFields,
		),
		SourceId:          sourceId,
		SourceStatus:      sourceStatus,
		CDCQueryExtractor: cdcQueryExtractor,
	}
}

func NewDependentWebhookType[ST, T any](
	sourceTypeId, fiberyId, fiberyName string,
	itemId func(ST, T) string,
	itemCheck func(ST, T) bool,
	itemExtractor func(ST) []T,
	sourceId func(ST) string,
	sourceBuilder func(id string) ST,
	batchItemExtractor func(quickbooks.BatchItemResponse) ST,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []ST,
	cdcQueryExtractor func(quickbooks.CDCQueryResponse) []ST,
	addlFields map[string]DependentFieldDef[ST, T],
) *DependentWebhookTypeDef[ST, T] {
	return &DependentWebhookTypeDef[ST, T]{
		DependentTypeDef: *NewDependentType(
			sourceTypeId,
			fiberyId, fiberyName,
			itemId,
			itemCheck,
			itemExtractor,
			batchItemExtractor,
			batchQueryExtractor,
			addlFields,
		),
		SourceId:      sourceId,
		SourceBuilder: sourceBuilder,
	}
}

func NewDependentDualType[ST, T any](
	sourceTypeId, fiberyId, fiberyName string,
	itemId func(ST, T) string,
	itemCheck func(ST, T) bool,
	itemExtractor func(ST) []T,
	sourceId func(ST) string,
	sourceStatus func(ST) string,
	sourceBuilder func(id string) ST,
	batchItemExtractor func(quickbooks.BatchItemResponse) ST,
	batchQueryExtractor func(quickbooks.BatchQueryResponse) []ST,
	cdcQueryExtractor func(quickbooks.CDCQueryResponse) []ST,
	addlFields map[string]DependentFieldDef[ST, T],
) *DependentDualTypeDef[ST, T] {
	return &DependentDualTypeDef[ST, T]{
		DependentCDCTypeDef: DependentCDCTypeDef[ST, T]{
			DependentTypeDef: *NewDependentType(
				sourceTypeId,
				fiberyId, fiberyName,
				itemId,
				itemCheck,
				itemExtractor,
				batchItemExtractor,
				batchQueryExtractor,
				addlFields,
			),
			SourceId:          sourceId,
			SourceStatus:      sourceStatus,
			CDCQueryExtractor: cdcQueryExtractor,
		},
		SourceBuilder: sourceBuilder,
	}
}

func NewUnionType(
	sourceTypes []StandardType,
	fiberyId, fiberyName string,
	fields map[string]UnionFieldDef,
) *UnionTypeDef {
	return &UnionTypeDef{
		SourceTypes: sourceTypes,
		FiberyId:    fiberyId,
		FiberyName:  fiberyName,
		Fields:      fields,
	}
}

func NewStaticType(
	fiberyId, fiberyName string,
	fields map[string]fibery.Field,
	data []map[string]any,
) *StaticTypeDef {
	return &StaticTypeDef{
		FiberyId:   fiberyId,
		FiberyName: fiberyName,
		Fields:     fields,
		Data:       data,
	}
}

// --- StandardTypeDef[T] methods ---

func (t *StandardTypeDef[T]) Id() string {
	return t.FiberyId
}

func (t *StandardTypeDef[T]) Name() string {
	return t.FiberyName
}

func (t *StandardTypeDef[T]) Schema() map[string]fibery.Field {
	schema := make(map[string]fibery.Field, len(t.Fields))
	for id, field := range t.Fields {
		schema[id] = field.Params
	}
	return schema
}

func (t *StandardTypeDef[T]) Type() string {
	return t.TypeId
}

func (t *StandardTypeDef[T]) Attachables(fieldId string) bool {
	_, ok := t.Fields[fieldId]
	return ok
}

func (t *StandardTypeDef[T]) Convert(data StandardData[T]) (map[string]any, error) {
	output := make(map[string]any, len(t.Fields))
	for id, field := range t.Fields {
		fieldValue, err := field.Convert(data)
		if err != nil {
			return nil, err
		}
		output[id] = fieldValue
	}
	return output, nil
}

func (t *StandardTypeDef[T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T {
	return quickbooks.BatchQueryExtractor(batch, t.BatchQueryExtractor)
}

func (t *StandardTypeDef[T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, bool, error) {
	input := t.extractBatchQuery(batch)

	more := len(input) == pageSize

	output := make([]map[string]any, 0, len(input))

	for _, i := range input {
		data := StandardData[T]{
			Item:        i,
			Attachables: attachables,
		}

		o, err := t.Convert(data)
		if err != nil {
			return nil, more, fmt.Errorf("error converting input data: %w", err)
		}

		output = append(output, o)
	}

	return output, more, nil
}

// --- CDCTypeDef[T] methods ---

// methods embedded from StandardTypeDef[T]:

// func (t *CDCTypeDef[T]) Id() string

// func (t *CDCTypeDef[T]) Name() string

// func (t *CDCTypeDef[T]) Schema() map[string]fibery.field

// func (t *CDCTypeDef[T]) TypeId() string

// func (t *CDCTypeDef[T]) Attachables(fieldId string) bool

func (t *CDCTypeDef[T]) Convert(data StandardData[T]) (map[string]any, error) {
	output := make(map[string]any, len(t.Fields))
	for id, field := range t.Fields {
		fieldValue, err := field.Convert(data)
		if err != nil {
			return nil, err
		}
		output[id] = fieldValue
	}
	return output, nil
}

// func (t *CDCTypeDef[T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T

// func (t *CDCTypeDef[T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, bool, error)

func (t *CDCTypeDef[T]) extractCDCQuery(cdc *quickbooks.ChangeDataCapture) []T {
	return quickbooks.CDCQueryExtractor(cdc, t.CDCQueryExtractor)
}

func (t *CDCTypeDef[T]) ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, error) {
	input := t.extractCDCQuery(cdc)

	if len(input) == pageSize {
		return nil, fmt.Errorf("cdc response for %s is equal to pageSize, please force full sync", t.Type())
	}

	output := make([]map[string]any, 0, len(input))

	for _, item := range input {
		if t.ItemStatus(item) == "Deleted" {
			o := map[string]any{
				"id":           t.ItemId(item),
				"__syncAction": fibery.REMOVE,
			}

			output = append(output, o)
		} else {
			data := StandardData[T]{
				Item:        item,
				Attachables: attachables,
			}

			o, err := t.Convert(data)
			if err != nil {
				return nil, fmt.Errorf("error converting input data: %w", err)
			}

			output = append(output, o)
		}
	}

	return output, nil
}

// --- methods for WebhookTypeDef[T] ---

// methods embedded from StandardTypeDef[T]:

// func (t *WebhookTypeDef[T]) Id() string

// func (t *WebhookTypeDef[T]) Name() string

// func (t *WebhookTypeDef[T]) Schema() map[string]fibery.field

// func (t *WebhookTypeDef[T]) TypeId() string

// func (t *WebhookTypeDef[T]) Attachables(fieldId string) bool

// func (t *WebhookTypeDef[T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T

// func (t *WebhookTypeDef[T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, bool, error)

func (t *WebhookTypeDef[T]) ProcessWebhookDeletions(ids []string) ([]map[string]any, error) {
	output := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		item := t.ItemBuilder(id)

		outputId := t.ItemId(item)

		output = append(output, map[string]any{
			"id":           outputId,
			"__syncAction": fibery.REMOVE,
		})
	}

	return output, nil
}

func (t *WebhookTypeDef[T]) RelatedTypes() map[string]CDCType {
	output := make(map[string]CDCType, len(t.Related))
	for _, typ := range t.Related {
		output[typ.Id()] = typ
	}
	return output
}

// --- methods for DualTypeDef[T] ---

// methods embedded from StandardTypeDef[T]:

// func (t *DualTypeDef[T]) Id() string

// func (t *DualTypeDef[T]) Name() string

// func (t *DualTypeDef[T]) Schema() map[string]fibery.field

// func (t *DualTypeDef[T]) TypeId() string

// func (t *DualTypeDef[T]) Attachables(fieldId string) bool

// func (t *DualTypeDef[T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T

// func (t *DualTypeDef[T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, bool, error)

// methods embedded from CDCTypeDef[T]:

// func (t *DualTypeDef[T]) extractCDCQuery(cdc *quickbooks.ChangeDataCapture) []T

// func (t *DualTypeDef[T]) ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, attachables map[string][]quickbooks.Attachable, pageSize int) ([]map[string]any, error)

func (t *DualTypeDef[T]) ProcessWebhookDeletions(ids []string) ([]map[string]any, error) {
	output := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		item := t.ItemBuilder(id)

		outputId := t.ItemId(item)

		output = append(output, map[string]any{
			"id":           outputId,
			"__syncAction": fibery.REMOVE,
		})
	}

	return output, nil
}

func (t *DualTypeDef[T]) RelatedTypes() map[string]CDCType {
	output := make(map[string]CDCType, len(t.Related))
	for _, typ := range t.Related {
		output[typ.Id()] = typ
	}
	return output
}

// --- methods for DependentTypeDef[ST, T] ---

func (t *DependentTypeDef[ST, T]) Id() string {
	return t.FiberyId
}

func (t *DependentTypeDef[ST, T]) Name() string {
	return t.FiberyName
}

func (t *DependentTypeDef[ST, T]) Schema() map[string]fibery.Field {
	schema := make(map[string]fibery.Field, len(t.Fields))
	for id, field := range t.Fields {
		schema[id] = field.Params
	}
	return schema
}

func (t *DependentTypeDef[ST, T]) SourceType() string {
	return t.SourceTypeId
}

func (t *DependentTypeDef[ST, T]) Convert(data DependentData[ST, T]) (map[string]any, error) {
	output := make(map[string]any, len(t.Fields))
	for id, field := range t.Fields {
		fieldValue, err := field.Convert(data)
		if err != nil {
			return nil, err
		}
		output[id] = fieldValue
	}
	return output, nil
}

func (t *DependentTypeDef[ST, T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []ST {
	return quickbooks.BatchQueryExtractor(batch, t.BatchQueryExtractor)
}

func (t *DependentTypeDef[ST, T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, idCache *IdCache, pageSize int) ([]map[string]any, bool, error) {
	input := t.extractBatchQuery(batch)

	more := len(input) == pageSize

	output := []map[string]any{}

	for _, source := range input {
		items := t.ItemExtractor(source)
		for _, i := range items {
			data := DependentData[ST, T]{
				SourceItem: source,
				Item:       i,
			}

			o, err := t.Convert(data)
			if err != nil {
				return nil, more, fmt.Errorf("error converting input data: %w", err)
			}

			output = append(output, o)
		}
	}

	return output, more, nil
}

// --- methods for DependentCDCTypeDef[ST, T] ---

// methods embedded from DependentTypeDef[ST, T]:

// func (t *DependentCDCTypeDef[ST, T]) Id() string

// func (t *DependentCDCTypeDef[ST, T]) Name() string

// func (t *DependentCDCTypeDef[ST, T]) Schema() map[string]fibery.field

// func (t *DependentCDCTypeDef[ST, T]) SourceTypeId() string

// func (t *DependentCDCTypeDef[ST, T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T

func (t *DependentCDCTypeDef[ST, T]) sourceKey(source ST) IdKey {
	return IdKey{EntityType: t.SourceType(), EntityId: t.SourceId(source)}
}

func (t *DependentCDCTypeDef[ST, T]) sourceMap(source ST) map[string]struct{} {
	items := t.ItemExtractor(source)

	output := make(map[string]struct{}, len(items))

	for _, i := range items {
		output[t.ItemId(source, i)] = struct{}{}
	}

	return output
}

func (t *DependentCDCTypeDef[ST, T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, idCache *IdCache, pageSize int) ([]map[string]any, bool, error) {
	input := t.extractBatchQuery(batch)

	more := len(input) == pageSize

	output := []map[string]any{}

	for _, source := range input {
		items := t.ItemExtractor(source)
		for _, i := range items {
			data := DependentData[ST, T]{
				SourceItem: source,
				Item:       i,
			}

			o, err := t.Convert(data)
			if err != nil {
				return nil, more, fmt.Errorf("error converting input data: %w", err)
			}

			idCache.SetIds(t.sourceKey(source), t.Id(), t.sourceMap(source))

			output = append(output, o)
		}
	}

	return output, more, nil
}

func (t *DependentCDCTypeDef[ST, T]) extractCDCQuery(cdc *quickbooks.ChangeDataCapture) []ST {
	return quickbooks.CDCQueryExtractor(cdc, t.CDCQueryExtractor)
}

func (t *DependentCDCTypeDef[ST, T]) ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, idCache *IdCache, pageSize int) ([]map[string]any, error) {
	input := t.extractCDCQuery(cdc)

	if len(input) == pageSize {
		return nil, fmt.Errorf("cdc response for %s is equal to pageSize, please force full sync", t.Id())
	}

	output := []map[string]any{}

	for _, source := range input {
		sourceKey := t.sourceKey(source)
		cachedIds, exists := idCache.GetIdsByType(sourceKey, t.Id())
		if !exists {
			cachedIds = make(map[string]struct{})
		}

		newMap := t.sourceMap(source)

		if t.SourceStatus(source) == "Deleted" {
			for cachedId := range cachedIds {
				output = append(output, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
			idCache.RemoveEntityType(sourceKey, t.Id())
			continue
		}

		items := t.ItemExtractor(source)
		for _, i := range items {
			data := DependentData[ST, T]{
				SourceItem: source,
				Item:       i,
			}

			o, err := t.Convert(data)
			if err != nil {
				return nil, fmt.Errorf("error converting input data: %w", err)
			}

			output = append(output, o)
		}

		for cachedId := range cachedIds {
			if _, ok := newMap[cachedId]; !ok {
				output = append(output, map[string]any{
					"id":           cachedId,
					"__syncAction": fibery.REMOVE,
				})
			}
		}

		idCache.AddIds(sourceKey, t.Id(), newMap)
	}

	return output, nil
}

// --- methods for DependentWebhookTypeDef[ST, T] ---

// methods embedded from DependentTypeDef[ST, T]:

// func (t *DependentWebhookTypeDef[ST, T]) Id() string

// func (t *DependentWebhookTypeDef[ST, T]) Name() string

// func (t *DependentWebhookTypeDef[ST, T]) Schema() map[string]fibery.field

// func (t *DependentWebhookTypeDef[ST, T]) SourceTypeId() string

// func (t *DependentWebhookTypeDef[ST, T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T

func (t *DependentWebhookTypeDef[ST, T]) sourceKey(source ST) IdKey {
	return IdKey{EntityType: t.SourceType(), EntityId: t.SourceId(source)}
}

func (t *DependentWebhookTypeDef[ST, T]) sourceMap(source ST) map[string]struct{} {
	items := t.ItemExtractor(source)

	output := make(map[string]struct{}, len(items))

	for _, i := range items {
		output[t.ItemId(source, i)] = struct{}{}
	}

	return output
}

func (t *DependentWebhookTypeDef[ST, T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, idCache *IdCache, pageSize int) ([]map[string]any, bool, error) {
	input := t.extractBatchQuery(batch)

	more := len(input) == pageSize

	output := []map[string]any{}

	for _, source := range input {
		items := t.ItemExtractor(source)
		for _, i := range items {
			data := DependentData[ST, T]{
				SourceItem: source,
				Item:       i,
			}

			o, err := t.Convert(data)
			if err != nil {
				return nil, more, fmt.Errorf("error converting input data: %w", err)
			}

			idCache.SetIds(t.sourceKey(source), t.Id(), t.sourceMap(source))

			output = append(output, o)
		}
	}

	return output, more, nil
}

func (t *DependentWebhookTypeDef[ST, T]) ProcessWebhookDeletions(sourceIds []string, idCache *IdCache) []map[string]any {
	output := []map[string]any{}
	for _, sourceId := range sourceIds {
		sourceKey := t.sourceKey(t.SourceBuilder(sourceId))

		cachedIds, exists := idCache.GetIdsByType(sourceKey, t.Id())
		if !exists {
			return nil
		}

		for cachedId := range cachedIds {
			output = append(output, map[string]any{
				"id":           cachedId,
				"__syncAction": fibery.REMOVE,
			})
		}

		idCache.RemoveEntityType(sourceKey, t.Id())
	}

	return output
}

// -- methods for DependentDualTypeDef[ST, T] ---

// methods embedded from DependentTypeDef[ST, T]:

// func (t *DependentDualTypeDef[ST, T]) Id() string

// func (t *DependentDualTypeDef[ST, T]) Name() string

// func (t *DependentDualTypeDef[ST, T]) Schema() map[string]fibery.field

// func (t *DependentDualTypeDef[ST, T]) SourceTypeId() string

// func (t *DependentDualTypeDef[ST, T]) extractBatchQuery(batch *quickbooks.BatchItemResponse) []T

// methods embedded from DependentCDCTypeDef[ST, T]:

// func (t *DependentDualTypeDef[ST, T]) sourceKey(source ST) IdKey

// func (t *DependentDualTypeDef[ST, T]) sourceMap(source ST) map[string]struct{}

// func (t *DependentDualTypeDef[ST, T]) ProcessBatchQuery(batch *quickbooks.BatchItemResponse, idCache *IdCache, pageSize int) ([]map[string]any, bool, error)

// func (t *DependentDualTypeDef[ST, T]) extractCDCQuery(cdc *quickbooks.ChangeDataCapture) []ST

// func (t *DependentDualTypeDef[ST, T]) ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, idCache *IdCache, pageSize int) ([]map[string]any, error)

func (t *DependentDualTypeDef[ST, T]) ProcessWebhookDeletions(sourceIds []string, idCache *IdCache) []map[string]any {
	output := []map[string]any{}
	for _, sourceId := range sourceIds {
		sourceKey := t.sourceKey(t.SourceBuilder(sourceId))

		cachedIds, exists := idCache.GetIdsByType(sourceKey, t.Id())
		if !exists {
			return nil
		}

		for cachedId := range cachedIds {
			output = append(output, map[string]any{
				"id":           cachedId,
				"__syncAction": fibery.REMOVE,
			})
		}

		idCache.RemoveEntityType(sourceKey, t.Id())
	}

	return output
}

// --- methods for UnionTypeDef ---

func (t *UnionTypeDef) Id() string {
	return t.FiberyId
}

func (t *UnionTypeDef) Name() string {
	return t.FiberyName
}

func (t *UnionTypeDef) Schema() map[string]fibery.Field {
	schema := make(map[string]fibery.Field, len(t.Fields))
	for id, field := range t.Fields {
		schema[id] = field.Params
	}
	return schema
}

func (t *UnionTypeDef) Convert(typeId string, item map[string]any) (map[string]any, error) {
	output := make(map[string]any, len(t.Fields))
	for id, field := range t.Fields {
		fieldValue, err := field.Convert(typeId, item)
		if err != nil {
			return nil, err
		}
		output[id] = fieldValue
	}
	return output, nil
}

func (t *UnionTypeDef) Types() []StandardType {
	return t.SourceTypes
}

func (t *UnionTypeDef) CDC() bool {
	for _, source := range t.SourceTypes {
		if _, ok := source.(CDCType); !ok {
			return false
		}
	}
	return true
}

func (t *UnionTypeDef) Webhook() bool {
	for _, source := range t.SourceTypes {
		if _, ok := source.(WebhookType); !ok {
			return false
		}
	}
	return true
}

func (t *UnionTypeDef) ProcessBatchQuery(batches map[string]*quickbooks.BatchItemResponse, pageSize int) ([]map[string]any, map[string]struct{}, error) {
	output := make([]map[string]any, 0)
	more := make(map[string]struct{})
	for _, source := range t.SourceTypes {
		if batch, ok := batches[source.Type()]; ok {
			input, m, err := source.ProcessBatchQuery(batch, nil, pageSize)
			if err != nil {
				return nil, nil, fmt.Errorf("error processing batch for %s: %w", source.Type(), err)
			}

			typeOutput := make([]map[string]any, 0, len(input))
			for _, item := range input {
				o, err := t.Convert(source.Type(), item)
				if err != nil {
					return nil, nil, fmt.Errorf("error converting %s: %w", source.Type(), err)
				}

				typeOutput = append(typeOutput, o)
			}

			output = append(output, typeOutput...)
			if m {
				more[source.Type()] = struct{}{}
			}
		}
	}
	return output, more, nil
}

func (t *UnionTypeDef) ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, pageSize int) ([]map[string]any, error) {
	output := []map[string]any{}
	for _, source := range t.SourceTypes {
		if cdcSource, ok := source.(CDCType); ok {
			input, err := cdcSource.ProcessCDCQuery(cdc, nil, pageSize)
			if err != nil {
				return nil, fmt.Errorf("error processing cdc for %s", cdcSource.Type())
			}

			typeOutput := make([]map[string]any, 0, len(input))
			for _, item := range input {
				o, err := t.Convert(source.Type(), item)
				if err != nil {
					return nil, fmt.Errorf("error converting %s: %w", source.Type(), err)
				}

				typeOutput = append(typeOutput, o)
			}
			output = append(output, typeOutput...)
		} else {
			return nil, fmt.Errorf("source %s is not a CDCType", source.Type())
		}
	}

	return output, nil
}

func (t *UnionTypeDef) ProcessWebhookDeletions(deletedSources map[string][]string) ([]map[string]any, error) {
	if field, ok := t.Fields["id"]; ok {
		output := make([]map[string]any, 0)
		for typeId, ids := range deletedSources {
			for _, id := range ids {
				idField := map[string]any{"id": id}
				outputId, err := field.Convert(typeId, idField)
				if err != nil {
					return nil, fmt.Errorf("error converting id value: %w", err)
				}

				output = append(output, map[string]any{
					"id":           outputId,
					"__syncAction": fibery.REMOVE,
				})
			}
		}

		return output, nil
	} else {
		return nil, fmt.Errorf("%s does not have an 'id' field", t.Id())
	}
}

func (t *StaticTypeDef) Id() string {
	return t.FiberyId
}

func (t *StaticTypeDef) Name() string {
	return t.FiberyName
}

func (t *StaticTypeDef) Schema() map[string]fibery.Field {
	return t.Fields
}

func (t *StaticTypeDef) GetData() []map[string]any {
	return t.Data
}
