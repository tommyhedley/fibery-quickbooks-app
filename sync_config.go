package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
)

type Filter struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Datalist bool   `json:"datalist,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Secured  bool   `json:"secured,omitempty"`
}

func SyncConfigHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Account qbo.FiberyAccountInfo `json:"account"`
	}

	type configTypes struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type syncWebhook struct {
		Enabled bool   `json:"enabled,omitempty"`
		Type    string `json:"type,omitempty"`
	}

	type syncConfig struct {
		Types    []configTypes `json:"types"`
		Filters  []Filter      `json:"filters"`
		Webhooks syncWebhook   `json:"webhooks,omitempty"`
	}

	availableTypes := []configTypes{}

	for _, t := range qbo.FiberyTypes {
		typ := t()
		availableTypes = append(availableTypes, configTypes{
			ID:   typ.ID(),
			Name: typ.Name(),
		})
	}

	config := syncConfig{
		Types:   availableTypes,
		Filters: []Filter{},
		Webhooks: syncWebhook{
			Enabled: true,
			Type:    "ui",
		},
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, config)
}

func ValidateFiltersHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Types   []string              `json:"types"`
		Filter  map[string]any        `json:"filter"`
		Account qbo.FiberyAccountInfo `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
