package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tommyhedley/fibery/fibery-tsheets-integration/oauth2"
)

var SlogConfig sloggerConfig

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	loggerStyle := os.Getenv("LOGGER_STYLE")

	SlogConfig = newSlogConfig(loggerLevel, loggerStyle)
	httpLogger := SlogConfig.Create()

	err := oauth2.GetDiscovery()
	if err != nil {
		log.Fatalf("unable to get discovery info: %v", err)
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewServer(httpLogger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server %+v\n", err)
		}
		close(done)
	}()

	log.Printf("Server starting at port %s...\n", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :%s %+v\n", port, err)
	}

	<-done
	log.Println("Server stopped")
}
