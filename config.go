package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var version = "dev"

var loggerLevels = map[string]slog.Level{
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
}

type ProgramConfig struct {
	Mode                       string
	Port                       string
	LoggerLevel                slog.Level
	LoggerStyle                string
	RefreshSecBeforeExpriation int
	QuickBooks                 struct {
		PageSize                    int
		MinorVersion                string
		Scope                       string
		WebhookToken                string
		DiscoveryEndpointSandbox    string
		DiscoveryEndpointProduction string
		EndpointSandbox             string
		EndpointProduction          string
		OauthClientIdSandbox        string
		OauthClientSecretSandbox    string
		OauthClientIdProduction     string
		OauthClientSecretProduction string
	}
}

func NewProgramConfig(pageSize, refreshSecBeforeExpiration int) ProgramConfig {
	config := ProgramConfig{}
	config.Load(pageSize, refreshSecBeforeExpiration)
	return config
}

func (c *ProgramConfig) Load(pageSize, refreshSecBeforeExpiration int) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to find .env file")
	}

	mode := os.Getenv("MODE")
	if mode != "production" && mode != "sandbox" {
		log.Fatalf("invalid environment variable: MODE = %s\nvalid options: 'sandbox' or 'production'\n", mode)
	}
	c.Mode = mode

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("invalid environment variable: PORT is required")
	}
	c.Port = port

	loggerLevel := os.Getenv("LOGGER_LEVEL")
	slogLevel := slog.LevelInfo
	if level, exists := loggerLevels[loggerLevel]; exists {
		slogLevel = level
	} else {
		log.Printf("invalid environment variable: LOGGER_LEVEL = %s, 'info' used by default\nvalid options: 'info', 'debug', 'error', 'warn'\n", slogLevel)
	}
	c.LoggerLevel = slogLevel

	loggerStyle := os.Getenv("LOGGER_STYLE")
	c.LoggerStyle = loggerStyle

	if refreshSecBeforeExpiration == 0 {
		log.Fatalln("invalid config variable: refreshSecBeforeExpiration must be more than 0")
	}
	c.RefreshSecBeforeExpriation = refreshSecBeforeExpiration

	if pageSize < 1 || pageSize > 1000 {
		log.Fatalln("missing or invalid config variable: pageSize must be more than 0 and less than 1000")
	}
	c.QuickBooks.PageSize = pageSize

	minorVersion := os.Getenv("MINOR_VERSION")
	if minorVersion == "" {
		log.Fatalln("missing environment variable: SCOPE is required")
	}
	c.QuickBooks.MinorVersion = minorVersion

	scope := os.Getenv("SCOPE")
	if scope == "" {
		log.Fatalln("missing environment variable: SCOPE is required")
	}
	c.QuickBooks.Scope = scope

	webhookToken := os.Getenv("WEBHOOK_TOKEN")
	if webhookToken == "" {
		log.Fatalln("missing environment variable: WEBHOOK_TOKEN is required")
	}
	c.QuickBooks.WebhookToken = webhookToken

	discoveryEndpointSandbox := os.Getenv("DISCOVERY_ENDPOINT_SANDBOX")
	if mode == "sandbox" && discoveryEndpointSandbox == "" {
		log.Fatalln("missing environment variable: DISCOVERY_ENDPOINT_SANDBOX is required")
	}
	c.QuickBooks.DiscoveryEndpointSandbox = discoveryEndpointSandbox

	discoveryEndpointProduction := os.Getenv("DISCOVERY_ENDPOINT_PRODUCTION")
	if mode == "production" && discoveryEndpointProduction == "" {
		log.Fatalln("missing environment variable: DISCOVERY_ENDPOINT_PRODUCTION is required")
	}
	c.QuickBooks.DiscoveryEndpointProduction = discoveryEndpointProduction

	endpointSandbox := os.Getenv("ENDPOINT_SANDBOX")
	if mode == "sandbox" && endpointSandbox == "" {
		log.Fatalln("missing environment variable: ENDPOINT_SANDBOX is required")
	}
	c.QuickBooks.EndpointSandbox = endpointSandbox

	endpointProduction := os.Getenv("ENDPOINT_PRODUCTION")
	if mode == "production" && endpointProduction == "" {
		log.Fatalln("missing environment variable: ENDPOINT_PRODUCTION is required")
	}
	c.QuickBooks.EndpointProduction = endpointProduction

	oauthClientIdSandbox := os.Getenv("OAUTH_CLIENT_ID_SANDBOX")
	if mode == "sandbox" && oauthClientIdSandbox == "" {
		log.Fatalln("missing environment variable: OAUTH_CLIENT_ID_SANDBOX is required")
	}
	c.QuickBooks.OauthClientIdSandbox = oauthClientIdSandbox

	oauthClientSecretSandbox := os.Getenv("OAUTH_CLIENT_SECRET_SANDBOX")
	if mode == "sandbox" && oauthClientSecretSandbox == "" {
		log.Fatalln("missing environment variable: OAUTH_CLIENT_SECRET_SANDBOX is required")
	}
	c.QuickBooks.OauthClientSecretSandbox = oauthClientSecretSandbox

	oauthClientIdProduction := os.Getenv("OAUTH_CLIENT_ID_PRODUCTION")
	if mode == "production" && oauthClientIdProduction == "" {
		log.Fatalln("missing environment variable: OAUTH_CLIENT_ID_PRODUCTION is required")
	}
	c.QuickBooks.OauthClientIdProduction = oauthClientIdProduction

	oauthClientSecretProduction := os.Getenv("OAUTH_CLIENT_SECRET_PRODUCTION")
	if mode == "production" && oauthClientSecretProduction == "" {
		log.Fatalln("missing environment variable: OAUTH_CLIENT_SECRET_PRODUTION is required")
	}
	c.QuickBooks.OauthClientSecretProduction = oauthClientSecretProduction
}

