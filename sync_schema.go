package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
)

func SchemaHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Types   []string              `json:"types"`
		Filter  Filter                `json:"filter"`
		Account qbo.FiberyAccountInfo `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	selectedSchemas := make(map[string]map[string]qbo.Field)

	for _, t := range req.Types {
		selectedSchemas[t] = qbo.Types[t].Schema()
	}

	RespondWithJSON(w, http.StatusOK, selectedSchemas)
}
