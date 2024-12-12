// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package qbo

import (
	"sync"
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

type CacheEntry[t any] struct {
	Data           []t
	ProcessedTypes map[string]bool
	More           bool
	mu             sync.Mutex
}

func allSubtypesProcessed(processedTypes map[string]bool) bool {
	for _, processed := range processedTypes {
		if !processed {
			return false
		}
	}
	return true
}

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

// Data Request Definitions & Parameters

type FullSyncRequest struct {
	Cache         *cache.Cache
	Group         *singleflight.Group
	Token         *BearerToken
	StartPosition int
	OperationID   string
	RealmID       string
	Filter        map[string]any
}

type DeltaSyncRequest struct {
	Cache       *cache.Cache
	Group       *singleflight.Group
	Token       *BearerToken
	OperationID string
	RealmID     string
	Types       []string
	LastSynced  time.Time
	Filter      map[string]any
}

type WebhookRequest struct {
}

type FiberyType interface {
	// TypeInfo returns a TypeArray containing the FiberyType ID and Name for implemented type as required by the Fibery /api/v1/synchronizer/config endpoint.
	TypeInfo() TypeArray
	// Schema returns a map of FiberyType fields for implemented type as required by the Fibery /api/v1/synchronizer/config endpoint.
	Schema() map[string]Field
	// TransformData transforms the data from the QuickBooks API into the format required by the Fibery /api/v1/synchronizer/data endpoint.
	TransformData(params ...any) (any, error)
	// FullSync performs a full sync of the given type data from the QuickBooks API into the Fibery /api/v1/synchronizer/data endpoint.
	// This function supports pagination and will return a boolean indicating if there is more data to be fetched.
	FullSync(*FullSyncRequest) ([]map[string]any, bool, error)
	// DeltaSync performs a delta sync using the ChangeDataCapture of the given type data from the QuickBooks API into the Fibery /api/v1/synchronizer/data endpoint.
	// Since CDC does not support pagination, we return an error and suggest a full sync if CDC amount is greater than allowed (1000).
	DeltaSync(*DeltaSyncRequest) ([]map[string]any, error)
	// Webhook performs a request for the notifcation items of the given data type from the QuickBooks API into the Fibery /api/v1/synchronizer/webhooks/transformer endpoint.
	Webhook(*WebhookRequest) ([]map[string]any, error)
}

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

type TypeArray struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BaseType bool   `json:"-"`
}

var Types = map[string]FiberyType{}
var TypeInfo = []TypeArray{}
var BaseTypes = map[string]bool{}
var Schema = make(map[string]map[string]Field)

func RegisterType(t FiberyType) {
	Types[t.TypeInfo().ID] = t
	TypeInfo = append(TypeInfo, t.TypeInfo())
	BaseTypes[t.TypeInfo().ID] = t.TypeInfo().BaseType
	Schema[t.TypeInfo().ID] = t.Schema()
}
