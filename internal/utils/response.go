package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type RequestError struct {
	StatusCode int
	Err        error
	TryLater   bool
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}

func NewRequestError(statusCode int, err error, tryLater bool) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		TryLater:   tryLater,
		Err:        err,
	}
}

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

func RespondWithTryLater(w http.ResponseWriter, code int, err error) {
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
