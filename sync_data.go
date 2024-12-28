package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

		var syncType qbo.SyncType
		var lastSyncTime time.Time

		if params.LastSyncronizedAt == "" {
			syncType = qbo.FullSync
		} else {
			syncType = qbo.DeltaSync
			lastSyncTime, err = time.Parse(time.RFC3339, params.LastSyncronizedAt)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to parse lastSyncronizedAt: %w", err))
				return
			}
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

		var items []map[string]any
		var more bool

		if syncType == qbo.DeltaSync {
			// CDC Request
			CDCRequestTypes := []string{}
			for _, t := range params.Types {
				if qbo.BaseTypes[t] {
					CDCRequestTypes = append(CDCRequestTypes, t)
				}
			}

			req := qbo.DeltaSyncRequest{
				Cache:       c,
				Group:       group,
				Token:       &params.Account.BearerToken,
				OperationID: params.OperationID,
				RealmID:     params.Account.RealmID,
				Types:       CDCRequestTypes,
				LastSynced:  lastSyncTime,
				Filter:      params.Filter,
			}

			datatype := qbo.Types[params.RequestedType]
			if datatype == nil {
				RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type was not found: %s", params.RequestedType))
				return
			}

			items, err = datatype.DeltaSync(&req)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to retrieve delta sync data: %w", err))
				return
			}
		} else {
			// Full Sync Request
			req := qbo.FullSyncRequest{
				Cache:         c,
				Group:         group,
				Token:         &params.Account.BearerToken,
				StartPosition: startPosition,
				OperationID:   params.OperationID,
				RealmID:       params.Account.RealmID,
				Filter:        params.Filter,
			}

			datatype := qbo.Types[params.RequestedType]
			if datatype == nil {
				RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type was not found: %s", params.RequestedType))
				return
			}

			items, more, err = datatype.FullSync(&req)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to retrieve full sync data: %w", err))
				return
			}
		}

		resp := responseBody{
			Items: items,
			Pagination: pagination{
				HasNext: more,
				NextPageConfig: nextPageConfig{
					StartPosition: startPosition + qbo.QueryPageSize,
				},
			},
			SynchronizationType: syncType,
		}
		RespondWithJSON(w, http.StatusOK, resp)
	}
}
