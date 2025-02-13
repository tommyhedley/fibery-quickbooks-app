package main

import (
	"encoding/json"
	"fmt"
	"github.com/tommyhedley/fibery-quickbooks-app/data"
	"github.com/tommyhedley/quickbooks-go"
	"log/slog"
	"net/http"
	"os"
)

func RespondWithError(w http.ResponseWriter, code int, err error) {
	if code >= 500 {
		slog.Error(err.Error(), "StatusCode", code)
	}
	w.Header().Set("Content-Type", "application/json")
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: err.Error(),
	})
}

func RespondWithRequestLimit(w http.ResponseWriter, code int, err error) {
	if code >= 500 {
		slog.Error(err.Error(), "StatusCode", code)
	}
	type errorResponse struct {
		Error    string `json:"error"`
		TryLater bool   `json:"tryLater"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error:    err.Error(),
		TryLater: true,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Error marshalling json", "StatusCode", 500)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func ConvertToCDCTypes(types map[string]*data.Type, requestedTypes []string) []string {
	typeSet := make(map[string]struct{})
	var CDCTypes []string

	for _, reqType := range requestedTypes {
		if typePointer, ok := types[reqType]; ok {
			datatype := *typePointer
			if cdcType, ok := datatype.(data.CDCQueryable); ok {
				if _, exists := typeSet[cdcType.GetId()]; !exists {
					typeSet[cdcType.GetId()] = struct{}{}
					CDCTypes = append(CDCTypes, cdcType.GetId())
				}
			}
			if depCDCType, ok := datatype.(data.DepCDCQueryable); ok {
				if _, exists := typeSet[depCDCType.GetSourceId()]; !exists {
					typeSet[depCDCType.GetSourceId()] = struct{}{}
					CDCTypes = append(CDCTypes, depCDCType.GetSourceId())
				}
			}
		}
	}

	return CDCTypes
}

func NewClientRequest(discovery *quickbooks.DiscoveryAPI, token *quickbooks.BearerToken, realmId string) (quickbooks.ClientRequest, error) {
	clientRequest := quickbooks.ClientRequest{
		DiscoveryAPI: discovery,
		Token:        token,
		RealmId:      realmId,
	}
	switch os.Getenv("MODE") {
	case "production":
		clientRequest.ClientId = os.Getenv("OAUTH_CLIENT_ID_PRODUCTION")
		clientRequest.ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET_PRODUCTION")
		clientRequest.Endpoint = os.Getenv("PRODUCTION_ENDPOINT")
	case "sandbox":
		clientRequest.ClientId = os.Getenv("OAUTH_CLIENT_ID_SANDBOX")
		clientRequest.ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET_SANDBOX")
		clientRequest.Endpoint = os.Getenv("SANDBOX_ENDPOINT")
	default:
		return quickbooks.ClientRequest{}, fmt.Errorf("invalid MODE setting: %s", os.Getenv("MODE"))
	}

	return clientRequest, nil
}
