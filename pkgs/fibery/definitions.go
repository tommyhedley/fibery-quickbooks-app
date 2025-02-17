package fibery

import "time"

// AppConfig defines the integration app
type AuthField struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional,omitempty"`
	Value       string `json:"value,omitempty"`
}

type Authentication struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Fields      []AuthField `json:"fields,omitempty"`
}

type ResponsibleFor struct {
	DataSynchronization bool `json:"dataSynchronization"`
	Automations         bool `json:"automations,omitempty"`
}

type ActionArg struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Type         string `json:"type"`
	TextTemplate bool   `json:"textTemplateSupported,omitempty"`
}

type Action struct {
	ActionID    string      `json:"action"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Args        []ActionArg `json:"args"`
}

type AppConfig struct {
	ID             string           `json:"id"`
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
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SyncFilter struct {
	ID       string `json:"id"`
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

const DateFormat = time.RFC3339

type Webhook struct {
	WebhookID   string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
}
