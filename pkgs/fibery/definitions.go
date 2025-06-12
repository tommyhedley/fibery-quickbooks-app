package fibery

import "time"

// AppConfig defines the integration app
type AuthField struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional,omitempty"`
	Value       string `json:"value,omitempty"`
}

type Authentication struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Fields      []AuthField `json:"fields,omitempty"`
}

type ResponsibleFor struct {
	DataSynchronization bool `json:"dataSynchronization"`
	Automations         bool `json:"automations,omitempty"`
}

type ActionArgumentType string

const (
	TextArg     ActionArgumentType = "text"
	TextAreaArg ActionArgumentType = "textarea"
)

type ActionArg struct {
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description,omitempty"`
	ArgType      ActionArgumentType `json:"type"`
	TextTemplate bool               `json:"textTemplateSupported,omitempty"`
}

type Action struct {
	ActionId    string      `json:"action"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Args        []ActionArg `json:"args"`
}

type AppConfig struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	Website        string           `json:"website,omitempty"`
	Version        string           `json:"version"`
	Description    string           `json:"description"`
	Authentication []Authentication `json:"authentication"`
	Sources        []string         `json:"sources"`
	ResponsibleFor ResponsibleFor   `json:"responsibleFor"`
	Actions        []Action         `json:"actions"`
}

type SyncConfigTypes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SyncFilter struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Datalist bool   `json:"datalist,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Secured  bool   `json:"secured,omitempty"`
}

type SyncConfigWebhook struct {
	Enabled bool   `json:"enabled,omitempty"`
	Type    string `json:"type,omitempty"`
}

type SyncConfig struct {
	Types    []SyncConfigTypes `json:"types"`
	Filters  []SyncFilter      `json:"filters"`
	Webhooks SyncConfigWebhook `json:"webhooks,omitempty"`
}

// Fibery Schema Definitions
type FieldType string

const (
	Id        FieldType = "id"
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
	File         FieldSubtype = "file"
	Daterange    FieldSubtype = "date-range"
	Title        FieldSubtype = "title"
	SingleSelect FieldSubtype = "single-select"
	MultiSelect  FieldSubtype = "multi-select"
	Day          FieldSubtype = "day"
)

type CardinalityType string

const (
	OTO CardinalityType = "one-to-one"
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
	Name        string           `json:"name"`
	Ignore      bool             `json:"ignore,omitempty"`
	Description string           `json:"description,omitempty"`
	ReadOnly    bool             `json:"readonly,omitempty"`
	Type        FieldType        `json:"type,omitempty"`
	SubType     FieldSubtype     `json:"subType,omitempty"`
	Format      map[string]any   `json:"format,omitempty"`
	Options     []map[string]any `json:"options,omitempty"`
	Relation    *Relation        `json:"relation,omitempty"`
}

type SyncType string

const (
	Delta SyncType = "delta"
	Full  SyncType = "full"
)

type SyncAction string

const (
	SET    SyncAction = "SET"
	REMOVE SyncAction = "REMOVE"
)

type NextPageConfig struct {
	Page int `json:"page"`
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

const DateFormat = time.RFC3339

type Webhook struct {
	WebhookId   string `json:"id"`
	WorkspaceId string `json:"workspaceId"`
}

type WebhookData map[string][]map[string]any

type WebhookTransformResponse struct {
	Data WebhookData `json:"data"`
}

type Type interface {
	Id() string
	Name() string
	Schema() map[string]Field
}
