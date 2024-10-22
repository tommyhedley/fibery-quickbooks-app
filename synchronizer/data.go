package synchronizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

type SyncType string

const (
	DeltaSync SyncType = "delta"
	FullSync  SyncType = "full"
)

type DeltaSyncAction int

const (
	GetNothing DeltaSyncAction = iota
	GetItems
	GetDeleted
	GetBoth
)

func DataHandler(w http.ResponseWriter, r *http.Request) {
	type nextPageConfig struct {
		Page int `json:"page"`
		Type int `json:"type"`
	}
	type pagination struct {
		HasNext        bool           `json:"hasNext"`
		NextPageConfig nextPageConfig `json:"nextPageConfig"`
	}
	type parameters struct {
		RequestedType string         `json:"requestedType"`
		Types         []string       `json:"types"`
		Filter        map[string]any `json:"filter"`
		Account       struct {
			AccessToken string `json:"access_token"`
		} `json:"account"`
		LastSyncronized string                               `json:"lastSynchronizedAt"`
		Pagination      nextPageConfig                       `json:"pagination"`
		Schema          map[string]map[string]map[string]any `json:"schema"`
	}
	type response[T any] struct {
		Items               []T        `json:"items"`
		Pagination          pagination `json:"pagination"`
		SynchronizationType string     `json:"synchronizationType"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	var lastSyncronized string

	if params.LastSyncronized != "" {
		lastSyncronizedTime, err := time.Parse(time.RFC3339, params.LastSyncronized)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to parse last sync time: %w", err))
			return
		}
		lastSyncronized = lastSyncronizedTime.Format("2006-01-02T15:04:05-07:00")
	}

	sync := DeltaSync
	if lastSyncronized == "" {
		sync = FullSync
	}

	var page int

	if params.Pagination.Page == 0 {
		page = 1
	} else {
		page = params.Pagination.Page
	}

	switch params.RequestedType {
	case "group":
		type groupRequest struct {
			Active           string `url:"active,omitempty"`
			Page             int    `url:"page"`
			SupplementalData string `url:"supplemental_data"`
			ModifiedSince    string `url:"modified_since,omitempty"`
		}

		groupReq := groupRequest{
			Active:           "both",
			Page:             page,
			SupplementalData: "no",
			ModifiedSince:    lastSyncronized,
		}

		groups, more, requestError := utils.APIRequest[groupRequest, any, GroupData](&groupReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/groups", params.Account.AccessToken, "groups")
		if requestError.Err != nil {
			if requestError.TryLater {
				utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
				return
			}
			utils.RespondWithError(w, requestError.StatusCode, fmt.Errorf("error with group data request: %w", requestError.Err))
			return
		}

		if sync != FullSync {
			for i := range groups {
				groups[i].SyncAction = "SET"
			}
		}

		resp := response[GroupData]{
			Items: groups,
			Pagination: pagination{
				HasNext: more,
				NextPageConfig: nextPageConfig{
					Page: page + 1,
				},
			},
			SynchronizationType: string(sync),
		}

		utils.RespondWithJSON(w, http.StatusOK, resp)
		return
	case "user":
		type userRequest struct {
			Active           string `url:"active"`
			Page             int    `url:"page"`
			SupplementalData string `url:"supplemental_data"`
			ModifiedSince    string `url:"modified_since,omitempty"`
		}

		userReq := userRequest{
			Active:           "both",
			Page:             page,
			SupplementalData: "no",
			ModifiedSince:    lastSyncronized,
		}

		users, more, requestError := utils.APIRequest[userRequest, any, UserData](&userReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/users", params.Account.AccessToken, "users")
		if requestError.Err != nil {
			if requestError.TryLater {
				utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
				return
			}
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with user data request: %w", requestError.Err))
			return
		}

		var items = []map[string]any{}

		for _, user := range users {
			var name strings.Builder
			if user.FirstName != "" {
				name.WriteString(user.FirstName)
				if user.LastName != "" {
					name.WriteString(" " + user.LastName)
				}
			} else if user.LastName != "" {
				name.WriteString(user.LastName)
			}
			item := map[string]any{
				"id":           user.ID.String(),
				"name":         name.String(),
				"display_name": user.DisplayName,
				"first_name":   user.FirstName,
				"last_name":    user.LastName,
				"active":       user.Active,
				"last_active":  user.LastActive,
				"group_id":     user.GroupID.String(),
				"email":        user.Email,
			}

			if customfields, ok := user.CustomFields.(map[string]any); ok {
				for key, value := range customfields {
					item[key] = value
				}
			}
			if sync != FullSync {
				item["__syncAction"] = "SET"
			}
			items = append(items, item)
		}

		resp := response[map[string]any]{
			Items: items,
			Pagination: pagination{
				HasNext: more,
				NextPageConfig: nextPageConfig{
					Page: page + 1,
				},
			},
			SynchronizationType: string(sync),
		}

		utils.RespondWithJSON(w, http.StatusOK, resp)
		return
	case "timesheet":
		var timesheetStartString string
		otc := "no"

		includeOTC, ok := params.Filter["includeOTC"]
		if !ok {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("filter key was not found in the request: includeOTC"))
			return
		}

		startVal, ok := params.Filter["timesheetStart"]
		if !ok {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("filter key was not found in the request: timesheetStart"))
			return
		}

		if includeOTC != nil {
			otcBool, ok := includeOTC.(bool)
			if !ok {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("filter value was the not correct type: includeOTC\n valid type: bool\n returned type: %T", includeOTC))
				return
			}
			if otcBool {
				otc = "both"
			}
		}

		if startVal == nil {
			timesheetStart := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
			timesheetStartString = timesheetStart.Format("2006-01-02")
		} else {
			startString, ok := startVal.(string)
			if !ok {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("filter value was the not correct type: timesheetStart\n valid type: string\n returned type: %T", startVal))
				return
			}
			timesheetStart, err := time.Parse(time.RFC3339, startString)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("unable to parse timesheetStart string value: %w", err))
				return
			}
			timesheetStartString = timesheetStart.Format("2006-01-02")
		}

		if sync == FullSync {
			type timesheetFullRequest struct {
				StartDate        string `url:"start_date"`
				Page             int    `url:"page"`
				SupplementalData string `url:"supplemental_data"`
				OnTheClock       string `url:"on_the_clock"`
			}

			timesheetFullReq := timesheetFullRequest{
				StartDate:        timesheetStartString,
				Page:             params.Pagination.Page,
				SupplementalData: "no",
				OnTheClock:       otc,
			}

			timesheets, more, requestError := utils.APIRequest[timesheetFullRequest, any, TimesheetData](&timesheetFullReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/timesheets", params.Account.AccessToken, "timesheets")
			if requestError.Err != nil {
				if requestError.TryLater {
					utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
					return
				}
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with timesheet data full request: %w", requestError.Err))
				return
			}

			var items = []map[string]any{}

			for _, timesheet := range timesheets {
				locked := false
				lockedNum, err := timesheet.Locked.Int64()
				if err != nil {
					utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error decoding locked value into int: %w", err))
					return
				}
				if lockedNum > 0 {
					locked = true
				}
				duration := time.Duration(timesheet.DurationSeconds) * time.Second
				item := map[string]any{
					"id":                 timesheet.ID.String(),
					"name":               timesheet.Notes,
					"user_id":            timesheet.UserID.String(),
					"created_by_user_id": timesheet.CreatedByUserID.String(),
					"jobcode_id":         timesheet.JobcodeID.String(),
					"locked":             locked,
					"last_modified":      timesheet.LastModified,
					"type":               timesheet.Type,
					"start":              timesheet.Start,
					"end":                timesheet.End,
					"date":               timesheet.Date,
					"duration":           duration.Seconds(),
					"duration_minutes":   duration.Minutes(),
					"duration_hours":     duration.Hours(),
					"on_the_clock":       timesheet.OnTheClock,
				}
				if customfields, ok := timesheet.CustomFields.(map[string]any); ok {
					for key, value := range customfields {
						item[key] = value
					}
				}
				items = append(items, item)
			}

			resp := response[map[string]any]{
				Items: items,
				Pagination: pagination{
					HasNext: more,
					NextPageConfig: nextPageConfig{
						Page: page + 1,
					},
				},
				SynchronizationType: string(FullSync),
			}

			utils.RespondWithJSON(w, http.StatusOK, resp)
			return
		} else {
			type timesheetDeltaRequest struct {
				ModifiedSince    string `url:"modified_since"`
				Page             int    `url:"page"`
				SupplementalData string `url:"supplemental_data"`
				OnTheClock       string `url:"on_the_clock"`
			}

			timesheetDeltaReq := timesheetDeltaRequest{
				ModifiedSince:    lastSyncronized,
				Page:             page,
				SupplementalData: "no",
				OnTheClock:       otc,
			}

			type timesheetDeletedRequest struct {
				ModifiedSince    string `url:"modified_since"`
				Page             int    `url:"page"`
				SupplementalData string `url:"supplemental_data"`
			}

			timesheetDeletedReq := timesheetDeletedRequest{
				ModifiedSince:    lastSyncronized,
				Page:             page,
				SupplementalData: "no",
			}

			var timesheets []TimesheetData
			var items []map[string]any
			var deletedTimesheets []TimesheetData
			var deletedItems []map[string]any
			var moreTimesheets, moreDeletedTimesheets bool
			var requestError *utils.RequestError

			if params.Pagination.Type == int(GetNothing) || params.Pagination.Type == int(GetBoth) {
				timesheets, moreTimesheets, requestError = utils.APIRequest[timesheetDeltaRequest, any, TimesheetData](&timesheetDeltaReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/timesheets", params.Account.AccessToken, "timesheets")
				if requestError.Err != nil {
					if requestError.TryLater {
						utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
						return
					}
					utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with timesheet data delta request: %w", requestError.Err))
					return
				}

				for _, timesheet := range timesheets {
					locked := false
					lockedNum, err := timesheet.Locked.Int64()
					if err != nil {
						utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error decoding locked value into int: %w", err))
						return
					}
					if lockedNum > 0 {
						locked = true
					}
					duration := time.Duration(timesheet.DurationSeconds) * time.Second
					item := map[string]any{
						"id":                 timesheet.ID.String(),
						"name":               timesheet.Notes,
						"user_id":            timesheet.UserID.String(),
						"created_by_user_id": timesheet.CreatedByUserID.String(),
						"jobcode_id":         timesheet.JobcodeID.String(),
						"locked":             locked,
						"last_modified":      timesheet.LastModified,
						"type":               timesheet.Type,
						"start":              timesheet.Start,
						"end":                timesheet.End,
						"date":               timesheet.Date,
						"duration":           duration.Seconds(),
						"duration_minutes":   duration.Minutes(),
						"duration_hours":     duration.Hours(),
						"on_the_clock":       timesheet.OnTheClock,
						"__syncAction":       "SET",
					}
					if customfields, ok := timesheet.CustomFields.(map[string]any); ok {
						for key, value := range customfields {
							item[key] = value
						}
					}
					items = append(items, item)
				}

				deletedTimesheets, moreDeletedTimesheets, requestError = utils.APIRequest[timesheetDeletedRequest, any, TimesheetData](&timesheetDeletedReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/timesheets_deleted", params.Account.AccessToken, "timesheets_deleted")
				if requestError.Err != nil {
					if requestError.TryLater {
						utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
						return
					}
					utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with timesheet data delta request: %w", requestError.Err))
					return
				}

				for _, deletedTimesheet := range deletedTimesheets {
					deletedItem := map[string]any{
						"id":           deletedTimesheet.ID.String(),
						"__syncAction": "REMOVE",
					}
					deletedItems = append(deletedItems, deletedItem)
				}

			} else if params.Pagination.Type == int(GetDeleted) {
				deletedTimesheets, moreDeletedTimesheets, requestError = utils.APIRequest[timesheetDeletedRequest, any, TimesheetData](&timesheetDeletedReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/timesheets_deleted", params.Account.AccessToken, "timesheets_deleted")
				if requestError.Err != nil {
					if requestError.TryLater {
						utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
						return
					}
					utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with timesheet data delta request: %w", requestError.Err))
					return
				}

				for _, deletedTimesheet := range deletedTimesheets {
					deletedItem := map[string]any{
						"id":           deletedTimesheet.ID.String(),
						"__syncAction": "REMOVE",
					}
					deletedItems = append(deletedItems, deletedItem)
				}
			} else if params.Pagination.Type == int(GetItems) {
				timesheets, moreTimesheets, requestError = utils.APIRequest[timesheetDeltaRequest, any, TimesheetData](&timesheetDeltaReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/timesheets", params.Account.AccessToken, "timesheets")
				if requestError.Err != nil {
					if requestError.TryLater {
						utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
						return
					}
					utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with timesheet data delta request: %w", requestError.Err))
					return
				}

				for _, timesheet := range timesheets {
					locked := false
					lockedNum, err := timesheet.Locked.Int64()
					if err != nil {
						utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error decoding locked value into int: %w", err))
						return
					}
					if lockedNum > 0 {
						locked = true
					}
					duration := time.Duration(timesheet.DurationSeconds) * time.Second
					item := map[string]any{
						"id":                 timesheet.ID.String(),
						"name":               timesheet.Notes,
						"user_id":            timesheet.UserID.String(),
						"created_by_user_id": timesheet.CreatedByUserID.String(),
						"jobcode_id":         timesheet.JobcodeID.String(),
						"locked":             locked,
						"last_modified":      timesheet.LastModified,
						"type":               timesheet.Type,
						"start":              timesheet.Start,
						"end":                timesheet.End,
						"date":               timesheet.Date,
						"duration":           duration.Seconds(),
						"duration_minutes":   duration.Minutes(),
						"duration_hours":     duration.Hours(),
						"on_the_clock":       timesheet.OnTheClock,
						"__syncAction":       "SET",
					}
					if customfields, ok := timesheet.CustomFields.(map[string]any); ok {
						for key, value := range customfields {
							item[key] = value
						}
					}
					items = append(items, item)
				}
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid type code in pagination parameter: %d", params.Pagination.Type))
				return
			}

			allItems := append(items, deletedItems...)

			more := moreTimesheets && moreDeletedTimesheets
			typesToSync := GetNothing
			if moreTimesheets && moreDeletedTimesheets {
				typesToSync = GetBoth
			} else if moreDeletedTimesheets {
				typesToSync = GetDeleted
			} else if moreTimesheets {
				typesToSync = GetItems
			}

			resp := response[map[string]any]{
				Items: allItems,
				Pagination: pagination{
					HasNext: more,
					NextPageConfig: nextPageConfig{
						Page: page + 1,
						Type: int(typesToSync),
					},
				},
				SynchronizationType: string(DeltaSync),
			}

			utils.RespondWithJSON(w, http.StatusOK, resp)
			return
		}
	case "jobcode":
		type jobcodeRequest struct {
			Active           string `url:"active"`
			Page             int    `url:"page"`
			SupplementalData string `url:"supplemental_data"`
			ModifiedSince    string `url:"modified_since,omitempty"`
			Type             string `url:"type"`
		}

		jobcodeReq := jobcodeRequest{
			Active:           "both",
			Page:             page,
			SupplementalData: "no",
			ModifiedSince:    lastSyncronized,
			Type:             "all",
		}

		jobcodes, more, requestError := utils.APIRequest[jobcodeRequest, any, JobcodeData](&jobcodeReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/jobcodes", params.Account.AccessToken, "jobcodes")
		if requestError.Err != nil {
			if requestError.TryLater {
				utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
				return
			}
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with jobcode data request: %w", requestError.Err))
			return
		}

		var items = []map[string]any{}

		for _, jobcode := range jobcodes {
			item := map[string]any{
				"id":                        jobcode.ID.String(),
				"name":                      jobcode.Name,
				"parent_id":                 jobcode.ParentID.String(),
				"type":                      jobcode.Type,
				"billable":                  jobcode.Billable,
				"active":                    jobcode.Active,
				"connected_with_quickbooks": jobcode.ConnectWithQuickbooks,
			}
			if customfields, ok := jobcode.CustomFields.(map[string]any); ok {
				for key, value := range customfields {
					item[key] = value
				}
			}
			if sync != FullSync {
				item["__syncAction"] = "SET"
			}
			items = append(items, item)
		}

		resp := response[map[string]any]{
			Items: items,
			Pagination: pagination{
				HasNext: more,
				NextPageConfig: nextPageConfig{
					Page: page + 1,
				},
			},
			SynchronizationType: string(sync),
		}

		utils.RespondWithJSON(w, http.StatusOK, resp)
		return
	default:
		if slices.Contains(params.Types, params.RequestedType) {
			type customfieldRequest struct {
				Active           string `url:"active,omitempty"`
				Page             int    `url:"page"`
				SupplementalData string `url:"supplemental_data"`
				ModifiedSince    string `url:"modified_since,omitempty"`
				CustomfieldID    string `url:"customfield_id"`
			}

			customfieldReq := customfieldRequest{
				Active:           "both",
				Page:             page,
				SupplementalData: "no",
				ModifiedSince:    lastSyncronized,
				CustomfieldID:    params.RequestedType,
			}

			customfieldItems, more, requestError := utils.APIRequest[customfieldRequest, any, CustomfieldData](&customfieldReq, nil, http.MethodGet, "https://rest.tsheets.com/api/v1/customfielditems", params.Account.AccessToken, "customfielditems")
			if requestError.Err != nil {
				if requestError.TryLater {
					utils.RespondWithTryLater(w, http.StatusTooManyRequests, fmt.Errorf("rate limit reached: %w", requestError.Err))
					return
				}
				utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("error with customField data request: %w", requestError.Err))
				return
			}

			if sync != FullSync {
				for i := range customfieldItems {
					customfieldItems[i].SyncAction = "SET"
				}
			}

			resp := response[CustomfieldData]{
				Items: customfieldItems,
				Pagination: pagination{
					HasNext: more,
					NextPageConfig: nextPageConfig{
						Page: page + 1,
					},
				},
				SynchronizationType: string(sync),
			}

			utils.RespondWithJSON(w, http.StatusOK, resp)
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid requested datatype"))
		return
	}
}
