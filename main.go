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
	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery-quickbooks-app/data"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
	"golang.org/x/sync/singleflight"
)

var version = "dev"

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	loggerStyle := os.Getenv("LOGGER_STYLE")

	c := cache.New(12*time.Hour, 12*time.Hour)
	var group singleflight.Group

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

	integration := Integration{
		appConfig: fibery.AppConfig{
			ID:          "qbo",
			Name:        "QuickBooks Online",
			Website:     "https://quickbooks.intuit.com",
			Version:     version,
			Description: "Integrate QuickBooks Online data with Fibery",
			Authentication: []fibery.Authentication{
				{
					ID:          "oauth2",
					Name:        "OAuth v2 Authentication",
					Description: "OAuth v2-based authentication and authorization for access to QuickBooks Online",
					Fields: []fibery.AuthField{
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
			ResponsibleFor: fibery.ResponsibleFor{
				DataSynchronization: true,
			},
		},
		syncConfig: fibery.SyncConfig{
			Types:   []fibery.SyncConfigTypes{},
			Filters: []fibery.SyncFilter{},
			Webhooks: fibery.SyncConfigWebhook{
				Enabled: true,
				Type:    "ui",
			},
		},
		types:        data.Types,
		cache:        c,
		group:        &group,
		discoveryAPI: discoveryAPI,
		client:       client,
	}

	for _, datatype := range data.Types {
		integration.syncConfig.Types = append(integration.syncConfig.Types, fibery.SyncConfigTypes{
			ID:   (*datatype).GetId(),
			Name: (*datatype).GetName(),
		})
	}

	SlogConfig := newSlogConfig(loggerLevel, loggerStyle)
	logger := SlogConfig.Create()
	slog.SetDefault(logger)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewServer(&integration),
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