func (c *ProgramConfig) BuildLogger() *slog.Logger {
	var handler slog.Handler
	switch c.LoggerStyle {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: c.LoggerLevel})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.LoggerLevel})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.LoggerLevel, AddSource: true})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: c.LoggerLevel})
		log.Printf("invalid environment variable: LOGGER_STYLE = %s, 'text' used by default\nvalid options: 'text', 'json', 'dev'\n", c.LoggerStyle)
	}

	return slog.New(handler)
}

func (c *ProgramConfig) NewClientRequest(discovery *quickbooks.DiscoveryAPI, client *http.Client) quickbooks.ClientRequest {
	clientRequest := quickbooks.ClientRequest{
		Client:       client,
		DiscoveryAPI: discovery,
	}
	switch c.Mode {
	case "production":
		clientRequest.ClientId = c.QuickBooks.OauthClientIdProduction
		clientRequest.ClientSecret = c.QuickBooks.OauthClientSecretProduction
		clientRequest.Endpoint = c.QuickBooks.EndpointProduction
	case "sandbox":
		clientRequest.ClientId = c.QuickBooks.OauthClientIdSandbox
		clientRequest.ClientSecret = c.QuickBooks.OauthClientSecretSandbox
		clientRequest.Endpoint = c.QuickBooks.EndpointSandbox
	}

	return clientRequest
}

func (c *ProgramConfig) DiscoverURL() string {
	switch c.Mode {
	case "production":
		return c.QuickBooks.DiscoveryEndpointProduction
	case "sandbox":
		return c.QuickBooks.DiscoveryEndpointSandbox
	default:
		return ""
	}
}

func AppConfig(version string) fibery.AppConfig {
	return fibery.AppConfig{
		Id:          "qbo",
		Name:        "QuickBooks Online",
		Website:     "https://quickbooks.intuit.com",
		Version:     version,
		Description: "Integrate QuickBooks Online data with Fibery",
		Authentication: []fibery.Authentication{
			{
				Id:          "oauth2",
				Name:        "OAuth v2 Authentication",
				Description: "OAuth v2-based authentication and authorization for access to QuickBooks Online",
				Fields: []fibery.AuthField{
					{
						Id:          "callback_uri",
						Title:       "callback_uri",
						Description: "OAuth post-auth redirect URI",
						Type:        "oauth",
					},
				},
			},
		},
		Sources: []string{},
		ResponsibleFor: fibery.ResponsibleFor{
			DataSynchronization: true,
		},
	}
}

func SyncConfig(types TypeRegistry) fibery.SyncConfig {
	syncConfig := fibery.SyncConfig{
		Types:   make([]fibery.SyncConfigTypes, len(types)),
		Filters: []fibery.SyncFilter{},
		Webhooks: fibery.SyncConfigWebhook{
			Enabled: true,
			Type:    "ui",
		},
	}

	i := 0
	for _, typ := range types {
		syncConfig.Types[i] = fibery.SyncConfigTypes{
			Id:   typ.Id(),
			Name: typ.Name(),
		}
		i++
	}

	return syncConfig
}
