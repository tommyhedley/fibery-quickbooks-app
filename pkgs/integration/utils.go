package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

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

func startPosition(page, pageSize int) int {
	return ((page - 1) * pageSize) + 1
}
