package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tommyhedley/quickbooks-go"
)

func main() {
	config := NewProgramConfig(1000, 600)
	types := NewTypeRegistry()

	discoveryAPI, err := quickbooks.CallDiscoveryAPI(config.DiscoverURL())
	if err != nil {
		log.Fatalf("error calling discovery API: %s\n", err.Error())
	}

	clientReq := config.NewClientRequest(discoveryAPI, http.DefaultClient)

	client, err := quickbooks.NewClient(clientReq)
	if err != nil {
		log.Fatalf("error creating quickbooks client: %s\n", err.Error())
	}

	operationTTL := time.Duration(30 * time.Second)
	idCacheTTL := time.Duration(24 * time.Hour)

	integration := NewIntegration(AppConfig(version), SyncConfig(types), config, types, client, operationTTL, idCacheTTL)
	slog.SetDefault(config.BuildLogger())

	server := &http.Server{
		Addr:         ":" + config.Port,
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

	slog.Info(fmt.Sprintf("Server starting at port %s...", config.Port))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("Could not listen on :%s %+v", config.Port, err))
	}

	<-done
	slog.Info("Server stopped")
}
