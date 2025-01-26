package data

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/fibery"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/qbo"
	"golang.org/x/sync/singleflight"
)

type Request struct {
	StartPosition  int
	OperationId    string
	RealmId        string
	LastSynced     time.Time
	RequestedType  string
	RequestedTypes []string
	CDCTypes       []string
	Filter         map[string]any
	Cache          *cache.Cache
	Group          *singleflight.Group
	Token          *qbo.BearerToken
}

type Response struct {
	Data     any
	MoreData bool
}

type IdCache struct {
	Mu          sync.Mutex
	OperationId string
	Entries     map[string]map[string]bool
}

const IdCacheLifetime = 4 * time.Hour

type Type interface {
	GetId() string
	GetName() string
	GetSchema() map[string]fibery.Field
	Query(req Request) (Response, error)
	ProcessQuery(array any) ([]map[string]any, error)
}

type DependentType interface {
	Type
	GetSourceId() string
}

type CDCQueryable interface {
	Type
	ProcessCDC(cdc qbo.ChangeDataCapture) ([]map[string]any, error)
}

type DepCDCQueryable interface {
	DependentType
	MapType(sourceArray any) (map[string]map[string]bool, error)
	ProcessCDC(cdc qbo.ChangeDataCapture, cacheEntry *IdCache) ([]map[string]any, error)
}

type WHQueryable interface {
	Type
	ProcessWHBatch(itemResponse qbo.BatchItemResponse, request *map[string][]map[string]any) error
	ProcessWHDelete(deleteIds []string, request *map[string][]map[string]any, cache *cache.Cache, realmId string) error
}

type FiberyType struct {
	id     string
	name   string
	schema map[string]fibery.Field
}

func (f FiberyType) GetId() string {
	return f.id
}

func (f FiberyType) GetName() string {
	return f.name
}

func (f FiberyType) GetSchema() map[string]fibery.Field {
	return f.schema
}

type SchemaGenFunc func(entity any) (map[string]any, error)

type QueryProcessorFunc func(entityArray any, schemaGen SchemaGenFunc) ([]map[string]any, error)

type QuickbooksType struct {
	FiberyType
	schemaGen      SchemaGenFunc
	dataQuery      func(req Request) (Response, error)
	queryProcessor QueryProcessorFunc
}

func (t QuickbooksType) Query(req Request) (Response, error) {
	return t.dataQuery(req)
}

func (t QuickbooksType) ProcessQuery(array any) ([]map[string]any, error) {
	return t.queryProcessor(array, t.schemaGen)
}

type CDCProcessorFunc func(cdc qbo.ChangeDataCapture, schemaGen SchemaGenFunc) ([]map[string]any, error)

type QuickbooksCDCType struct {
	QuickbooksType
	changeDataCaptureProcessor CDCProcessorFunc
}

func (t QuickbooksCDCType) ProcessCDC(cdc qbo.ChangeDataCapture) ([]map[string]any, error) {
	return t.changeDataCaptureProcessor(cdc, t.schemaGen)
}

type WHQueryProcessorFunc func(itemResponse qbo.BatchItemResponse, response *map[string][]map[string]any, queryProcessor QueryProcessorFunc, schemaGen SchemaGenFunc, id string) error
type whDeleteProcessorFunc func(deleteIds []string, response *map[string][]map[string]any, cache *cache.Cache, realmId string, id string) error

type QuickbooksWHType struct {
	QuickbooksType
	whQueryProcessor  WHQueryProcessorFunc
	whDeleteProcessor whDeleteProcessorFunc
}

func (t QuickbooksWHType) ProcessWHBatch(itemResponse qbo.BatchItemResponse, response *map[string][]map[string]any) error {
	return t.whQueryProcessor(itemResponse, response, t.queryProcessor, t.schemaGen, t.GetId())
}

