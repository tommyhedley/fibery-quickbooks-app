package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/data"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/fibery"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/qbo"
	"golang.org/x/sync/singleflight"
)

type Integration struct {
	cache      *cache.Cache
	group      *singleflight.Group
	client     *qbo.Client
	appConfig  fibery.AppConfig
	syncConfig fibery.SyncConfig
	types      map[string]*data.Type
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
		typePointer, ok := i.types[t]
		if !ok {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("type %s not found in registered types", t))
			return
		}
		datatype := *typePointer
		requestedSchemas[t] = datatype.GetSchema()
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

	lastSynced := time.Time{}
	if params.LastSyncronizedAt != "" {
		lastSynced, err = time.Parse(fibery.DateFormat, params.LastSyncronizedAt)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to convert lastSynchronizedAt value: %w", err))
		}
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
		LastSynced:     lastSynced,
		RequestedType:  params.RequestedType,
		RequestedTypes: params.Types,
		CDCTypes:       ConvertToCDCTypes(i.types, params.Types),
		Filter:         params.Filter,
		Cache:          i.cache,
		Group:          i.group,
		Token:          &params.Account.BearerToken,
	}

	var responseBody fibery.DataHandlerResponse
	types := *typePointer

	switch datatype := types.(type) {
	case data.DependentType:
		if cdctype, ok := datatype.(data.DepCDCQueryable); ok {
			cacheKey := fmt.Sprintf("%s:%s", req.RealmId, req.RequestedType)
			cacheEntry, found := i.cache.Get(cacheKey)
			if req.LastSynced.IsZero() || !found {
				groupKey := fmt.Sprintf("%s:%s:%d", req.OperationId, datatype.GetSourceId(), req.StartPosition)
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

				dataResponse := res.(data.Response)

				idMap, err := cdctype.MapType(dataResponse.Data)
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
				items, err := datatype.ProcessQuery(dataResponse.Data)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process data: %w", err))
					return
				}
				responseBody = fibery.DataHandlerResponse{
					Items: items,
					Pagination: fibery.Pagination{
						HasNext: dataResponse.MoreData,
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
					client, err := qbo.NewClient(req.RealmId, req.Token)
					if err != nil {
						return nil, fmt.Errorf("unable to create quickbooks client: %w", err)
					}

					data, err := client.ChangeDataCapture(req.CDCTypes, req.LastSynced)
					if err != nil {
						return nil, err
					}

					return data, nil
				})

				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get change data capture: %w", err))
					return
				}

				items, err := cdctype.ProcessCDC(res.(qbo.ChangeDataCapture), cacheEntry)
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
		} else {
			groupKey := fmt.Sprintf("%s:%s:%d", req.OperationId, datatype.GetSourceId(), req.StartPosition)
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

			dataResponse := res.(data.Response)

			items, err := datatype.ProcessQuery(dataResponse.Data)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process data: %w", err))
				return
			}
			responseBody = fibery.DataHandlerResponse{
				Items: items,
				Pagination: fibery.Pagination{
					HasNext: dataResponse.MoreData,
					NextPageConfig: fibery.NextPageConfig{
						StartPosition: req.StartPosition + qbo.QueryPageSize,
					},
				},
				SynchronizationType: fibery.FullSync,
			}
		}
	default:
		if datatype, ok := datatype.(data.CDCQueryable); !ok || req.LastSynced.IsZero() {
			groupKey := fmt.Sprintf("%s:%s:%d", req.OperationId, params.RequestedType, req.StartPosition)
			res, err, _ := i.group.Do(groupKey, func() (interface{}, error) {
				data, err := datatype.Query(req)
				if err != nil {
					return nil, err
				}
				return data, nil
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to retrieve data: %w", err))
				return
			}

			dataResponse := res.(data.Response)

			items, err := datatype.ProcessQuery(dataResponse.Data)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process data: %w", err))
				return
			}
			responseBody = fibery.DataHandlerResponse{
				Items: items,
				Pagination: fibery.Pagination{
					HasNext: dataResponse.MoreData,
					NextPageConfig: fibery.NextPageConfig{
						StartPosition: req.StartPosition + qbo.QueryPageSize,
					},
				},
				SynchronizationType: fibery.FullSync,
			}
		} else {
			groupKey := params.OperationID
			res, err, _ := i.group.Do(groupKey, func() (interface{}, error) {
				client, err := qbo.NewClient(req.RealmId, req.Token)
				if err != nil {
					return nil, fmt.Errorf("unable to create quickbooks client: %w", err)
				}

				data, err := client.ChangeDataCapture(req.CDCTypes, req.LastSynced)
				if err != nil {
					return nil, err
				}
				return data, nil
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get change data capture: %w", err))
				return
			}

			items, err := datatype.ProcessCDC(res.(qbo.ChangeDataCapture))
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
	}

	slog.Info(fmt.Sprintf("Getting data for %s", params.RequestedType))

	RespondWithJSON(w, http.StatusOK, responseBody)
}

