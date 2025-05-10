package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
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

func GetAttachmentSources(schema map[string]map[string]fibery.Field, attachablesFieldName string) map[string]bool {
	typeMap := map[string]bool{}
	for typeId, fields := range schema {
		for _, field := range fields {
			if field.SubType == fibery.File {
				typeMap[typeId] = true
			}
		}
	}
	return typeMap
}

func GenerateAttachablesURL(a quickbooks.Attachable) string {
	return fmt.Sprintf("app://resource?type=%s&id=%s", "attachable", a.Id)
}
