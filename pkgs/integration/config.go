package integration

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tommyhedley/quickbooks-go"
)

var loggerLevels = map[string]slog.Level{
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
}

type Config struct {
	Mode                       string
	Port                       string
	LoggerLevel                slog.Level
	LoggerStyle                string
	RefreshSecBeforeExpriation int
	AttachableFieldId          string
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

type Parameters struct {
	PageSize                   int
	RefreshSecBeforeExpiration int
	Version                    string
	AttachableFieldId          string
	OperationTTL               time.Duration
	IdCacheTTL                 time.Duration
}

func BuildConfig(params Parameters) (Config, error) {
	config := Config{}
	err := config.Load(params.PageSize, params.RefreshSecBeforeExpiration, params.AttachableFieldId)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) Load(pageSize, refreshSecBeforeExpiration int, attachableFieldId string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("unable to find .env file")
	}

	mode := os.Getenv("MODE")
	if mode != "production" && mode != "sandbox" {
		return fmt.Errorf("invalid environment variable: MODE = %s\nvalid options: 'sandbox' or 'production'\n", mode)
	}
	c.Mode = mode

	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("invalid environment variable: PORT is required")
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
		return fmt.Errorf("invalid config variable: refreshSecBeforeExpiration must be more than 0")
	}
	c.RefreshSecBeforeExpriation = refreshSecBeforeExpiration

	if attachableFieldId == "" {
		return fmt.Errorf("no specified id for an quickbooks.Attachables schema field")
	}
	c.AttachableFieldId = attachableFieldId

	if pageSize < 1 || pageSize > 1000 {
		return fmt.Errorf("missing or invalid config variable: pageSize must be more than 0 and less than 1000")
	}
	c.QuickBooks.PageSize = pageSize

	minorVersion := os.Getenv("MINOR_VERSION")
	if minorVersion == "" {
		return fmt.Errorf("missing environment variable: SCOPE is required")
	}
	c.QuickBooks.MinorVersion = minorVersion

	scope := os.Getenv("SCOPE")
	if scope == "" {
		return fmt.Errorf("missing environment variable: SCOPE is required")
	}
	c.QuickBooks.Scope = scope

	webhookToken := os.Getenv("WEBHOOK_TOKEN")
	if webhookToken == "" {
		return fmt.Errorf("missing environment variable: WEBHOOK_TOKEN is required")
	}
	c.QuickBooks.WebhookToken = webhookToken

	discoveryEndpointSandbox := os.Getenv("DISCOVERY_ENDPOINT_SANDBOX")
	if mode == "sandbox" && discoveryEndpointSandbox == "" {
		return fmt.Errorf("missing environment variable: DISCOVERY_ENDPOINT_SANDBOX is required")
	}
	c.QuickBooks.DiscoveryEndpointSandbox = discoveryEndpointSandbox

	discoveryEndpointProduction := os.Getenv("DISCOVERY_ENDPOINT_PRODUCTION")
	if mode == "production" && discoveryEndpointProduction == "" {
		return fmt.Errorf("missing environment variable: DISCOVERY_ENDPOINT_PRODUCTION is required")
	}
	c.QuickBooks.DiscoveryEndpointProduction = discoveryEndpointProduction

	endpointSandbox := os.Getenv("ENDPOINT_SANDBOX")
	if mode == "sandbox" && endpointSandbox == "" {
		return fmt.Errorf("missing environment variable: ENDPOINT_SANDBOX is required")
	}
	c.QuickBooks.EndpointSandbox = endpointSandbox

	endpointProduction := os.Getenv("ENDPOINT_PRODUCTION")
	if mode == "production" && endpointProduction == "" {
		return fmt.Errorf("missing environment variable: ENDPOINT_PRODUCTION is required")
	}
	c.QuickBooks.EndpointProduction = endpointProduction

	oauthClientIdSandbox := os.Getenv("OAUTH_CLIENT_ID_SANDBOX")
	if mode == "sandbox" && oauthClientIdSandbox == "" {
		return fmt.Errorf("missing environment variable: OAUTH_CLIENT_ID_SANDBOX is required")
	}
	c.QuickBooks.OauthClientIdSandbox = oauthClientIdSandbox

	oauthClientSecretSandbox := os.Getenv("OAUTH_CLIENT_SECRET_SANDBOX")
	if mode == "sandbox" && oauthClientSecretSandbox == "" {
		return fmt.Errorf("missing environment variable: OAUTH_CLIENT_SECRET_SANDBOX is required")
	}
	c.QuickBooks.OauthClientSecretSandbox = oauthClientSecretSandbox

	oauthClientIdProduction := os.Getenv("OAUTH_CLIENT_ID_PRODUCTION")
	if mode == "production" && oauthClientIdProduction == "" {
		return fmt.Errorf("missing environment variable: OAUTH_CLIENT_ID_PRODUCTION is required")
	}
	c.QuickBooks.OauthClientIdProduction = oauthClientIdProduction

	oauthClientSecretProduction := os.Getenv("OAUTH_CLIENT_SECRET_PRODUCTION")
	if mode == "production" && oauthClientSecretProduction == "" {
		return fmt.Errorf("missing environment variable: OAUTH_CLIENT_SECRET_PRODUTION is required")
	}
	c.QuickBooks.OauthClientSecretProduction = oauthClientSecretProduction

	return nil
}

func (c *Config) BuildLogger() *slog.Logger {
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

func (c *Config) NewClientRequest(discovery *quickbooks.DiscoveryAPI, client *http.Client) quickbooks.ClientRequest {
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

func (c *Config) DiscoverURL() string {
	switch c.Mode {
	case "production":
		return c.QuickBooks.DiscoveryEndpointProduction
	case "sandbox":
		return c.QuickBooks.DiscoveryEndpointSandbox
	default:
		return ""
	}
}
