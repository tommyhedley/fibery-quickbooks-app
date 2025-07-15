package app

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var loggerLevels = map[string]slog.Level{
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
}

func getRequiredEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("missing environment variable: %s", key)
	}
	return v, nil
}

func parseDurationEnv(key string) (time.Duration, error) {
	raw, err := getRequiredEnv(key)
	if err != nil {
		return 0, err
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}

func parseIntEnv(key string) (int, error) {
	raw, err := getRequiredEnv(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return n, nil
}

type Config struct {
	fiberyApp          fibery.AppConfig
	fiberySync         fibery.SyncConfig
	Version            string
	Mode               string
	Port               string
	LoggerLevel        slog.Level
	LoggerStyle        string
	TokenRefreshWindow time.Duration
	PageSize           int
	AttachableFieldId  string
	OperationTTL       time.Duration
	IdCacheTTL         time.Duration
	QuickBooks         struct {
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

func NewConfig(version string) (Config, error) {
	config := Config{Version: version}
	err := config.Load()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) Load() error {
	filtered := []string{os.Args[0]}
	for _, a := range os.Args[1:] {
		if a == "--dotenv" || a == "-dotenv" {
			if err := godotenv.Load(); err != nil {
				return fmt.Errorf("unable to load .env: %w", err)
			}
			continue
		}
		filtered = append(filtered, a)
	}
	os.Args = filtered

	flag.StringVar(&c.Mode, "mode", os.Getenv("MODE"), "run the server in 'sandbox' or 'production'")
	flag.StringVar(&c.Port, "port", os.Getenv("PORT"), "http listen port")

	var logLevelStr string
	flag.StringVar(&logLevelStr, "log_level", os.Getenv("LOGGER_LEVEL"), "slog logger level: 'info','debug','error','warn'")
	flag.StringVar(&c.LoggerStyle, "log_style", os.Getenv("LOGGER_STYLE"), "slog logger style: 'text','json','dev'")

	flag.DurationVar(&c.TokenRefreshWindow, "token_refresh", 0, "duration before token expiration to refresh token")
	flag.DurationVar(&c.OperationTTL, "op_ttl", 0, "operation time to live")
	flag.DurationVar(&c.IdCacheTTL, "cache_ttl", 0, "cache time to live")

	flag.IntVar(&c.QuickBooks.PageSize, "page_size", 0, "quickbooks query page size → max 1000")
	flag.StringVar(&c.AttachableFieldId, "attachable_field", os.Getenv("ATTACHABLE_FIELD_ID"), "attachables field id")

	flag.Parse()

	if c.Mode != "sandbox" && c.Mode != "production" {
		return fmt.Errorf("MODE must be 'sandbox' or 'production'; got %q", c.Mode)
	}
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	if lvl, ok := loggerLevels[logLevelStr]; ok {
		c.LoggerLevel = lvl
	} else {
		log.Printf("invalid LOGGER_LEVEL=%q, defaulting to 'info'\n", logLevelStr)
		c.LoggerLevel = slog.LevelInfo
	}
	switch c.LoggerStyle {
	case "text", "json", "dev":
	default:
		log.Printf("invalid LOGGER_STYLE=%q, defaulting to 'text'\n", c.LoggerStyle)
		c.LoggerStyle = "text"
	}

	if c.TokenRefreshWindow == 0 {
		d, err := parseDurationEnv("TOKEN_REFRESH_WINDOW")
		if err != nil {
			return err
		}
		c.TokenRefreshWindow = d
	}
	if c.OperationTTL == 0 {
		d, err := parseDurationEnv("OPERATION_TTL")
		if err != nil {
			return err
		}
		c.OperationTTL = d
	}
	if c.IdCacheTTL == 0 {
		d, err := parseDurationEnv("CACHE_TTL")
		if err != nil {
			return err
		}
		c.IdCacheTTL = d
	}

	if c.QuickBooks.PageSize == 0 {
		n, err := parseIntEnv("PAGE_SIZE")
		if err != nil {
			return err
		}
		c.QuickBooks.PageSize = n
	}
	if c.QuickBooks.PageSize < 1 || c.QuickBooks.PageSize > 1000 {
		return fmt.Errorf("page_size must be between 1 and 1000")
	}

	if c.AttachableFieldId == "" {
		return fmt.Errorf("ATTACHABLE_FIELD_ID is required")
	}

	var err error
	c.QuickBooks.MinorVersion, err = getRequiredEnv("MINOR_VERSION")
	if err != nil {
		return err
	}
	c.QuickBooks.Scope, err = getRequiredEnv("SCOPE")
	if err != nil {
		return err
	}
	c.QuickBooks.WebhookToken, err = getRequiredEnv("WEBHOOK_TOKEN")
	if err != nil {
		return err
	}

	c.QuickBooks.DiscoveryEndpointSandbox = os.Getenv("DISCOVERY_ENDPOINT_SANDBOX")
	c.QuickBooks.DiscoveryEndpointProduction = os.Getenv("DISCOVERY_ENDPOINT_PRODUCTION")
	c.QuickBooks.EndpointSandbox = os.Getenv("ENDPOINT_SANDBOX")
	c.QuickBooks.EndpointProduction = os.Getenv("ENDPOINT_PRODUCTION")
	c.QuickBooks.OauthClientIdSandbox = os.Getenv("OAUTH_CLIENT_ID_SANDBOX")
	c.QuickBooks.OauthClientSecretSandbox = os.Getenv("OAUTH_CLIENT_SECRET_SANDBOX")
	c.QuickBooks.OauthClientIdProduction = os.Getenv("OAUTH_CLIENT_ID_PRODUCTION")
	c.QuickBooks.OauthClientSecretProduction = os.Getenv("OAUTH_CLIENT_SECRET_PRODUCTION")

	if c.Mode == "sandbox" {
		if c.QuickBooks.DiscoveryEndpointSandbox == "" ||
			c.QuickBooks.EndpointSandbox == "" ||
			c.QuickBooks.OauthClientIdSandbox == "" ||
			c.QuickBooks.OauthClientSecretSandbox == "" {
			return fmt.Errorf("all sandbox QuickBooks settings are required")
		}
	} else {
		if c.QuickBooks.DiscoveryEndpointProduction == "" ||
			c.QuickBooks.EndpointProduction == "" ||
			c.QuickBooks.OauthClientIdProduction == "" ||
			c.QuickBooks.OauthClientSecretProduction == "" {
			return fmt.Errorf("all production QuickBooks settings are required")
		}
	}

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