func (Integration) WebhookInitHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Account integrationAccountInfo `json:"account"`
		Types   []string               `json:"types"`
		Filter  map[string]any         `json:"filter"`
		Webhook fibery.Webhook         `json:"webhook"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	webhookID := uuid.New().String()
	RespondWithJSON(w, http.StatusOK, fibery.Webhook{WebhookID: webhookID, WorkspaceID: params.Account.RealmID})
}

func (Integration) WebhookPreProcessHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		EventNotifications []struct {
			RealmID         string `json:"realmId"`
			DataChangeEvent struct {
				Entities []struct {
					ID          string    `json:"id"`
					Operation   string    `json:"operation"`
					Name        string    `json:"name"`
					LastUpdated time.Time `json:"lastUpdated"`
				} `json:"entities"`
			} `json:"dataChangeEvent"`
		} `json:"eventNotifications"`
	}

	type responseBody struct {
		Reply        map[string]string `json:"reply"`
		WorkspaceIds []string          `json:"workspaceIds"`
	}

	verifierToken := os.Getenv("WEBHOOK_TOKEN")
	if verifierToken == "" {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("missing verifier token"))
		return
	}

	intuitSignature := r.Header.Get("intuit-signature")
	if intuitSignature == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("missing intuit-signature header"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to read request body: %w", err))
		return
	}

	mac := hmac.New(sha256.New, []byte(verifierToken))
	mac.Write(bodyBytes)
	computedHash := mac.Sum(nil)

	decodedSignature, err := base64.StdEncoding.DecodeString(intuitSignature)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid base64 signature: %w", err))
		return
	}

	if !hmac.Equal(computedHash, decodedSignature) {
		RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("signature and payload do not match"))
		return
	}

	var params requestBody
	if err := json.Unmarshal(bodyBytes, &params); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	workspaceIDs := []string{}
	for _, event := range params.EventNotifications {
		workspaceIDs = append(workspaceIDs, event.RealmID)
	}

	RespondWithJSON(w, http.StatusOK, responseBody{
		Reply:        map[string]string{"challenge": "success"},
		WorkspaceIds: workspaceIDs,
	})
}

func (i *Integration) WebhookTransformHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Params struct {
			Connection                      string    `json:"connection"`
			XForwardedPort                  string    `json:"x-forwarded-port"`
			XForwardedPath                  string    `json:"x-forwarded-path"`
			XForwardedPrefix                string    `json:"x-forwarded-prefix"`
			XRealIP                         string    `json:"x-real-ip"`
			UserAgent                       string    `json:"user-agent"`
			ContentType                     string    `json:"content-type"`
			Accept                          string    `json:"accept"`
			IntuitSignature                 string    `json:"intuit-signature"`
			IntuitCreatedTime               time.Time `json:"intuit-created-time"`
			IntuitTID                       string    `json:"intuit-t-id"`
			IntuitNotificationSchemaVersion string    `json:"intuit-notification-schema-version"`
			AcceptEncoding                  string    `json:"accept-encoding"`
			Authorization                   string    `json:"authorization"`
		} `json:"params"`
		Types  []string `json:"types"`
		Filter struct {
		} `json:"filter"`
		Account integrationAccountInfo `json:"account"`
		Payload struct {
			EventNotifications []struct {
				RealmID         string `json:"realmId"`
				DataChangeEvent struct {
					Entities []struct {
						ID          string    `json:"id"`
						Operation   string    `json:"operation"`
						Name        string    `json:"name"`
						LastUpdated time.Time `json:"lastUpdated"`
					} `json:"entities"`
				} `json:"dataChangeEvent"`
			} `json:"eventNotifications"`
		} `json:"payload"`
	}

	type responseBody struct {
		Data map[string][]map[string]any `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	queryEntities := map[string][]string{}
	deleteEntities := map[string][]string{}

	for _, event := range params.Payload.EventNotifications {
		if event.RealmID != params.Account.RealmID {
			continue
		}
		for _, e := range event.DataChangeEvent.Entities {
			if _, ok := i.types[e.Name]; ok {
				switch e.Operation {
				case "Create", "Update", "Emailed", "Void":
					queryEntities[e.Name] = append(queryEntities[e.Name], e.ID)
				case "Delete", "Merge":
					deleteEntities[e.Name] = append(deleteEntities[e.Name], e.ID)
				}
			}
		}
	}

	// handle removing deleted entities from cache

	var response responseBody
	response.Data = map[string][]map[string]any{}

	for typ, ids := range deleteEntities {
		for _, id := range ids {
			response.Data[typ] = append(response.Data[typ], map[string]any{
				"id":           id,
				"__syncAction": fibery.REMOVE,
			})
		}
		if dependents, ok := data.SourceDependents[typ]; ok {
			for _, dependentPtr := range dependents {
				dependent := *dependentPtr
				cacheKey := fmt.Sprintf("%s:%s", params.Account.RealmID, dependent.GetId())
				if cacheEntry, found := i.cache.Get(cacheKey); found {
					cacheEntry, ok := cacheEntry.(*data.IdCache)
					if !ok {
						RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to convert cache entry to IdCache"))
						return
					}

					cacheEntry.Mu.Lock()
					defer cacheEntry.Mu.Unlock()
					for _, id := range ids {
						cachedIds := cacheEntry.Entries[id]
						for cachedId := range cachedIds {
							response.Data[dependent.GetId()] = append(response.Data[dependent.GetId()], map[string]any{
								"id":           cachedId,
								"__syncAction": fibery.REMOVE,
							})
						}
						delete(cacheEntry.Entries, id)
						if _, ok := cacheEntry.Entries[id]; !ok {
							fmt.Printf("cache entry for invoice %s deleted\n", id)
						}
					}
				}

			}
		}
	}

	batchRequest := []qbo.BatchItemRequest{}

	for typ, ids := range queryEntities {
		req := qbo.BatchItemRequest{
			BID:   typ,
			Query: fmt.Sprintf("SELECT * FROM %s WHERE Id IN ('%s')", typ, strings.Join(ids, "','")),
		}
		batchRequest = append(batchRequest, req)
	}

	client, err := qbo.NewClient(params.Account.RealmID, &params.Account.BearerToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create new client: %w", err))
		return
	}

	batchResponse, err := client.BatchRequest(batchRequest)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to make batch request: %w", err))
		return
	}

	for _, itemResponse := range batchResponse {
		datatype := *i.types[itemResponse.BID]
		whDatatype, ok := datatype.(data.WHQueryable)
		if !ok {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to convert datatype to WHQueryable"))
			return
		}
		err := whDatatype.ProcessWHBatch(itemResponse, &response.Data, i.cache, params.Account.RealmID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process webhook batch: %w", err))
			return
		}
	}

	fmt.Println(data.FormatJSON(response))

	RespondWithJSON(w, http.StatusOK, response)
}

