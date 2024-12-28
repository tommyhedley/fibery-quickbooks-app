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
	"golang.org/x/sync/singleflight"
)

func main() {
	port := os.Getenv("PORT")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	loggerStyle := os.Getenv("LOGGER_STYLE")

	c := cache.New(12*time.Hour, 12*time.Hour)
	var group singleflight.Group

	SlogConfig := newSlogConfig(loggerLevel, loggerStyle)
	logger := SlogConfig.Create()
	slog.SetDefault(logger)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewServer(c, &group),
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
