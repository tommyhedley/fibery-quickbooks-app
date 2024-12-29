// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package qbo

import (
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
	Cache         *cache.Cache
	Group         *singleflight.Group
	Token         *BearerToken
	StartPosition int
	OperationID   string
	RealmID       string
	LastSynced    string
	Types         []string
	Filter        map[string]any
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

type DataResponse[t any] struct {
	More bool
	Data t
}

type IDCacheEntry struct {
	OperationID string
	ItemIDs     map[string]map[string]bool
}

const IDCacheLifetime = 4 * time.Hour

// FiberyType represents any datatype that can be turned into a Fibery type/database.
type FiberyType interface {
	// TypeInfo returns a TypeArray containing the FiberyType ID and Name for implemented type as required by the Fibery /api/v1/synchronizer/config endpoint.
	// Returned ID value should match the datatype's name in the QuickBooks API.
	TypeInfo() TypeArray
	// Schema returns a map of FiberyType fields for implemented type as required by the Fibery /api/v1/synchronizer/config endpoint.
	Schema() map[string]Field
	// GetData handles data retreiaval for the given type and determines what type of sync is required.
	GetData(*DataRequest) (DataHandlerResponse, error)
}

// QBOPrimaryType represents Fibery types or databases that correspond to QuickBooks objects.
// Since they are queryable using the Quickbooks API, they have simpler data transformation requirements.
type QBOPrimaryType interface {
	FiberyType
	TransformItem() (map[string]any, error)
	TransformDataFS(data DataResponse[[]any]) ([]map[string]any, error)
	TransformDataDS(data ChangeDataCapture) ([]map[string]any, error)
	FullSync(*DataRequest) (DataResponse[[]any], error)
	DeltaSync(*DataRequest) (ChangeDataCapture, error)
}

// QBOSubtype represents Fibery types or databases that correspond to objects that are part of Quickbooks objects but not directly queryable.
// Invoice => InvoiceLine where InvoiceLine is a FiberySubtype of Invoice. This type may requuire caching to properly handle Delta and Webhook syncs.
type QBOSubtype interface {
	FiberyType
	TransformItem(parent QBOPrimaryType) (map[string]any, error)
	TransformDataFS(data DataResponse[[]any]) ([]map[string]any, error)
	TransformDataDS(data ChangeDataCapture, idCache IDCacheEntry) ([]map[string]any, error)
	FullSync(*DataRequest) (DataResponse[[]any], error)
	DeltaSync(*DataRequest) (ChangeDataCapture, error)
}

type TypeArray struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Types = map[string]FiberyType{}
var TypeInfo = []TypeArray{}
var Schema = make(map[string]map[string]Field)
var BaseTypes = map[string]bool{}

func RegisterType(t QBOPrimaryType) {
	Types[t.TypeInfo().ID] = t
	TypeInfo = append(TypeInfo, t.TypeInfo())
	Schema[t.TypeInfo().ID] = t.Schema()
	if _, ok := t.(QBOPrimaryType); ok {
		BaseTypes[t.TypeInfo().ID] = true
	}
}
