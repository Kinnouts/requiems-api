package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sentry "github.com/getsentry/sentry-go"

	"requiems-api/app"
	"requiems-api/platform/config"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	if cfg.Environment != "development" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			Environment:      cfg.Environment,
			TracesSampleRate: 0.01,
		}); err != nil {
			log.Printf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	appInstance, err := app.New(ctx, cfg)

	if err != nil {
		log.Printf("failed to initialise app: %v", err)
		sentry.Flush(2 * time.Second)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Println("API server listening on", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           appInstance.Handler(),
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
