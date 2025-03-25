package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tommyhedley/fibery-quickbooks-app/data"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

var version = "dev"

type Integration struct {
	appConfig  fibery.AppConfig
	syncConfig fibery.SyncConfig
	client     *quickbooks.Client
	dataCache  *data.Cache
	idCache    *data.IdCache
}

func NewIntegration(appConfig fibery.AppConfig, syncConfig fibery.SyncConfig, client *quickbooks.Client) *Integration {
	dataCache := data.NewDataCache(30 * time.Second)
	idCache := data.NewIdCache(24 * time.Hour)
	integration := &Integration{
		appConfig:  appConfig,
		syncConfig: syncConfig,
		client:     client,
		dataCache:  dataCache,
		idCache:    idCache,
	}
	integration.StartCacheCleaner()
	return integration
}

func (i *Integration) Cleanup() {
	i.dataCache.CleanupExpired()
	i.idCache.CleanupExpired()
}

func (i *Integration) StartCacheCleaner() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			i.Cleanup()
		}
	}()
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	var discoveryUrl string
	switch os.Getenv("MODE") {
	case "production":
		discoveryUrl = os.Getenv("DISCOVERY_PRODUCTION_ENDPOINT")
	case "sandbox":
		discoveryUrl = os.Getenv("DISCOVERY_SANDBOX_ENDPOINT")
	default:
		slog.Error("Invalid mode", "error", os.Getenv("MODE"))
		os.Exit(1)
	}

	discoveryAPI, err := quickbooks.CallDiscoveryAPI(discoveryUrl)
	if err != nil {
		slog.Error("Error calling discovery API", "error", err)
		os.Exit(1)
	}

	clientReq, err := NewClientRequest(discoveryAPI, http.DefaultClient)
	if err != nil {
		slog.Error("Error creating client request", "error", err)
		os.Exit(1)
	}

	client, err := quickbooks.NewClient(clientReq)
	if err != nil {
		slog.Error("Error creating client", "error", err)
		os.Exit(1)
	}

	appConfig := fibery.AppConfig{
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

	syncConfig := fibery.SyncConfig{
		Types:   []fibery.SyncConfigTypes{},
		Filters: []fibery.SyncFilter{},
		Webhooks: fibery.SyncConfigWebhook{
			Enabled: true,
			Type:    "ui",
		},
	}

	for _, typ := range data.Types.All {
		syncConfig.Types = append(syncConfig.Types, fibery.SyncConfigTypes{
			Id:   (*typ).Id(),
			Name: (*typ).Name(),
		})
	}

	integration := NewIntegration(appConfig, syncConfig, client)
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	loggerStyle := os.Getenv("LOGGER_STYLE")
	SlogConfig := newSlogConfig(loggerLevel, loggerStyle)
	logger := SlogConfig.Create()
	slog.SetDefault(logger)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewIntegrationHandler(integration),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		slog.Info("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("Could not gracefully shutdown the server %+v", err))
		}
		close(done)
	}()

	slog.Info(fmt.Sprintf("Server starting at port %s...", port))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("Could not listen on :%s %+v", port, err))
	}

	<-done
	slog.Info("Server stopped")
}
