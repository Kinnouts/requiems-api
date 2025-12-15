package main

import (
	"context"
	"log"
	"net/http"

	"requiems-api/internal/app"
	"requiems-api/internal/config"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	a, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	addr := ":" + cfg.Port
	
	log.Printf("API server listening on %s\n", addr)

	if err := http.ListenAndServe(addr, a.Handler()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

