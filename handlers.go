package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/data"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/fibery"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/qbo"
	"golang.org/x/sync/singleflight"
)

type Integration struct {
	cache      *cache.Cache
	group      *singleflight.Group
	appConfig  fibery.AppConfig
	syncConfig fibery.SyncConfig
}

func (i *Integration) AppConfigHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, i.appConfig)
}

func (Integration) AccountValidateHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Id     string `json:"id"`
		Fields struct {
			integrationAccountInfo
		} `json:"fields"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	client, err := qbo.NewClient(reqBody.Fields.RealmID, &reqBody.Fields.BearerToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create new client: %w", err))
		return
	}

	token := reqBody.Fields.BearerToken

	refreshNeeded, err := token.RefreshNeeded()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to determine if token refresh is needed: %w", err))
		return
	}

	if refreshNeeded {
		newToken, err := client.RefreshToken(token.RefreshToken)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to refresh token: %w", err))
			return
		}
		token = *newToken
	}

	info, err := client.FindCompanyInfo()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to find company info: %w", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, integrationAccountInfo{
		Name:        info.CompanyName,
		RealmID:     reqBody.Fields.RealmID,
		BearerToken: token,
	},
	)
}

func (i *Integration) SyncConfigHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, i.syncConfig)
}

func (i *Integration) SyncSchemaHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Types   []string               `json:"types"`
		Filter  fibery.SyncFilter      `json:"filter"`
		Account integrationAccountInfo `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	requestedSchemas := make(map[string]map[string]fibery.Field)

	for _, t := range req.Types {
		requestedSchemas[t] = i.fiberyTypes[t]().Schema()
	}

	RespondWithJSON(w, http.StatusOK, requestedSchemas)
}

func (i *Integration) SyncDataHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		RequestedType     string                               `json:"requestedType"`
		OperationID       string                               `json:"operationId"`
		Types             []string                             `json:"types"`
		Filter            map[string]any                       `json:"filter"`
		Account           integrationAccountInfo               `json:"account"`
		LastSyncronizedAt string                               `json:"lastSynchronizedAt"`
		Pagination        fibery.NextPageConfig                `json:"pagination"`
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

	typePointer, ok := data.Types[params.RequestedType]
	if !ok {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type was not found: %s", params.RequestedType))
		return
	}

	req := data.Request{
		StartPosition:  startPosition,
		OperationId:    params.OperationID,
		RealmId:        params.Account.RealmID,
		LastSynced:     params.LastSyncronizedAt,
		RequestedType:  params.RequestedType,
		RequestedTypes: params.Types,
		CDCTypes:       reqTypesToCDCTypes(params.Types),
		Filter:         params.Filter,
		Cache:          i.cache,
		Group:          i.group,
		Token:          &params.Account.BearerToken,
	}

	responseBody := fibery.DataHandlerResponse{}
	allTypes := *typePointer

	switch datatype := allTypes.(type) {
	case data.DependentDataType:
		cacheKey := fmt.Sprintf("%s:%s", req.RealmId, req.RequestedType)
		cacheEntry, found := i.cache.Get(cacheKey)
		if params.LastSyncronizedAt == "" || !found {
			groupKey := fmt.Sprintf("%s:%s:%d", req.OperationId, datatype.Source.GetId(), req.StartPosition)
			res, err, _ := i.group.Do(groupKey, func() (interface{}, error) {
				data, err := datatype.Query(req)
				if err != nil {
					return nil, err
				}

				return data, nil
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to retrieve parent data: %w", err))
				return
			}

			parentRequest := res.(data.Response)

			idMap, err := datatype.TypeMapper(parentRequest.Data, datatype.SourceMapper)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to map ids: %w", err))
				return
			}

			if !found {
				idEntry := data.IdCache{
					OperationId: req.OperationId,
					Entries:     idMap,
				}
				err = req.Cache.Add(cacheKey, &idEntry, data.IdCacheLifetime)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to add cache entry: %w", err))
					return
				}
			} else {
				cacheEntry, ok := cacheEntry.(*data.IdCache)
				if !ok {
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to convert cache entry to IdCache"))
					return
				}
				cacheEntry.Mu.Lock()
				defer cacheEntry.Mu.Unlock()
				if !ok {
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to convert cache entry to IdCache"))
					return
				}
				if cacheEntry.OperationId == req.OperationId {
					for sourceId, sourceMap := range idMap {
						cacheEntry.Entries[sourceId] = sourceMap
					}
					req.Cache.Set(cacheKey, cacheEntry, data.IdCacheLifetime)
				} else {
					newCacheEntry := data.IdCache{
						OperationId: req.OperationId,
						Entries:     idMap,
					}
					req.Cache.Set(cacheKey, &newCacheEntry, data.IdCacheLifetime)
				}
			}
			items, err := datatype.QueryProcessor(parentRequest.Data, datatype.SchemaTransformer)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process data: %w", err))
				return
			}
			responseBody = fibery.DataHandlerResponse{
				Items: items,
				Pagination: fibery.Pagination{
					HasNext: parentRequest.MoreData,
					NextPageConfig: fibery.NextPageConfig{
						StartPosition: req.StartPosition + qbo.QueryPageSize,
					},
				},
				SynchronizationType: fibery.FullSync,
			}
		} else {
			groupKey := params.OperationID
			cacheEntry, ok := cacheEntry.(*data.IdCache)
			if !ok {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to convert cache entry to IdCache"))
				return
			}
			res, err, _ := i.group.Do(groupKey, func() (interface{}, error) {
				// get change data capture for all valid types
				return nil, nil
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get change data capture: %w", err))
				return
			}

			items, err := datatype.ChangeDataCaptureProcessor(res.(qbo.ChangeDataCapture), cacheEntry, datatype.SourceMapper, datatype.SchemaTransformer)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process change data capture: %w", err))
				return
			}

			responseBody = fibery.DataHandlerResponse{
				Items: items,
				Pagination: fibery.Pagination{
					HasNext: false,
					NextPageConfig: fibery.NextPageConfig{
						StartPosition: req.StartPosition + qbo.QueryPageSize,
					},
				},
				SynchronizationType: fibery.DeltaSync,
			}
		}
	default:
		if params.LastSyncronizedAt == "" {
			groupKey := fmt.Sprintf("%s:%s:%d", params.OperationID, params.RequestedType, params.StartPosition)
			res, err, _ := i.group.Do(groupKey, func() (interface{}, error) {
				// get query from datatype
				return nil, nil
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to retrieve data: %w", err))
				return
			}

			// do a for loop and convert all data to schema
			// set value on response body
			// set synctype to full
		} else {
			groupKey := params.OperationID
			res, err, _ := i.group.Do(groupKey, func() (interface{}, error) {
				// get change data capture for all valid types
				return nil, nil
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get change data capture: %w", err))
				return
			}

			// walking through the cdc response should be the same for all other types
			// create a function in datatype the takes a single corresponding datatype and converts it to a map[string]any matching the corresponding schema
			// append the map[string]any for each entity to the response body
			// set synctype to delta
		}
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
		Cache:          i.cache,
		Group:          i.group,
		Token:          &params.Account.BearerToken,
	}

	datatype := data.FiberyTypes[params.RequestedType]()
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

	fmt.Printf("DataHandler response for datatype: %s \n%s\n", params.RequestedType, qbo.FormatJSON(res))

	RespondWithJSON(w, http.StatusOK, res)
}
