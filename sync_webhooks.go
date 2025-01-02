package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/qbo"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type webhook struct {
		WebhookID   string `json:"id"`
		WorkspaceID string `json:"workspaceId"`
	}
	type requestBody struct {
		Account qbo.FiberyAccountInfo `json:"account"`
		Types   []string              `json:"types"`
		Filter  map[string]any        `json:"filter"`
		Webhook webhook               `json:"webhook"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	webhookID := uuid.New().String()
	RespondWithJSON(w, http.StatusOK, webhook{WebhookID: webhookID, WorkspaceID: params.Account.RealmID})
}

func PreProcessHandler(w http.ResponseWriter, r *http.Request) {
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

func TransformHandler(w http.ResponseWriter, r *http.Request) {
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
		Account qbo.FiberyAccountInfo `json:"account"`
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
		// Only process notifications for the matching realm
		if event.RealmID != params.Account.RealmID {
			continue
		}
		for _, e := range event.DataChangeEvent.Entities {
			if _, ok := qbo.Types[e.Name]; ok {
				switch e.Operation {
				case "Create", "Update", "Emailed", "Void":
					queryEntities[e.Name] = append(queryEntities[e.Name], e.ID)
					// If needed, skip old timestamps here
				case "Delete", "Merge":
					deleteEntities[e.Name] = append(deleteEntities[e.Name], e.ID)
				}
			}
		}
	}

	result := responseBody{
		Data: map[string][]map[string]any{},
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

	for _, item := range batchResponse {
	}

	// Handle creates/updates in parallel (one goroutine per type)
	var wg sync.WaitGroup
	for typ, ids := range queryEntities {
		wg.Add(1)
		go func(t string, IDs []string) {
			defer wg.Done()

			// Look up the QBO type handler
			dt := qbo.Types[t]
			if dt == nil {
				return
			}

			// Example: fetch or transform data for these IDs
			items, err := fetchAndTransform(realmID, dt, IDs)
			if err == nil {
				result.Data[t] = append(result.Data[t], items...)
			}
		}(typ, ids)
	}

	// Handle deletes/merges (can be done immediately or in parallel)
	for typ, ids := range deleteEntities {
		// Remove or mark deletions in your data store
		// Possibly use your ID cache
		for _, id := range ids {
			dt := qbo.Types[typ]
			if dt == nil {
				continue
			}
			// Example: just store a “deleted” entry
			delItem := map[string]any{
				"id":           id,
				"__syncAction": "REMOVE",
			}
			result.Data[typ] = append(result.Data[typ], delItem)
		}
	}

	wg.Wait()

	RespondWithJSON(w, http.StatusOK, struct{}{})
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, struct{}{})
}
