package main

import (
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type Type interface {
	SourceId() string
	FiberyId() string
	FiberyName() string
	Schema() map[string]fibery.Field
	ProcessBatchQuery(batch *quickbooks.BatchItemResponse, idCache *IdCache, attachables []quickbooks.Attachable) (fibery.DataHandlerResponse, error)
}

type CDC interface {
	ProcessCDCQuery(cdc *quickbooks.ChangeDataCapture, idCache *IdCache, attachables []quickbooks.Attachable) (fibery.DataHandlerResponse, error)
}

type Webhook interface {
	ProcessWebhookDeletions(ids []string, idCache *IdCache) []map[string]any
}

type StandardData[T any] struct {
	Item        T
	Attachables []quickbooks.Attachable
}

type DependentData[ST, T any] struct {
	SourceItem  ST
	Item        T
	Attachables []quickbooks.Attachable
}

type FieldDef[T any] struct {
	id      string
	params  fibery.Field
	process func(StandardData[T]) (any, error)
}

type DependentFieldDef[ST, T any] struct {
	id      string
	params  fibery.Field
	process func(DependentData[ST, T]) (any, error)
}

type itemFieldValueFunc[T, V any] func(T) V
type sourceMapperFunc[T any] func(T) map[string]struct{}
type batchItemExtractorFunc[T any] func(quickbooks.BatchItemResponse) T
type batchQueryExtractorFunc[T any] func(quickbooks.BatchQueryResponse) []T
type cdcQueryExtractorFunc[T any] func(quickbooks.CDCQueryResponse) []T

type StandardType[T any] struct {
	sourceId            string
	fiberyId            string
	fiberyName          string
	fields              []FieldDef[T]
	batchItemExtractor  batchItemExtractorFunc[T]
	batchQueryExtractor batchQueryExtractorFunc[T]
}

func (t *StandardType[T]) Schema() map[string]fibery.Field {
	schema := make(map[string]fibery.Field, len(t.fields))
	for _, field := range t.fields {
		schema[field.id] = field.params
	}
	return schema
}

func (t *StandardType[T]) Process(data StandardData[T]) (map[string]any, error) {
	output := make(map[string]any, len(t.fields))
	for _, field := range t.fields {
		fieldValue, err := field.process(data)
		if err != nil {
			return nil, err
		}
		output[field.id] = fieldValue
	}
	return output, nil
}

func (t *StandardType[T]) ExtractBatchQuery(batch *quickbooks.BatchItemResponse) []T {
	return quickbooks.BatchQueryExtractor[T](batch, t.batchQueryExtractor)
}

type CDCType[T any] struct {
	StandardType[T]
	itemId            itemFieldValueFunc[T, string]
	itemStatus        itemFieldValueFunc[T, string]
	cdcQueryExtractor cdcQueryExtractorFunc[T]
}

// Schema() is embedded from StandardType[T]

func (t *CDCType[T]) Process(data StandardData[T]) (map[string]any, error) {
	if t.itemStatus(data.Item) == "Deleted" {
		return map[string]any{
			"id":           t.itemId(data.Item),
			"__syncAction": fibery.REMOVE,
		}, nil
	} else {
		output := make(map[string]any, len(t.fields))
		for _, field := range t.fields {
			fieldValue, err := field.process(data)
			if err != nil {
				return nil, err
			}
			output[field.id] = fieldValue
		}
		return output, nil
	}
}

// ExtractBatchQuery() is embedded from StandardType[T]

func (t *CDCType[T]) ExtractCDCQuery(cdc *quickbooks.ChangeDataCapture) []T {
	return quickbooks.CDCQueryExtractor[T](cdc, t.cdcQueryExtractor)
}

type WebhookType[T any] struct {
	StandardType[T]
}

type DualType[T any] struct {
	CDCType[T]
}

type DependentType[ST, T any] struct {
	sourceId            string
	fiberyId            string
	fiberyName          string
	fields              []DependentFieldDef[ST, T]
	batchItemExtractor  batchItemExtractorFunc[ST]
	batchQueryExtractor batchQueryExtractorFunc[ST]
}

type DependentCDCType[ST, T any] struct {
	DependentType[ST, T]
	itemId            itemFieldValueFunc[ST, string]
	itemStatus        itemFieldValueFunc[ST, string]
	cdcQueryExtractor cdcQueryExtractorFunc[ST]
	sourceMap         sourceMapperFunc[ST]
}

type DependentWebhookType[ST, T any] struct {
	DependentType[ST, T]
	sourceMap sourceMapperFunc[T]
}

type DependentDualType[ST, T any] struct {
	DependentCDCType[ST, T]
}

type UnionType struct{}
