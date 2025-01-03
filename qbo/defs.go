// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package qbo

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

type CustomField struct {
	DefinitionId string `json:"DefinitionId,omitempty"`
	StringValue  string `json:"StringValue,omitempty"`
	Type         string `json:"Type,omitempty"`
	Name         string `json:"Name,omitempty"`
}

// Date represents a Quickbooks date
type Date struct {
	time.Time `json:",omitempty"`
}

// UnmarshalJSON removes time from parsed date
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	d.Time, err = time.Parse(qboDateFormat, string(b))
	if err != nil {
		d.Time, err = time.Parse(qboDayFormat, string(b))
	}

	return err
}

func (d Date) String() string {
	return d.Format(qboDateFormat)
}

// QueryResponse defines the response format for Quickbooks queries
// This format is used for single type query responses, Change Data Capture responses, and batch query responses
type QueryResponse struct {
	Items         []QuickbooksDataType
	StartPosition int `json:"startPosition"`
	MaxResults    int `json:"maxResults"`
	TotalCount    int `json:"totalCount,omitempty"`
}

func (qr *QueryResponse) UnmarshalJSON(data []byte) error {
	type alias QueryResponse
	var temp struct {
		alias
	}
	if err := json.Unmarshal(data, &temp.alias); err != nil {
		return fmt.Errorf("failed to unmarshal known fields in QueryResponse: %w", err)
	}

	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		return fmt.Errorf("failed to unmarshal entire QueryResponse: %w", err)
	}

	delete(generic, "startPosition")
	delete(generic, "maxResults")
	delete(generic, "totalCount")

	var allItems []QuickbooksDataType
	for key, raw := range generic {
		constructor, ok := quickbooksTypes[key]
		if !ok {
			continue
		}

		var rawObjects []json.RawMessage
		if err := json.Unmarshal(raw, &rawObjects); err != nil {
			return fmt.Errorf("failed to unmarshal QueryResponse Items, array was expected but not found: %w", err)
		}

		for _, rawObj := range rawObjects {
			obj := constructor()
			if err := json.Unmarshal(rawObj, obj); err != nil {
				return fmt.Errorf("failed to unmarshal QueryResponse Item: %w", err)
			}
			allItems = append(allItems, obj)
		}
	}

	// 4) Store known fields and parsed items
	qr.StartPosition = temp.alias.StartPosition
	qr.MaxResults = temp.alias.MaxResults
	qr.TotalCount = temp.alias.TotalCount
	qr.Items = allItems
	return nil
}

// EmailAddress represents a QuickBooks email address.
type EmailAddress struct {
	Address string `json:",omitempty"`
}

// EndpointUrl specifies the endpoint to connect to
type EndpointUrl string

const (
	QueryPageSize    = 1000
	qboDateFormat    = "2006-01-02T15:04:05-07:00"
	qboDayFormat     = "2006-01-02"
	fiberyDateFormat = "2020-01-22T01:02:23.977Z"
)

func (u EndpointUrl) String() string {
	return string(u)
}

// MemoRef represents a QuickBooks MemoRef object.
type MemoRef struct {
	Value string `json:"value,omitempty"`
}

// MetaData is a timestamp of genesis and last change of a Quickbooks object.
type MetaData struct {
	CreateTime      Date `json:",omitempty"`
	LastUpdatedTime Date `json:",omitempty"`
}

// PhysicalAddress represents a QuickBooks address.
type PhysicalAddress struct {
	Id string `json:"Id,omitempty"`
	// These lines are context-dependent! Read the QuickBooks API carefully.
	Line1   string `json:",omitempty"`
	Line2   string `json:",omitempty"`
	Line3   string `json:",omitempty"`
	Line4   string `json:",omitempty"`
	Line5   string `json:",omitempty"`
	City    string `json:",omitempty"`
	Country string `json:",omitempty"`
	// A.K.A. State.
	CountrySubDivisionCode string `json:",omitempty"`
	PostalCode             string `json:",omitempty"`
	Lat                    string `json:",omitempty"`
	Long                   string `json:",omitempty"`
}

