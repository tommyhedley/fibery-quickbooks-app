package synchronizer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommyhedley/fibery/fibery-tsheets-integration/internal/utils"
)

type SyncConfig struct {
	Types    []DB        `json:"types"`
	Filters  []Filter    `json:"filters"`
	Webhooks SyncWebhook `json:"webhooks,omitempty"`
}

type DB struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Filter struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Datalist bool   `json:"datalist,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Secured  bool   `json:"secured,omitempty"`
}

type SyncWebhook struct {
	Enabled bool `json:"enabled,omitempty"`
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := SyncConfig{
		Types: []DB{
			{
				ID:   "expense",
				Name: "Expense",
			},
			{
				ID:   "expenseItem",
				Name: "Expense Item",
			},
			{
				ID:   "invoice",
				Name: "Invoice",
			},
			{
				ID:   "invoiceItem",
				Name: "Invoice Item",
			},
			{
				ID:   "vendor",
				Name: "Vendor",
			},
			{
				ID:   "customer",
				Name: "Customer",
			},
			{
				ID:   "customerType",
				Name: "Customer Type",
			},
			{
				ID:   "account",
				Name: "Account",
			},
			{
				ID:   "item",
				Name: "Item",
			},
			{
				ID:   "class",
				Name: "Class",
			},
			{
				ID:   "taxCode",
				Name: "Tax Code",
			},
			{
				ID:   "taxExemption",
				Name: "Tax Exemption",
			},
			{
				ID:   "bill",
				Name: "Bill",
			},
			{
				ID:   "term",
				Name: "Term",
			},
			{
				ID:   "employee",
				Name: "Employee",
			},
		},
		Filters: []Filter{},
		Webhooks: SyncWebhook{
			Enabled: true,
		},
	}

	type requestBody struct {
		Types   []string `json:"types"`
		Account struct {
			AccessToken string `json:"access_token"`
		} `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to decode request parameters: %w", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, config)
}

func ValidateFiltersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Types   []string       `json:"types"`
		Filter  map[string]any `json:"filter"`
		Account struct {
			Token string `json:"token"`
		} `json:"account"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("couldn't decode parameters"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil)
}
