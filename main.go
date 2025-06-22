package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/app"
	_ "github.com/tommyhedley/fibery-quickbooks-app/pkgs/app/types"
)

func main() {
	shutdownCtx, shutdownCancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer shutdownCancel()

	params := app.Parameters{
		Version:                    "dev-v0.0.3",
		PageSize:                   1000,
		RefreshSecBeforeExpiration: 600,
		AttachableFieldId:          "attachables",
		OperationTTL:               time.Duration(15 * time.Second),
		IdCacheTTL:                 time.Duration(24 * time.Hour),
	}

	a, err := app.New(shutdownCtx, params)
	if err != nil {
		log.Fatalf("unable to create new integration: %s", err.Error())
	}

	server := &http.Server{
		Addr:    ":" + a.Port(),
		Handler: app.NewHandler(a),
		BaseContext: func(net.Listener) context.Context {
			return shutdownCtx
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-shutdownCtx.Done()
		slog.Info("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("Could not gracefully shutdown the server %+v", err))
		}
	}()

	slog.Info(fmt.Sprintf("Server starting at port %s...", a.Port()))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("Could not listen on :%s %+v", a.Port(), err))
	}

	slog.Info("Server stopped")
}
