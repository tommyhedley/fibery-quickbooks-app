package sync

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/internal/utils"
)

type TypeArray struct {
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

type SyncConfig struct {
	Types    []TypeArray `json:"types"`
	Filters  []Filter    `json:"filters"`
	Webhooks SyncWebhook `json:"webhooks,omitempty"`
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Account struct {
			AccessToken string `json:"access_token"`
		} `json:"account"`
	}

	config := SyncConfig{
		Types:   Types,
		Filters: []Filter{},
		Webhooks: SyncWebhook{
			Enabled: true,
		},
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, config)
}

func ValidateFiltersHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Types   []string       `json:"types"`
		Filter  map[string]any `json:"filter"`
		Account struct {
			Token string `json:"token"`
		} `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil)
}
