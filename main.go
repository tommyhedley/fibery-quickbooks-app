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

	"github.com/patrickmn/go-cache"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/data"
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/fibery"
	"golang.org/x/sync/singleflight"
)

func main() {
	port := os.Getenv("PORT")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	loggerStyle := os.Getenv("LOGGER_STYLE")

	c := cache.New(12*time.Hour, 12*time.Hour)
	var group singleflight.Group

	integration := Integration{
		cache: c,
		group: &group,
		appConfig: fibery.AppConfig{
			ID:          "qbo",
			Name:        "QuickBooks Online",
			Website:     "https://quickbooks.intuit.com",
			Version:     "0.1.0",
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
		types: data.Types,
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
