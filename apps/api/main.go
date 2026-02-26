package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"requiems-api/app"
	"requiems-api/platform/config"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	appInstance, err := app.New(ctx, cfg)

	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
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
