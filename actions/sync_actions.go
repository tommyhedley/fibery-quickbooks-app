package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tommyhedley/fibery/fibery-qbo-integration/internal/utils"
)

type timesheet struct {
	TimesheetID json.Number
	UserID      json.Number `json:"user_id" type:"string"`
	JobcodeId   json.Number `json:"jobcode_id" type:"string"`
	Type        string      `json:"type"`
}

type manualTimesheet struct {
	timesheet
	Duration int    `json:"duration"`
	Date     string `json:"date"`
}

type regularTimesheet struct {
	timesheet
	Start string `json:"start"`
	End   string `json:"end"`
}

func SyncActionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("type") {
	case "manual-timesheet":
		decoder := json.NewDecoder(r.Body)
		params := []manualTimesheet{}
		err := decoder.Decode(&params)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to decode request parameters: %w", err))
			return
		}

	case "regular-timehsheet":
		decoder := json.NewDecoder(r.Body)
		params := []regularTimesheet{}
		err := decoder.Decode(&params)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to decode request parameters: %w", err))
			return
		}

	default:
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid request type"))
		return
	}
}

func SyncActionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		syncActionToken := os.Getenv("SYNC_ACTION_API_KEY")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("missing auth header"))
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("invalid auth header format"))
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if token == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("missing auth token"))
			return
		}
		if token != syncActionToken {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("invalid auth token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
