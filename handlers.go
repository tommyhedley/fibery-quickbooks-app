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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tommyhedley/fibery-quickbooks-app/data"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type QuickBooksAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmID string `json:"realmId,omitempty"`
	quickbooks.BearerToken
}

func (i *Integration) AppConfigHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, i.appConfig)
}

func (i *Integration) AccountValidateHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Id     string `json:"id"`
		Fields struct {
			QuickBooksAccountInfo
		} `json:"fields"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	token := &reqBody.Fields.BearerToken

	secBeforeExp, err := strconv.Atoi(os.Getenv("TOKEN_REFRESH_BEFORE_EXPIRATION"))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("invalid TOKEN_REFRESH_BEFORE_EXPIRATION value: %w", err))
		return
	}

	refreshNeeded := token.CheckExpiration(secBeforeExp)

	if refreshNeeded {
		newToken, err := i.client.RefreshToken(token.RefreshToken)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to refresh token: %w", err))
			return
		}
		token = newToken
	}

	params := quickbooks.RequestParameters{
		Ctx:     r.Context(),
		RealmId: reqBody.Fields.RealmID,
		Token:   token,
	}

	info, err := i.client.FindCompanyInfo(params)
	if err != nil {
		HandleRequestError(w, http.StatusInternalServerError, "unable to find company info", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, QuickBooksAccountInfo{
		Name:        info.CompanyName,
		RealmID:     reqBody.Fields.RealmID,
		BearerToken: *token,
	},
	)
}

func (i *Integration) SyncConfigHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, i.syncConfig)
}

func (i *Integration) SyncSchemaHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Types   []string              `json:"types"`
		Filter  fibery.SyncFilter     `json:"filter"`
		Account QuickBooksAccountInfo `json:"account"`
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
		typePointer, ok := data.Types.All[t]
		if !ok {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("type %s not found in registered types", t))
			return
		}
		datatype := *typePointer
		requestedSchemas[t] = datatype.Schema()
	}

	RespondWithJSON(w, http.StatusOK, requestedSchemas)
}

func (i *Integration) SyncDataHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		RequestedType     string                               `json:"requestedType"`
		OperationID       string                               `json:"operationId"`
		Types             []string                             `json:"types"`
		Filter            map[string]any                       `json:"filter"`
		Account           QuickBooksAccountInfo                `json:"account"`
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

	if params.OperationID == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("operationId is required"))
		return
	}

	if params.RequestedType == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requestedType is required"))
		return
	}

	lastSynced := time.Time{}
	if params.LastSyncronizedAt != "" {
		lastSynced, err = time.Parse(fibery.DateFormat, params.LastSyncronizedAt)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to convert lastSynchronizedAt value: %w", err))
			return
		}
	}

	opCache, exists := i.dataCache.Get(params.OperationID)
	if !exists {
		slog.Debug("opCache does not exist")
		requestedTypes := make(map[string]bool)
		syncTypes := make(map[string]fibery.SyncType)
		deltaRequestTypes := []string{}
		for _, rt := range params.Types {
			requestedTypes[rt] = false
			storedType := (*data.Types.All[rt])
			if !lastSynced.IsZero() {
				switch storedType.(type) {
				case data.CDCQueryable:
					syncTypes[rt] = fibery.Delta
					deltaRequestTypes = append(deltaRequestTypes, storedType.Id())
				case data.DepCDCQueryable:
					if _, ok := i.idCache.Get(params.Account.RealmID, storedType.Id()); ok {
						syncTypes[rt] = fibery.Delta
						deltaRequestTypes = append(deltaRequestTypes, storedType.Id())
					} else {
						syncTypes[rt] = fibery.Full
					}
				}
			} else {
				syncTypes[rt] = fibery.Full
			}
		}

		opCache = data.NewOperationCache(requestedTypes, syncTypes)
		i.dataCache.Set(params.OperationID, opCache)

		if len(deltaRequestTypes) > 0 {
			go func() {
				cdcParams := quickbooks.RequestParameters{
					Ctx:     r.Context(),
					RealmId: params.Account.RealmID,
					Token:   &params.Account.BearerToken,
				}

				cdc, err := i.client.ChangeDataCapture(cdcParams, deltaRequestTypes, lastSynced)
				if err != nil {
					slog.Error("CDC request failed", "error", err)
				} else {
					opCache.Lock()
					opCache.Results["cdc"] = &data.CacheEntry{
						Data: make(chan any, 1),
					}
					opCache.Results["cdc"].Data <- cdc
					opCache.Unlock()
				}

			}()
		}
	}

	req := data.Request{
		Filter:        params.Filter,
		LastSynced:    lastSynced,
		Ctx:           r.Context(),
		Types:         params.Types,
		RealmId:       params.Account.RealmID,
		StartPosition: startPosition,
		PageSize:      quickbooks.QueryPageSize,
		OpCache:       opCache,
		IdCache:       i.idCache,
		Client:        i.client,
		Token:         &params.Account.BearerToken,
	}

	requestType, ok := data.Types.All[params.RequestedType]
	if !ok {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requestedType was not found in i.types"))
		return
	}

	slog.Debug("getting data")

	resp, err := (*requestType).GetData(req)
	if err != nil {
		HandleRequestError(w, http.StatusInternalServerError, "data request failed", err)
		return
	}

	slog.Debug("data request completed")

	opCache.Lock()
	allComplete := true
	for _, complete := range opCache.RequestedTypes {
		if !complete {
			allComplete = false
			break
		}
	}
	opCache.Unlock()

	if allComplete {
		i.dataCache.Delete(params.OperationID)
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (Integration) WebhookInitHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Account QuickBooksAccountInfo `json:"account"`
		Types   []string              `json:"types"`
		Filter  map[string]any        `json:"filter"`
		Webhook fibery.Webhook        `json:"webhook"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	webhookID := uuid.New().String()
	RespondWithJSON(w, http.StatusOK, fibery.Webhook{WebhookId: webhookID, WorkspaceId: params.Account.RealmID})
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
		Types   []string              `json:"types"`
		Filter  map[string]any        `json:"filter"`
		Account QuickBooksAccountInfo `json:"account"`
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
			if _, ok := data.Types.All[e.Name]; ok {
				switch e.Operation {
				case "Create", "Update", "Emailed", "Void":
					queryEntities[e.Name] = append(queryEntities[e.Name], e.ID)
				case "Delete", "Merge":
					deleteEntities[e.Name] = append(deleteEntities[e.Name], e.ID)
				}
			}
		}
	}

	response := fibery.WebhookTransformResponse{
		Data: make(fibery.WebhookData),
	}

	for typeId, ids := range deleteEntities {
		for _, id := range ids {
			response.Data[typeId] = append(response.Data[typeId], map[string]any{
				"id":           id,
				"__syncAction": fibery.REMOVE,
			})
		}
		if depTypes, ok := data.Types.DepWHReceivable[typeId]; ok {
			for _, depPtr := range depTypes {
				for _, selectedType := range params.Types {
					if depType := *depPtr; depType.Id() == selectedType {
						idSet, exists := i.idCache.Get(params.Account.RealmID, depType.Id())
						if !exists {
							RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("no idCache exists for: %s, please full sync to build idCache", depType.Id()))
							return
						}
						for _, parentId := range ids {
							if depIds, exists := idSet[parentId]; exists {
								for _, depId := range depIds {
									response.Data[depType.Id()] = append(response.Data[depType.Id()], map[string]any{
										"id":           depId,
										"__syncAction": fibery.REMOVE,
									})
								}
								i.idCache.RemoveSource(params.Account.RealmID, depType.Id(), parentId)
							}
						}
					}
					continue
				}
			}
		}
	}

	batchRequest := []quickbooks.BatchItemRequest{}

	for typ, ids := range queryEntities {
		req := quickbooks.BatchItemRequest{
			BID:   typ,
			Query: fmt.Sprintf("SELECT * FROM %s WHERE Id IN ('%s')", typ, strings.Join(ids, "','")),
		}
		batchRequest = append(batchRequest, req)
	}

	reqParams := quickbooks.RequestParameters{
		Ctx:     r.Context(),
		RealmId: params.Account.RealmID,
		Token:   &params.Account.BearerToken,
	}

	batchResponse, err := i.client.BatchRequest(reqParams, batchRequest)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to make batch request: %w", err))
		return
	}

	req := data.Request{
		Filter:  params.Filter,
		Ctx:     r.Context(),
		Types:   params.Types,
		RealmId: params.Account.RealmID,
		IdCache: i.idCache,
		Client:  i.client,
		Token:   &params.Account.BearerToken,
	}

	for _, itemResponse := range batchResponse {
		if faultType := itemResponse.Fault.FaultType; faultType != "" {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("batch request error: %s", faultType))
			return
		}

		datatype, ok := (*data.Types.All[itemResponse.BID]).(data.WHReceivable)
		if !ok {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("requested datatype does not implement WHQueryable, please review you datatype implementations"))
			return
		}

		err = datatype.ProcessWH(req, &itemResponse, &response.Data)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process webhook batch: %w", err))
			return
		}
	}

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

	redirectURI, err := i.client.FindAuthorizationUrl(os.Getenv("SCOPE"), reqBody.State, reqBody.CallbackURI)
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

	token, err := i.client.RetrieveBearerToken(reqBody.Code, reqBody.Fields.CallbackURI)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to retreive bearer token: %w", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, QuickBooksAccountInfo{
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
