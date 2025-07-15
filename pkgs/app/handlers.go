package app

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
	"time"

	"github.com/google/uuid"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

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

	refreshNeeded := token.CheckExpiration(int(i.config.TokenRefreshWindow.Seconds()))

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
		RealmId:         reqBody.Fields.RealmId,
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
		RealmId:     reqBody.Fields.RealmId,
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
		storedType, ok := i.types.Get(typeId)
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
	req := SyncRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	if req.OperationId == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("operationId is required"))
		return
	}

	if req.RequestedType == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("requestedType is required"))
		return
	}

	if req.Pagination.Page == 0 {
		req.Pagination.Page = 1
	}

	op, err := i.opManager.GetOrAddOperation(req, i)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("issue getting/creating operation: %w", err))
		return
	}

	err = op.SubmitRequest(req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	key := ResponseChannelKey(req.RequestedType, req.Pagination.Page)

	slog.Debug(fmt.Sprintf("request submitted & key created"))

	ch, err := op.GetChannel(key)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	slog.Debug("channel found")

	select {
	case resp, ok := <-ch:
		if !ok {
			RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("channel prematurely closed for key: %s", key))
			return
		}

		if resp.Error != nil {
			RespondWithError(w, http.StatusInternalServerError, resp.Error)
			return
		}

		RespondWithJSON(w, http.StatusOK, resp.DataHandlerResponse)
	case <-op.ctx.Done():
		RespondWithError(w, http.StatusGatewayTimeout, fmt.Errorf("operation %s cancelled or timed out", req.OperationId))
		return
	}
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
	RespondWithJSON(w, http.StatusOK, fibery.Webhook{WebhookId: webhookID, WorkspaceId: params.Account.RealmId})
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
	decoder := json.NewDecoder(r.Body)
	req := WebhookRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request body: %w", err))
		return
	}

	wg, err := buildWebhookGroup(req, i, time.Duration(5*time.Second))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error building webhookGroup: %w", err))
		return
	}

	err = wg.fetchAll(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error fetching webhook data: %w", err))
		return
	}

	items, err := wg.process()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("error processing webhookGroup data: %w", err))
		return
	}

	resp := fibery.WebhookTransformResponse{
		Data: make(fibery.WebhookData, len(items)),
	}

	resp.Data = items

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
		RealmId:     realmId,
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
			Ctx:             r.Context(),
			RealmId:         req.Account.RealmId,
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
