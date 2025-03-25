package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/tommyhedley/quickbooks-go"
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

func RespondWithRateLimit(w http.ResponseWriter, code int, err error) {
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

func HandleRequestError(w http.ResponseWriter, code int, errMsg string, err error) {
	var responseError error
	if errMsg == "" {
		responseError = err
	} else {
		responseError = fmt.Errorf("%s: %w", errMsg, err)
	}

	var rateLimitError *quickbooks.RateLimitError
	if errors.As(err, &rateLimitError) {
		RespondWithRateLimit(w, http.StatusTooManyRequests, responseError)
	} else {
		RespondWithError(w, code, responseError)
	}
}

func NewClientRequest(discovery *quickbooks.DiscoveryAPI, client *http.Client) (quickbooks.ClientRequest, error) {
	clientRequest := quickbooks.ClientRequest{
		Client:       client,
		DiscoveryAPI: discovery,
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
