package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
	"golang.org/x/sync/singleflight"
)

func DataHandler(c *cache.Cache, group *singleflight.Group) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			RequestedType     string                               `json:"requestedType"`
			OperationID       string                               `json:"operationId"`
			Types             []string                             `json:"types"`
			Filter            map[string]any                       `json:"filter"`
			Account           qbo.FiberyAccountInfo                `json:"account"`
			LastSyncronizedAt string                               `json:"lastSynchronizedAt"`
			Pagination        qbo.NextPageConfig                   `json:"pagination"`
			Schema            map[string]map[string]map[string]any `json:"schema"`
		}

		decoder := json.NewDecoder(r.Body)
		params := requestBody{}
		err := decoder.Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
			return
		}

		startPosition := params.Pagination.StartPosition

		if startPosition == 0 {
			startPosition = 1
		}

		if params.RequestedType == "" {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requestedType is required"))
			return
		}

		if params.OperationID == "" {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("operationId is required"))
			return
		}

		req := qbo.DataRequest{
			Cache:         c,
			Group:         group,
			Token:         &params.Account.BearerToken,
			StartPosition: startPosition,
			OperationID:   params.OperationID,
			RealmID:       params.Account.RealmID,
			LastSynced:    params.LastSyncronizedAt,
			Types:         params.Types,
			Filter:        params.Filter,
		}

		datatype := qbo.Types[params.RequestedType]
		if datatype == nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type was not found: %s", params.RequestedType))
			return
		}

		slog.Info(fmt.Sprintf("Getting data for %s", params.RequestedType))

		res, err := datatype.GetData(&req)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to retrieve data: %w", err))
			return
		}

		RespondWithJSON(w, http.StatusOK, res)
	}
}
