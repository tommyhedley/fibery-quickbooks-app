package main

import (
	"net/http"
)

type AuthField struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional,omitempty"`
	Value       string `json:"value,omitempty"`
}

type Authentication struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Fields      []AuthField `json:"fields,omitempty"`
}

type ResponsibleFor struct {
	DataSynchronization bool `json:"dataSynchronization"`
	Automations         bool `json:"automations,omitempty"`
}

type ActionArg struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Type         string `json:"type"`
	TextTemplate bool   `json:"textTemplateSupported,omitempty"`
}

type Action struct {
	ActionID    string      `json:"action"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Args        []ActionArg `json:"args"`
}

type AppConfig struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Website        string           `json:"website,omitempty"`
	Version        string           `json:"version"`
	Description    string           `json:"description"`
	Authentication []Authentication `json:"authentication"`
	Sources        []string         `json:"sources"`
	ResponsibleFor ResponsibleFor   `json:"responsibleFor"`
	Actions        []Action         `json:"actions"`
}

func AppConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := AppConfig{
		ID:          "qbo",
		Name:        "QuickBooks Online",
		Website:     "https://quickbooks.intuit.com",
		Version:     "0.1.0",
		Description: "Integrate QuickBooks Online data with Fibery",
		Authentication: []Authentication{
			{
				ID:          "oauth2",
				Name:        "OAuth v2 Authentication",
				Description: "OAuth v2-based authentication and authorization for access to QuickBooks Online",
				Fields: []AuthField{
					{
						ID:          "callback_uri",
						Title:       "callback_uri",
						Description: "OAuth post-auth redirect URI",
						Type:        "oauth",
					},
				},
			},
		},
		Sources: []string{},
		ResponsibleFor: ResponsibleFor{
			DataSynchronization: true,
		},
	}
	RespondWithJSON(w, http.StatusOK, config)
}
