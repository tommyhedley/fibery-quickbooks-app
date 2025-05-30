package main

import (
	"context"
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
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type QuickBooksAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmID string `json:"realmId,omitempty"`
	quickbooks.BearerToken
}

type SyncDataRequest struct {
	RequestedType     string                             `json:"requestedType"`
	OperationID       string                             `json:"operationId"`
	Types             []string                           `json:"types"`
	Schema            map[string]map[string]fibery.Field `json:"schema"`
	Filter            map[string]any                     `json:"filter"`
	Account           QuickBooksAccountInfo              `json:"account"`
	LastSyncronizedAt time.Time                          `json:"lastSynchronizedAt"`
	Pagination        fibery.NextPageConfig              `json:"pagination"`
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

	refreshNeeded := token.CheckExpiration(i.config.RefreshSecBeforeExpriation)

	if refreshNeeded {
		newToken, err := i.client.RefreshToken(token.RefreshToken)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to refresh token: %w", err))
			return
		}
		token = newToken
	}

	params := quickbooks.RequestParameters{
		Ctx:             r.Context(),
		RealmId:         reqBody.Fields.RealmID,
		Token:           token,
		WaitOnRateLimit: true,
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
	type requestBody struct {
		Types   []string              `json:"types"`
		Filter  fibery.SyncFilter     `json:"filter"`
		Account QuickBooksAccountInfo `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	requestedSchemas := make(map[string]map[string]fibery.Field)

	for _, typeId := range params.Types {
		storedType, ok := i.types.GetType(typeId)
		if !ok {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("type %s not found in registered types", typeId))
			return
		}
		requestedSchemas[typeId] = storedType.Schema()
	}

	RespondWithJSON(w, http.StatusOK, requestedSchemas)
}

func (i *Integration) SyncDataHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := SyncDataRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	if params.OperationID == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("operationId is required"))
		return
	}

	if params.RequestedType == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requestedType is required"))
		return
	}

	op, err := i.opManager.GetOrAddOperation(params, i.types, i.idStore, i.client, i.config.QuickBooks.PageSize)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("issue getting/creating operation: %w", err))
		return
	}

	if params.Pagination.StartPosition == 0 {
		params.Pagination.StartPosition = 1
	}

	requestedType, exists := i.types.GetType(params.RequestedType)
	if !exists {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requested type: %s does not exist in the TypeRegistry", params.RequestedType))
		return
	}

	resp, err := requestedType.GetData(i.client, op, params.Pagination, i.config.QuickBooks.PageSize)
	if err != nil {
		HandleRequestError(w, http.StatusInternalServerError, "data request failed", err)
		return
	}

	if op.IsComplete() {
		i.opManager.DeleteOperation(params.OperationID)
		slog.Debug(fmt.Sprintf("operation %s deleted", params.OperationID))
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

func (i *Integration) WebhookPreProcessHandler(w http.ResponseWriter, r *http.Request) {
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

	mac := hmac.New(sha256.New, []byte(i.config.QuickBooks.WebhookToken))
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
		Types   []string                           `json:"types"`
		Schema  map[string]map[string]fibery.Field `json:"schema"`
		Filter  map[string]any                     `json:"filter"`
		Account QuickBooksAccountInfo              `json:"account"`
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

	idCache, exists := i.idStore.GetOrCreateIdCache(params.Account.RealmID)
	if !exists {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("no idCache was found for realmId: %s, perform a full sync before enabling webhooks", params.Account.RealmID))
		return
	}

	activeWebhookTypesBySource := make(map[string][]Type)
	activeRelatedTypesByWebhook := make(map[string][]CDCType)

	for _, typeId := range params.Types {
		storedType, ok := i.types.GetType(typeId)
		if !ok {
			RespondWithError(w, http.StatusBadRequest, fmt.Errorf("type %s not found in registered types", typeId))
			return
		} else {
			switch whType := storedType.(type) {
			case UnionType:
				for _, unionType := range whType.UnionTypes() {
					switch ut := unionType.(type) {
					case WebhookType:
						activeWebhookTypesBySource[ut.Id()] = append(activeWebhookTypesBySource[ut.Id()], storedType)
						activeRelatedTypesByWebhook[ut.Id()] = append(activeRelatedTypesByWebhook[ut.Id()], ut.GetRelatedTypes()...)
					case WebhookDepType:
						activeWebhookTypesBySource[ut.SourceId()] = append(activeWebhookTypesBySource[ut.SourceId()], storedType)
					}
				}
			case WebhookType:
				activeWebhookTypesBySource[whType.Id()] = append(activeWebhookTypesBySource[whType.Id()], storedType)
				activeRelatedTypesByWebhook[whType.Id()] = append(activeRelatedTypesByWebhook[whType.Id()], whType.GetRelatedTypes()...)
			case WebhookDepType:
				activeWebhookTypesBySource[whType.SourceId()] = append(activeWebhookTypesBySource[whType.SourceId()], storedType)
			}
		}
	}

	attachableSources := GetAttachmentSources(params.Schema, "Files")
	for source := range attachableSources {
		slog.Debug(fmt.Sprintf("attachable source: %s", source))
	}

	activeAttachableSources := make(map[string]bool)

	queryEntities := map[string][]string{}
	deleteEntities := map[string][]string{}
	cdcTypes := map[string]CDCType{}

	var oldestChange time.Time

	for _, event := range params.Payload.EventNotifications {
		if event.RealmID != params.Account.RealmID {
			continue
		}
		for _, e := range event.DataChangeEvent.Entities {
			if _, exists := activeWebhookTypesBySource[e.Name]; exists {
				switch e.Operation {
				case "Create", "Update", "Emailed", "Void":
					queryEntities[e.Name] = append(queryEntities[e.Name], e.ID)
					if attachableSources[e.Name] {
						activeAttachableSources[e.Name] = true
					}
				case "Delete", "Merge":
					deleteEntities[e.Name] = append(deleteEntities[e.Name], e.ID)
				}
			}
			if relatedTypes, exists := activeRelatedTypesByWebhook[e.Name]; exists {
				for _, typ := range relatedTypes {
					cdcTypes[typ.Id()] = typ
				}
				if oldestChange.IsZero() || e.LastUpdated.Before(oldestChange) {
					oldestChange = e.LastUpdated
				}
			}
		}
	}

	for source := range activeAttachableSources {
		slog.Debug(fmt.Sprintf("active attachable source: %s", source))
	}

	resp := &fibery.WebhookTransformResponse{
		Data: make(fibery.WebhookData),
	}

	if len(deleteEntities) > 0 {
		for typeId, ids := range deleteEntities {
			for _, storedType := range activeWebhookTypesBySource[typeId] {
				switch whType := storedType.(type) {
				case WebhookType:
					whType.ProcessWebhookDeletions(ids, resp)
				case WebhookDepType:
					whType.ProcessWebhookDeletions(ids, resp, idCache)
				default:
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("invalid type found in activeTypesBySource array: %s", storedType.Id()))
				}
			}
		}
	}

	if len(queryEntities) > 0 {
		batchRequest := make([]quickbooks.BatchItemRequest, 0, len(queryEntities))

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

		for _, itemResponse := range batchResponse {
			if faultType := itemResponse.Fault.FaultType; faultType != "" {
				RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("batch request error: %s", faultType))
				return
			}

			for _, storedType := range activeWebhookTypesBySource[itemResponse.BID] {
				switch webhookType := storedType.(type) {
				case WebhookType:
					webhookType.ProcessWebhookUpdate(&itemResponse, resp)
				case WebhookDepType:
					webhookType.ProcessWebhookUpdate(&itemResponse, resp, idCache)
				default:
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("invalid type found in activeTypesBySource array: %s", storedType.Id()))
				}
			}
		}
	}

	numCDCType := len(cdcTypes)
	numAttachableSources := len(activeAttachableSources)

	slog.Debug(fmt.Sprintf("number of CDC types to query: %d", numCDCType))
	slog.Debug(fmt.Sprintf("number of Attachable sources: %d", numAttachableSources))

	if numCDCType > 0 || numAttachableSources > 0 {
		cdcQueryEntities := make([]string, 0, numCDCType+numAttachableSources)

		if numCDCType > 0 {
			for typeId := range cdcTypes {
				cdcQueryEntities = append(cdcQueryEntities, typeId)
			}
		}

		if numAttachableSources > 0 {
			cdcQueryEntities = append(cdcQueryEntities, "Attachable")
		}

		reqParams := quickbooks.RequestParameters{
			Ctx:     r.Context(),
			RealmId: params.Account.RealmID,
			Token:   &params.Account.BearerToken,
		}

		oldestChange = oldestChange.Add(time.Second * -5)

		cdc, err := i.client.ChangeDataCapture(reqParams, cdcQueryEntities, oldestChange)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to make cdc request: %w", err))
			return
		}

		if numCDCType > 0 {
			for typeId, cdcType := range cdcTypes {
				items, err := cdcType.processCDC(&cdc)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to process cdc data: %w", err))
					return
				}
				resp.Data[typeId] = items
			}
		}

		if numAttachableSources > 0 {
			attachables := quickbooks.CDCQueryExtractor[quickbooks.Attachable](&cdc)

			for _, attachable := range attachables {
				if len(attachable.AttachableRef) == 0 {
					continue
				}

				attachableURL := GenerateAttachablesURL(attachable)
				fmt.Println(attachableURL)

				for _, ref := range attachable.AttachableRef {
					if !activeAttachableSources[ref.EntityRef.Type] {
						continue
					}

					items := resp.Data[ref.EntityRef.Type]

					for i, item := range items {
						if item["id"] == ref.EntityRef.Value {
							URLs := resp.Data[ref.EntityRef.Type][i]["Files"].([]string)
							resp.Data[ref.EntityRef.Type][i]["Files"] = append(URLs, attachableURL)
						}
					}
				}
			}
		}
	}

	RespondWithJSON(w, http.StatusOK, resp)
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

	slog.Debug(fmt.Sprintf("auth scope: %s", i.config.QuickBooks.Scope))
	redirectURI, err := i.client.FindAuthorizationUrl(i.config.QuickBooks.Scope, reqBody.State, reqBody.CallbackURI)
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

func (Integration) ActionHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, nil)
}

func (i *Integration) SyncResourceHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Types   []string              `json:"types"`
		Filter  map[string]any        `json:"filter"`
		Account QuickBooksAccountInfo `json:"account"`
		Params  struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"params"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	switch req.Params.Type {
	case "attachable":
		requestParams := quickbooks.RequestParameters{
			Ctx:             context.Background(),
			RealmId:         req.Account.RealmID,
			Token:           &req.Account.BearerToken,
			WaitOnRateLimit: false,
		}

		downloadUrl, err := i.client.GetAttachableDownloadURL(requestParams, req.Params.Id)
		if err != nil {
			HandleRequestError(w, http.StatusInternalServerError, "data request failed", err)
			return
		}

		resp, err := http.Get(downloadUrl.String())
		if err != nil {
			RespondWithError(w, http.StatusBadGateway, fmt.Errorf("download request failed: %w", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			RespondWithError(w, resp.StatusCode, fmt.Errorf("unexpected status %d from QuickBooks", resp.StatusCode))
			return
		}

		if ct := resp.Header.Get("Content-Type"); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		if cd := resp.Header.Get("Content-Disposition"); cd != "" {
			w.Header().Set("Content-Disposition", cd)
		}
		if cl := resp.Header.Get("Content-Length"); cl != "" {
			w.Header().Set("Content-Length", cl)
		}

		if _, err := io.Copy(w, resp.Body); err != nil {
			fmt.Printf("warning: streaming error: %v\n", err)
		}
	default:
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("resouce type: %s is not implemented", req.Params.Type))
		return
	}
	RespondWithJSON(w, http.StatusOK, nil)
}
