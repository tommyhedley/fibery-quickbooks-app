package data

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
	"golang.org/x/sync/singleflight"
)

type Request struct {
	StartPosition     int
	OperationId       string
	LastSynced        time.Time
	RequestedType     string
	RequestedTypes    []string
	RequestedCDCTypes []string
	Filter            map[string]any
	Cache             *cache.Cache
	Group             *singleflight.Group
	Client            *quickbooks.Client
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

// These interfaces group datatypes by what queries can be performed on them

// Type defines the most basic requirements for a datatype to be synced between QuickBooks and Fibery
type Type interface {
	GetId() string
	GetName() string
	GetSchema() map[string]fibery.Field
	Query(req Request) (Response, error)
	ProcessQuery(array any) ([]map[string]any, error)
}

// DependentType limits Types that are not individually queryable in QuickBooks (ex. invoice lines) but need to be seperated into their own Fibery type for proper structure
type DependentType interface {
	Type
	GetSourceId() string
}

// CDCQueryable limits Types to those that can be queried using QuickBooks Change Data Capture
type CDCQueryable interface {
	Type
	ProcessCDC(cdc quickbooks.ChangeDataCapture) ([]map[string]any, error)
}

// DepCDCQueryable limits Types to those whos source can be queried using QuickBooks Change Data Capture
type DepCDCQueryable interface {
	DependentType
	MapType(sourceArray any) (map[string]map[string]bool, error)
	ProcessCDC(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache) ([]map[string]any, error)
}

// WHQueryable limits Types to those that can send a Webhook notification on update
type WHQueryable interface {
	Type
	ProcessWHBatch(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string) error
}

type DepWHReceivable interface {
	DependentType
	ProcessWHBatch(sourceArray any, cacheEntry *IdCache) ([]map[string]any, error)
}

// FiberyType establishes the base information required to create a datatype in Fibery
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

// QuickBooksType establishes the base functions required to retreive, process, and convert a QuickBooks entity to a Fibery entity
type QuickBooksType struct {
	FiberyType
	schemaGen      schemaGenFunc
	query          func(req Request) (Response, error)
	queryProcessor queryProcessorFunc
}

type schemaGenFunc func(entity any) (map[string]any, error)

type queryProcessorFunc func(entityArray any, schemaGen schemaGenFunc) ([]map[string]any, error)

// Query requests Type data from QuickBooks based on the Request parameters
func (t QuickBooksType) Query(req Request) (Response, error) {
	return t.query(req)
}

// ProcessQuery takes an array of QuickBooks entities and returns them as an array of corresponding Fibery entities
func (t QuickBooksType) ProcessQuery(array any) ([]map[string]any, error) {
	return t.queryProcessor(array, t.schemaGen)
}

type cdcProcessorFunc func(cdc quickbooks.ChangeDataCapture, schemaGen schemaGenFunc) ([]map[string]any, error)

// QuickBooksCDCType established the additional function(s) required to process Change Data Capture responses
type QuickBooksCDCType struct {
	QuickBooksType
	cdcProcessor cdcProcessorFunc
}

// ProcessCDC takes a non-specific Change Data Capture response and returns entities of the relevant type converted to Fibery schema
func (t QuickBooksCDCType) ProcessCDC(cdc quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	return t.cdcProcessor(cdc, t.schemaGen)
}

type WHBatchProcessorFunc func(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string, queryProcessor queryProcessorFunc, schemaGen schemaGenFunc, typeId string) error

// QuickBooksWHType established the additional function(s) required to process a webhook notifcation
type QuickBooksWHType struct {
	QuickBooksType
	whBatchProcessor WHBatchProcessorFunc
}

func (t QuickBooksWHType) ProcessWHBatch(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string) error {
	return t.whBatchProcessor(itemResponse, response, cache, realmId, t.queryProcessor, t.schemaGen, t.GetId())
}

// QuickBooksDualType requires the functions from both QuickBooksCDCType and QuickBooksWHType
type QuickBooksDualType struct {
	QuickBooksType
	cdcProcessor     cdcProcessorFunc
	whBatchProcessor WHBatchProcessorFunc
}

// ProcessCDC takes a non-specific Change Data Capture response and returns entities of the relevant type converted to Fibery schema
func (t QuickBooksDualType) ProcessCDC(cdc quickbooks.ChangeDataCapture) ([]map[string]any, error) {
	return t.cdcProcessor(cdc, t.schemaGen)
}

func (t QuickBooksDualType) ProcessWHBatch(itemResponse quickbooks.BatchItemResponse, response *map[string][]map[string]any, cache *cache.Cache, realmId string) error {
	return t.whBatchProcessor(itemResponse, response, cache, realmId, t.queryProcessor, t.schemaGen, t.GetId())
}

type depSchemaGenFunc func(entity any, source any) (map[string]any, error)

// DependentBaseType established the base functions required to process, extract, and convert dependent data from an array of source entities
type DependentBaseType struct {
	FiberyType
	schemaGen      depSchemaGenFunc
	queryProcessor func(sourceArray any, schemaGen depSchemaGenFunc) ([]map[string]any, error)
}

func (t DependentBaseType) ProcessQuery(array any) ([]map[string]any, error) {
	return t.queryProcessor(array, t.schemaGen)
}

// DependentDataType corresponds to a QuickBooksType which can only be requested through a query or read operation
type DependentDataType struct {
	DependentBaseType
	source QuickBooksType
}

func (t DependentDataType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentDataType) GetSourceId() string {
	return t.source.GetId()
}

// sourceMapperFunc maps the dependent Ids of a single corresponding source entity
type sourceMapperFunc func(source any) (map[string]bool, error)

// typeMapperFunc maps an array of source entities using the sourceMapperFunc for each source entity
type typeMapperFunc func(sourceArray any, sourceMapper sourceMapperFunc) (map[string]map[string]bool, error)

type depCDCProcessorFunc func(cdc quickbooks.ChangeDataCapture, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error)

type DependentCDCType struct {
	DependentBaseType
	source       QuickBooksCDCType
	sourceMapper sourceMapperFunc
	typeMapper   typeMapperFunc
	cdcProcessor depCDCProcessorFunc
}

func (t DependentCDCType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentCDCType) GetSourceId() string {
	return t.source.GetId()
}

// ProcessCDC takes a non-specific Change Data Capture response and returns dependent entities if the source type is included
func (t DependentCDCType) ProcessCDC(cdc quickbooks.ChangeDataCapture, idEntry *IdCache) ([]map[string]any, error) {
	return t.cdcProcessor(cdc, idEntry, t.sourceMapper, t.schemaGen)
}

// MapType creates a map of source & dependent entity ids to track changes from Change Data Capture and Webhook notifications
func (t DependentCDCType) MapType(sourceArray any) (map[string]map[string]bool, error) {
	return t.typeMapper(sourceArray, t.sourceMapper)
}

type depWHBatchProcessorFunc func(sourceArray any, cacheEntry *IdCache, sourceMapper sourceMapperFunc, schemaGen depSchemaGenFunc) ([]map[string]any, error)

type DependentWHType struct {
	DependentBaseType
	source           QuickBooksWHType
	sourceMapper     sourceMapperFunc
	whBatchProcessor depWHBatchProcessorFunc
}

func (t DependentWHType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentWHType) GetSourceId() string {
	return t.source.GetId()
}

func (t DependentWHType) ProcessWHBatch(sourceArray any, cacheEntry *IdCache) ([]map[string]any, error) {
	return t.whBatchProcessor(sourceArray, cacheEntry, t.sourceMapper, t.schemaGen)
}

type DependentDualType struct {
	DependentBaseType
	source           QuickBooksDualType
	sourceMapper     sourceMapperFunc
	typeMapper       typeMapperFunc
	whBatchProcessor depWHBatchProcessorFunc
	cdcProcessor     depCDCProcessorFunc
}

func (t DependentDualType) Query(req Request) (Response, error) {
	return t.source.Query(req)
}

func (t DependentDualType) GetSourceId() string {
	return t.source.GetId()
}

// ProcessCDC takes a non-specific Change Data Capture response and returns dependent entities if the source type is included
func (t DependentDualType) ProcessCDC(cdc quickbooks.ChangeDataCapture, idEntry *IdCache) ([]map[string]any, error) {
	return t.cdcProcessor(cdc, idEntry, t.sourceMapper, t.schemaGen)
}

// MapType creates a map of source & dependent entity ids to track changes from Change Data Capture and Webhook notifications
func (t DependentDualType) MapType(sourceArray any) (map[string]map[string]bool, error) {
	return t.typeMapper(sourceArray, t.sourceMapper)
}

func (t DependentDualType) ProcessWHBatch(sourceArray any, cacheEntry *IdCache) ([]map[string]any, error) {
	return t.whBatchProcessor(sourceArray, cacheEntry, t.sourceMapper, t.schemaGen)
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
