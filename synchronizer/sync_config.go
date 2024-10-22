package synchronizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

type DBTypes string

type SyncConfig struct {
	Types    []DB        `json:"types"`
	Filters  []Filter    `json:"filters"`
	Webhooks SyncWebhook `json:"webhooks,omitempty"`
}
type DB struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Filter struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Datalist bool   `json:"datalist,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Secured  bool   `json:"secured,omitempty"`
}
type SyncWebhook struct {
	Enabled bool `json:"enabled,omitempty"`
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := SyncConfig{
		Types: []DB{
			{
				ID:   "user",
				Name: "User",
			},
			{
				ID:   "group",
				Name: "Group",
			},
			{
				ID:   "timesheet",
				Name: "Timesheet",
			},
			{
				ID:   "jobcode",
				Name: "Jobcode",
			},
		},
		Filters: []Filter{
			{
				ID:       "timesheetStart",
				Title:    "Sync timesheets on/after this date. Jan 1, 2020 will be used if selection is earlier or empty.",
				Type:     "datebox",
				Optional: true,
				Secured:  true,
			},
			{
				ID:      "includeOTC",
				Title:   "Include on the clock/active timesheets with sync?",
				Type:    "bool",
				Secured: true,
			},
		},
		Webhooks: SyncWebhook{
			Enabled: true,
		},
	}

	type parameters struct {
		Types   []string `json:"types"`
		Account struct {
			AccessToken string `json:"access_token"`
		} `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	type customFieldRequest struct {
		Active           string `url:"active"`
		SupplementalData string `url:"supplemental_data"`
	}

	type customFieldResponse struct {
		ID   json.Number `json:"id" type:"string"`
		Name string      `json:"name"`
	}

	customFieldReq := customFieldRequest{
		Active:           "yes",
		SupplementalData: "no",
	}

	customFields, _, requestError := utils.GetData[customFieldRequest, customFieldResponse](&customFieldReq, "https://rest.tsheets.com/api/v1/customfields", params.Account.AccessToken, "customfields")
	if requestError.Err != nil {
		if requestError.TryLater {
			utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with customfield name request: %w", requestError.Err))
		return
	}

	for _, customField := range customFields {
		config.Types = append(config.Types, DB{
			ID:   customField.ID.String(),
			Name: customField.Name,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, config)
}

func ValidateFiltersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Types   []string       `json:"types"`
		Filter  map[string]any `json:"filter"`
		Account struct {
			Token string `json:"token"`
		} `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	if params.Filter["timesheetStart"] != nil {
		if datesString, ok := params.Filter["timesheetStart"].(string); ok {
			_, err := time.Parse(time.RFC3339, datesString)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("couldn't parse date string"))
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, nil)
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Filter input value is an invalid type"))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, nil)
}
