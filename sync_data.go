package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
	"golang.org/x/sync/singleflight"
)

func DataHandler(c *cache.Cache, group *singleflight.Group) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type syncType string

		const (
			DeltaSync syncType = "delta"
			FullSync  syncType = "full"
		)

		type nextPageConfig struct {
			StartPosition int `json:"startPosition"`
		}

		type pagination struct {
			HasNext        bool           `json:"hasNext"`
			NextPageConfig nextPageConfig `json:"nextPageConfig"`
		}

		type requestBody struct {
			RequestedType   string                               `json:"requestedType"`
			OperationID     string                               `json:"operationId"`
			Types           []string                             `json:"types"`
			Filter          map[string]any                       `json:"filter"`
			Account         qbo.FiberyAccountInfo                `json:"account"`
			LastSyncronized string                               `json:"lastSynchronizedAt"`
			Pagination      nextPageConfig                       `json:"pagination"`
			Schema          map[string]map[string]map[string]any `json:"schema"`
		}

		type responseBody struct {
			Items               []map[string]any `json:"items"`
			Pagination          pagination       `json:"pagination"`
			SynchronizationType syncType         `json:"synchronizationType"`
		}

		decoder := json.NewDecoder(r.Body)
		params := requestBody{}
		err := decoder.Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
			return
		}

		startPosition := params.Pagination.StartPosition
		reqType := params.RequestedType
		opID := params.OperationID

		if reqType == "" || opID == "" {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("request parameters missing type: %s and/or operation ID (%s", reqType, opID))
			return
		}

		rf, exists := qbo.GetRequestFunctions(reqType)
		if !exists {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type was not found: %s", reqType))
			return
		}

		if startPosition == 0 {
			startPosition = 1
		}

		req := qbo.RequestParameters{
			Cache:         c,
			Group:         group,
			Token:         &params.Account.BearerToken,
			StartPosition: startPosition,
			OperationID:   opID,
			RealmID:       params.Account.RealmID,
			LastSynced:    params.LastSyncronized,
			Filter:        params.Filter,
		}

		// Add error type checking for respond with trylater
		items, more, err := (*rf)(req)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to retreive full sync data"))
			return
		}

		sync := DeltaSync
		if params.LastSyncronized == "" {
			sync = FullSync
		}

		resp := responseBody{
			Items: items,
			Pagination: pagination{
				HasNext: more,
				NextPageConfig: nextPageConfig{
					StartPosition: startPosition + qbo.QueryPageSize,
				},
			},
			SynchronizationType: sync,
		}
		RespondWithJSON(w, http.StatusOK, resp)
	}
}
