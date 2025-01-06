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
	LastSynced     string
	RequestedType  string
	RequestedTypes []string
	CDCTypes       []string
	Filter         map[string]any
	Cache          *cache.Cache
	Group          *singleflight.Group
	Token          *qbo.BearerToken
}

type Response struct {
	Data     []any
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
}

type FiberyType struct {
	Id     string
	Name   string
	Schema map[string]fibery.Field
}

func (f FiberyType) GetId() string {
	return f.Id
}

func (f FiberyType) GetName() string {
	return f.Name
}

func (f FiberyType) GetSchema() map[string]fibery.Field {
	return f.Schema
}

type QuickbooksType struct {
	FiberyType
	SchemaTransformer func(entity any) (map[string]any, error)
	DataQuery         func(req Request) (Response, error)
}

func (t QuickbooksType) Query(req Request) (Response, error) {
	return t.DataQuery(req)
}

type QuickbooksCDCType struct {
	QuickbooksType
	ChangeDataCaptureProcessor func(cdc qbo.ChangeDataCapture) ([]map[string]any, error)
}

type QuickbooksWHType struct {
	QuickbooksType
}

type QuickbooksDualType struct {
	QuickbooksType
	ChangeDataCaptureProcessor func(cdc qbo.ChangeDataCapture) ([]map[string]any, error)
}

type DepSchemaTransformerFunc func(entity any, source any) (map[string]any, error)
type SourceMapperFunc func(source any) (map[string]bool, error)
type TypeMapperFunc func(sourceArray any, sourceMapper SourceMapperFunc) (map[string]map[string]bool, error)

type DependentDataType struct {
	FiberyType
	Source                     Type
	SchemaTransformer          DepSchemaTransformerFunc
	SourceMapper               SourceMapperFunc
	TypeMapper                 TypeMapperFunc
	QueryProcessor             func(sourceArray any, schemaTransformer DepSchemaTransformerFunc) ([]map[string]any, error)
	ChangeDataCaptureProcessor func(cdc qbo.ChangeDataCapture, cacheEntry *IdCache, sourceMapper SourceMapperFunc, schemaTransformer DepSchemaTransformerFunc) ([]map[string]any, error)
}

func (t DependentDataType) Query(req Request) (Response, error) {
	return t.Source.Query(req)
}

var Types = map[string]*Type{}

func RegisterType(t Type) {
	Types[t.GetId()] = &t
}

func TestFiberyType(t Type) {
	t.GetId()
	t.GetName()
	t.GetSchema()
	t.Query(Request{})
}

func FormatJSON(data interface{}) string {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Failed to generate json", err)
		return ""
	}
	return string(prettyJSON)
}