func (i *Integration) Oauth2AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		CallbackURI string `json:"callback_uri"`
		State       string `json:"state"`
	}
	type responseBody struct {
		RedirectURI string `json:"redirect_uri"`
	}
	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	client, err := qbo.NewClient("", nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error creating new client: %w", err))
		return
	}

	redirectURI, err := client.FindAuthorizationUrl(os.Getenv("SCOPE"), reqBody.State, reqBody.CallbackURI)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error generating redirect uri: %w", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, responseBody{
		RedirectURI: redirectURI,
	})
}

func (i *Integration) Oauth2TokenHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Fields struct {
			ID          string `json:"_id"`
			App         string `json:"app"`
			Owner       string `json:"owner"`
			ExpireOn    string `json:"expireOn"`
			CallbackURI string `json:"callback_uri"`
			State       string `json:"state"`
		} `json:"fields"`
		Code    string `json:"code"`
		State   string `json:"state"`
		RealmID string `json:"realmId"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	realmId := reqBody.RealmID
	if realmId == "" {
		mode := os.Getenv("MODE")
		switch mode {
		case "production":
			realmId = os.Getenv("REALM_ID_PRODUCTION")
		case "sandbox":
			realmId = os.Getenv("REALM_ID_SANDBOX")
		default:
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("invalid mode: %s", mode))
		}
	}
	client, err := qbo.NewClient(realmId, nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create new client: %w", err))
		return
	}

	token, err := client.RetrieveBearerToken(reqBody.Code, reqBody.Fields.CallbackURI)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to retreive bearer token: %w", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, integrationAccountInfo{
		RealmID:     realmId,
		BearerToken: *token,
	})
}

func (Integration) LogoHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./logo.svg")
	if err != nil {
		http.Error(w, "Unable to open SVG file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	svgData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read SVG file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")

	w.Write(svgData)
}

func (Integration) SyncFilterValidateHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, nil)
}