func (t QuickbooksWHType) ProcessWHDelete(deleteIds []string, response *map[string][]map[string]any, cache *cache.Cache, realmId string) error {
	return t.whDeleteProcessor(deleteIds, response, cache, realmId, t.GetId())
}

type QuickbooksDualType struct {
	QuickbooksType
	changeDataCaptureProcessor CDCProcessorFunc
	whQueryProcessor           WHQueryProcessorFunc
	whDeleteProcessor          whDeleteProcessorFunc
}

func (t QuickbooksDualType) ProcessCDC(cdc qbo.ChangeDataCapture) ([]map[string]any, error) {
	return t.changeDataCaptureProcessor(cdc, t.schemaGen)
}

func (t QuickbooksDualType) ProcessWHBatch(itemResponse qbo.BatchItemResponse, response *map[string][]map[string]any) error {
	return t.whQueryProcessor(itemResponse, response, t.queryProcessor, t.schemaGen, t.GetId())
}

func (t QuickbooksDualType) ProcessWHDelete(deleteIds []string, response *map[string][]map[string]any, cache *cache.Cache, realmId string) error {
	return t.whDeleteProcessor(deleteIds, response, cache, realmId, t.GetId())
}

type DependentSchemaGenFunc func(entity any, source any) (map[string]any, error)

type DependentBaseType struct {
	FiberyType
	schemaGen      DependentSchemaGenFunc
	queryProcessor func(sourceArray any, schemaGen DependentSchemaGenFunc) ([]map[string]any, error)
}

func (t DependentBaseType) ProcessQuery(array any) ([]map[string]any, error) {
	return t.queryProcessor(array, t.schemaGen)
}

type DependentDataType struct {
	DependentBaseType
	source QuickbooksType
}

func (t DependentDataType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentDataType) GetSourceId() string {
	return t.source.GetId()
}

type SourceMapperFunc func(source any) (map[string]bool, error)
type TypeMapperFunc func(sourceArray any, sourceMapper SourceMapperFunc) (map[string]map[string]bool, error)
type DependentCDCProcessorFunc func(cdc qbo.ChangeDataCapture, cacheEntry *IdCache, sourceMapper SourceMapperFunc, schemaGen DependentSchemaGenFunc) ([]map[string]any, error)

type DependentCDCType struct {
	DependentBaseType
	source                     QuickbooksCDCType
	sourceMapper               SourceMapperFunc
	typeMapper                 TypeMapperFunc
	changeDataCaptureProcessor DependentCDCProcessorFunc
}

func (t DependentCDCType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentCDCType) GetSourceId() string {
	return t.source.GetId()
}

func (t DependentCDCType) ProcessCDC(cdc qbo.ChangeDataCapture, idEntry *IdCache) ([]map[string]any, error) {
	return t.changeDataCaptureProcessor(cdc, idEntry, t.sourceMapper, t.schemaGen)
}

func (t DependentCDCType) MapType(sourceArray any) (map[string]map[string]bool, error) {
	return t.typeMapper(sourceArray, t.sourceMapper)
}

type DependentWHType struct {
	DependentBaseType
	source QuickbooksWHType
}

func (t DependentWHType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentWHType) GetSourceId() string {
	return t.source.GetId()
}

type DependentDualType struct {
	DependentBaseType
	source                     QuickbooksDualType
	sourceMapper               SourceMapperFunc
	typeMapper                 TypeMapperFunc
	changeDataCaptureProcessor DependentCDCProcessorFunc
}

func (t DependentDualType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentDualType) GetSourceId() string {
	return t.source.GetId()
}

var Types = map[string]*Type{}
var SourceDependents = map[string][]*DependentType{}

func RegisterType(t Type) {
	Types[t.GetId()] = &t
	deptype, ok := t.(DependentType)
	if ok {
		SourceDependents[deptype.GetSourceId()] = append(SourceDependents[deptype.GetSourceId()], &deptype)
	}
}

func FormatJSON(data interface{}) string {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Failed to generate json", err)
		return ""
	}
	return string(prettyJSON)
}