// ReferenceType represents a QuickBooks reference to another object.
type ReferenceType struct {
	Value string `json:"value,omitempty"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}

// TelephoneNumber represents a QuickBooks phone number.
type TelephoneNumber struct {
	FreeFormNumber string `json:",omitempty"`
}

// WebSiteAddress represents a Quickbooks Website
type WebSiteAddress struct {
	URI string `json:",omitempty"`
}

// Response & Request Format For Fibery Account Info
type FiberyAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmID string `json:"realmId,omitempty"`
	BearerToken
}

// Fibery Schema Definitions
type FieldType string

const (
	ID        FieldType = "id"
	Text      FieldType = "text"
	Number    FieldType = "number"
	DateType  FieldType = "date"
	TextArray FieldType = "array[text]"
)

type FieldSubtype string

const (
	URL          FieldSubtype = "url"
	Integer      FieldSubtype = "integer"
	Email        FieldSubtype = "email"
	Boolean      FieldSubtype = "boolean"
	HTML         FieldSubtype = "html"
	MD           FieldSubtype = "md"
	Files        FieldSubtype = "files"
	Daterange    FieldSubtype = "date-range"
	Title        FieldSubtype = "title"
	SingleSelect FieldSubtype = "single-select"
	MultiSelect  FieldSubtype = "multi-select"
	Day          FieldSubtype = "day"
)

type CardinalityType string

const (
	OTO CardinalityType = "one-to-one"
	OTM CardinalityType = "one-to-many"
	MTO CardinalityType = "many-to-one"
	MTM CardinalityType = "many-to-many"
)

type Relation struct {
	Cardinality   CardinalityType `json:"cardinality"`
	Name          string          `json:"name"`
	TargetName    string          `json:"targetName"`
	TargetType    string          `json:"targetType"`
	TargetFieldID string          `json:"targetFieldId"`
}

type Field struct {
	Ignore      bool             `json:"ignore,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	ReadOnly    bool             `json:"readonly,omitempty"`
	Type        FieldType        `json:"type,omitempty"`
	SubType     FieldSubtype     `json:"subType,omitempty"`
	Format      map[string]any   `json:"format,omitempty"`
	Options     []map[string]any `json:"options,omitempty"`
	Relation    *Relation        `json:"relation,omitempty"`
}

// Data Type & Handler Definitions
type SyncType string

const (
	DeltaSync SyncType = "delta"
	FullSync  SyncType = "full"
)

type SyncAction string

const (
	SET    SyncAction = "SET"
	REMOVE SyncAction = "REMOVE"
)

type DataRequest struct {
	StartPosition  int
	OperationID    string
	RealmID        string
	LastSynced     string
	RequestedType  string
	RequestedTypes []string
	CDCTypes       []string
	Filter         map[string]any
	Cache          *cache.Cache
	Group          *singleflight.Group
	Token          *BearerToken
}

type NextPageConfig struct {
	StartPosition int `json:"startPosition"`
}

type Pagination struct {
	HasNext        bool           `json:"hasNext"`
	NextPageConfig NextPageConfig `json:"nextPageConfig"`
}

type DataHandlerResponse struct {
	Items               []map[string]any `json:"items"`
	Pagination          Pagination       `json:"pagination"`
	SynchronizationType SyncType         `json:"synchronizationType"`
}

type DataResponse struct {
	Data any
	More bool
}

type DependentDataIDCache struct {
	OperationID string
	IDs         map[string]map[string]bool
}

const IDCacheLifetime = 4 * time.Hour

// FiberyDataType defines all of the functions needed to be a type in Fibery
type FiberyDataType interface {
	ID() string
	Name() string
	Schema() map[string]Field
	GetData(req *DataRequest) (DataHandlerResponse, error)
}

// QuickbooksDataType defines the available base functions of all the Quickbooks data types for the integration
type QuickbooksDataType interface {
	FiberyDataType
	TransformItem() (map[string]any, error)
	getFullData(req *DataRequest) (DataResponse, error)
	transformFullData(data DataResponse) ([]map[string]any, error)
}

type CDCDataType interface {
	QuickbooksDataType
	transformChangeDataCapture(cdc ChangeDataCapture) ([]map[string]any, error)
}

type ParentDataType interface {
	QuickbooksDataType
	Dependents() map[string][]DependentDataType
	TransformItemAndDependents() (map[string][]map[string]any, error)
}

type DependentDataType interface {
	FiberyDataType
	Parent() string
	transformItem(parent any) (map[string]any, error)
	transformChangeDataCapture(cdc ChangeDataCapture, idCache *DependentDataIDCache) ([]map[string]any, error)
}

var FiberyTypes = map[string]func() FiberyDataType{}
var quickbooksTypes = map[string]func() QuickbooksDataType{}

func RegisterType(t FiberyDataType) {
	FiberyTypes[t.ID()] = func() FiberyDataType {
		typ := reflect.TypeOf(t)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		return reflect.New(typ).Interface().(FiberyDataType)
	}

	if qb, ok := t.(QuickbooksDataType); ok {
		quickbooksTypes[qb.ID()] = func() QuickbooksDataType {
			typ := reflect.TypeOf(qb)
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			return reflect.New(typ).Interface().(QuickbooksDataType)
		}
	}
}

// temp

func FormatJSON(data interface{}) string {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Failed to generate json", err)
		return ""
	}
	return string(prettyJSON)
}

func TestQuickbooksDataType(t QuickbooksDataType) {
	return
}

func TestCDCDataType(t CDCDataType) {
	return
}

func TestParentDataType(t ParentDataType) {
	return
}

func TestDependentDataType(t DependentDataType) {
	return
}
