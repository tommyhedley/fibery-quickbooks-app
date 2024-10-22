package synchronizer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

type FieldType string

const (
	ID        FieldType = "id"
	Text      FieldType = "text"
	Number    FieldType = "number"
	Date      FieldType = "date"
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

type Relation struct {
	Cardinality   string `json:"cardinality"`
	Name          string `json:"name"`
	TargetName    string `json:"targetName"`
	TargetType    string `json:"targetType"`
	TargetFieldID string `json:"targetFieldId"`
}

type Field struct {
	Ignore      bool           `json:"ignore,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	ReadOnly    bool           `json:"readonly,omitempty"`
	Type        FieldType      `json:"type,omitempty"`
	SubType     FieldSubtype   `json:"subType,omitempty"`
	Format      map[string]any `json:"format,omitempty"`
	Relation    *Relation      `json:"relation,omitempty"`
}

type CustomfieldSchema struct {
	ID         Field `json:"id"`
	Name       Field `json:"name"`
	Active     Field `json:"active"`
	SyncAction Field `json:"__syncAction"`
}

type CustomfieldData struct {
	ID         json.Number `json:"id" type:"string"`
	Name       string      `json:"name"`
	Active     bool        `json:"active"`
	SyncAction string      `json:"__syncAction,omitempty"`
}

type GroupSchema struct {
	ID         Field `json:"id"`
	Name       Field `json:"name"`
	Active     Field `json:"active"`
	SyncAction Field `json:"__syncAction"`
}

var group = GroupSchema{
	ID: Field{
		Name: "Id",
		Type: ID,
	},
	Name: Field{
		Name: "Name",
		Type: Text,
	},
	Active: Field{
		Name:     "Active",
		SubType:  Boolean,
		ReadOnly: false,
	},
	SyncAction: Field{
		Type: Text,
		Name: "Sync Action",
	},
}

type GroupData struct {
	ID         json.Number `json:"id" type:"string"`
	Name       string      `json:"name"`
	Active     bool        `json:"active"`
	SyncAction string      `json:"__syncAction,omitempty"`
}

// Timesheets, Users, & Jobcodes have dynamic data and schema based on selected custom fields
var user = map[string]Field{
	"id": {
		Name: "Id",
		Type: ID,
	},
	"name": {
		Name: "Name",
		Type: Text,
	},
	"display_name": {
		Name: "Display Name",
		Type: Text,
	},
	"first_name": {
		Name: "First Name",
		Type: Text,
	},
	"last_name": {
		Name: "Last Name",
		Type: Text,
	},
	"active": {
		Name:     "Active",
		SubType:  Boolean,
		ReadOnly: false,
	},
	"last_active": {
		Name: "Last Active",
		Type: Date,
	},
	"group_id": {
		Name: "Group ID",
		Type: Text,
		Relation: &Relation{
			Cardinality:   "many-to-one",
			Name:          "Group",
			TargetName:    "Users",
			TargetType:    "group",
			TargetFieldID: "id",
		},
	},
	"email": {
		Name:    "Email",
		SubType: Email,
	},
	"__syncAction": {
		Type: Text,
		Name: "Sync Action",
	},
}

type UserData struct {
	ID           json.Number `json:"id" type:"string"`
	Name         string      `json:"name"`
	DisplayName  string      `json:"display_name"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Active       bool        `json:"active"`
	LastActive   string      `json:"last_active"`
	GroupID      json.Number `json:"group_id" type:"string"`
	Email        string      `json:"email"`
	CustomFields any         `json:"customfields,omitempty"`
}

var timesheet = map[string]Field{
	"id": {
		Name: "Id",
		Type: ID,
	},
	"name": {
		Name: "Name",
		Type: Text,
	},
	"user_id": {
		Name: "User ID",
		Type: Text,
		Relation: &Relation{
			Cardinality:   "many-to-one",
			Name:          "User",
			TargetName:    "Timesheets",
			TargetType:    "user",
			TargetFieldID: "id",
		},
	},
	"created_by_user_id": {
		Name: "Creation User ID",
		Type: Text,
		Relation: &Relation{
			Cardinality:   "many-to-one",
			Name:          "Creation User",
			TargetName:    "Created Timesheets",
			TargetType:    "user",
			TargetFieldID: "id",
		},
	},
	"jobcode_id": {
		Name: "Jobcode ID",
		Type: Text,
		Relation: &Relation{
			Cardinality:   "many-to-one",
			Name:          "Jobcode",
			TargetName:    "Timesheets",
			TargetType:    "jobcode",
			TargetFieldID: "id",
		},
	},
	"locked": {
		Name:     "Locked",
		Type:     Text,
		SubType:  Boolean,
		ReadOnly: true,
	},
	"last_modified": {
		Name: "Last Modified",
		Type: Date,
	},
	"type": {
		Name:    "Timesheet Type",
		Type:    Text,
		SubType: SingleSelect,
	},
	"start": {
		Name: "Start",
		Type: Date,
	},
	"end": {
		Name: "End",
		Type: Date,
	},
	"date": {
		Name:    "Date",
		Type:    Date,
		SubType: Day,
	},
	"duration": {
		Name: "Duration (S)",
		Type: Number,
		Format: map[string]interface{}{
			"format":               "Number",
			"unit":                 "s",
			"hasThousandSeparator": true,
			"precision":            0,
		},
		SubType: Integer,
	},
	"duration_minutes": {
		Name: "Duration (M)",
		Type: Number,
		Format: map[string]interface{}{
			"format":               "Number",
			"unit":                 "m",
			"hasThousandSeparator": true,
			"precision":            2,
		},
	},
	"duration_hours": {
		Name: "Duration (H)",
		Type: Number,
		Format: map[string]interface{}{
			"format":               "Number",
			"unit":                 "h",
			"hasThousandSeparator": true,
			"precision":            2,
		},
	},
	"on_the_clock": {
		Name:     "On The Clock",
		Type:     Text,
		SubType:  Boolean,
		ReadOnly: true,
	},
	"__syncAction": {
		Type: Text,
		Name: "Sync Action",
	},
}

type TimesheetData struct {
	ID              json.Number `json:"id" type:"string"`
	UserID          json.Number `json:"user_id" type:"string"`
	CreatedByUserID json.Number `json:"created_by_user_id" type:"string"`
	JobcodeID       json.Number `json:"jobcode_id" type:"string"`
	Locked          json.Number `json:"locked" type:"string"`
	Notes           string      `json:"notes"`
	LastModified    string      `json:"last_modified"`
	Type            string      `json:"type"`
	Start           string      `json:"start"`
	End             string      `json:"end"`
	Date            string      `json:"date"`
	DurationSeconds int         `json:"duration"`
	OnTheClock      bool        `json:"on_the_clock"`
	CustomFields    any         `json:"customfields,omitempty"`
}

var jobcode = map[string]Field{
	"id": {
		Name: "Id",
		Type: ID,
	},
	"name": {
		Name: "Name",
		Type: Text,
	},
	"parent_id": {
		Name: "Parent ID",
		Type: Text,
		Relation: &Relation{
			Cardinality:   "many-to-one",
			Name:          "Parent",
			TargetName:    "Jobs",
			TargetType:    "jobcode",
			TargetFieldID: "id",
		},
	},
	"type": {
		Name:    "Type",
		Type:    Text,
		SubType: SingleSelect,
	},
	"billable": {
		Name:    "Billable",
		Type:    Text,
		SubType: Boolean,
	},
	"active": {
		Name:    "Active",
		Type:    Text,
		SubType: Boolean,
	},
	"connected_with_quickbooks": {
		Name:    "Connected with QuickBooks",
		Type:    Text,
		SubType: Boolean,
	},
	"__syncAction": {
		Type: Text,
		Name: "Sync Action",
	},
}

type JobcodeData struct {
	ID                    json.Number `json:"id" type:"string"`
	Name                  string      `json:"name"`
	ParentID              json.Number `json:"parent_id" type:"string"`
	Type                  string      `json:"type"`
	Billable              bool        `json:"billable"`
	Active                bool        `json:"active"`
	ConnectWithQuickbooks bool        `json:"connect_with_quickbooks"`
	CustomFields          any         `json:"customfields,omitempty"`
}

func SchemaHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Types   []string `json:"types"`
		Filter  Filter   `json:"filter"`
		Account struct {
			AccessToken string `json:"access_token"`
		} `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	type customFieldRequest struct {
		Active           string `url:"active"`
		SupplementalData string `url:"supplemental_data"`
	}

	type customFieldResponse struct {
		ID        json.Number `json:"id" type:"string"`
		Name      string      `json:"name"`
		AppliesTo string      `json:"applies_to"`
	}

	allType := map[string]interface{}{
		"user":      user,
		"group":     group,
		"timesheet": timesheet,
		"jobcode":   jobcode,
	}

	customFieldReq := customFieldRequest{
		Active:           "yes",
		SupplementalData: "no",
	}

	customfields, _, requestError := utils.APIRequest[customFieldRequest, any, customFieldResponse](&customFieldReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/customfields", params.Account.AccessToken, "customfields")
	if requestError.Err != nil {
		if requestError.TryLater {
			utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with customfield name request: %w", requestError.Err))
		return
	}

	for _, customfield := range customfields {
		allType[customfield.ID.String()] = CustomfieldSchema{
			ID: Field{
				Name: "Id",
				Type: ID,
			},
			Name: Field{
				Name: "Name",
				Type: Text,
			},
			Active: Field{
				Name:    "Active",
				Type:    Text,
				SubType: Boolean,
			},
			SyncAction: Field{
				Type: Text,
				Name: "Sync Action",
			},
		}
		fieldID := customfield.ID.String()
		fieldName := customfield.Name + " Name"
		switch customfield.AppliesTo {
		case "timesheet":
			timesheet[fieldID] = Field{
				Name: fieldName,
				Type: Text,
				Relation: &Relation{
					Cardinality:   "many-to-one",
					Name:          customfield.Name,
					TargetName:    "Timesheets",
					TargetType:    customfield.ID.String(),
					TargetFieldID: "name",
				},
			}
		case "user":
			user[fieldID] = Field{
				Name: fieldName,
				Type: Text,
				Relation: &Relation{
					Cardinality:   "many-to-one",
					Name:          customfield.Name,
					TargetName:    "Users",
					TargetType:    customfield.ID.String(),
					TargetFieldID: "name",
				},
			}
		case "jobcode":
			jobcode[fieldID] = Field{
				Name: fieldName,
				Type: Text,
				Relation: &Relation{
					Cardinality:   "many-to-one",
					Name:          customfield.Name,
					TargetName:    "Jobcodes",
					TargetType:    customfield.ID.String(),
					TargetFieldID: "name",
				},
			}
		default:

		}
	}

	returnType := make(map[string]interface{})

	for name, fields := range allType {
		for _, t := range params.Types {
			if name == t {
				returnType[name] = fields
			}
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, returnType)
}
