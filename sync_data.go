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
			StartPosition:  startPosition,
			OperationID:    params.OperationID,
			RealmID:        params.Account.RealmID,
			LastSynced:     params.LastSyncronizedAt,
			RequestedType:  params.RequestedType,
			RequestedTypes: params.Types,
			CDCTypes:       reqTypesToCDCTypes(params.Types),
			Filter:         params.Filter,
			Cache:          c,
			Group:          group,
			Token:          &params.Account.BearerToken,
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

func reqTypesToCDCTypes(requestedTypes []string) []string {
	typeSet := make(map[string]struct{})
	var CDCTypes []string

	for _, dataType := range requestedTypes {
		if _, ok := qbo.Types[dataType]; ok {
			if cdcType, ok := qbo.Types[dataType].(qbo.CDCDataType); ok {
				if _, exists := typeSet[cdcType.ID()]; !exists {
					typeSet[cdcType.ID()] = struct{}{}
					CDCTypes = append(CDCTypes, cdcType.ID())
				}
			}
			if depType, ok := qbo.Types[dataType].(qbo.DependentDataType); ok {
				if _, exists := typeSet[depType.ParentID()]; !exists {
					typeSet[depType.ParentID()] = struct{}{}
					CDCTypes = append(CDCTypes, depType.ParentID())
				}
			}
		}
	}

	return CDCTypes
}
