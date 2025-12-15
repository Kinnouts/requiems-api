package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"requiems-api/internal/advice"
	appdb "requiems-api/internal/db"
	"requiems-api/internal/quotes"
	"requiems-api/internal/words"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"
	}

	// Run migrations before starting the app.
	if err := appdb.Migrate(dsn, "infra/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	pool, err := appdb.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Feature modules
	adviceSvc := advice.NewService(pool)
	advice.RegisterHTTP(mux, adviceSvc)

	quotesSvc := quotes.NewService(pool)
	quotes.RegisterHTTP(mux, quotesSvc)

	wordsSvc := words.NewService(pool)
	words.RegisterHTTP(mux, wordsSvc)

	addr := ":" + envOrDefault("PORT", "8080")
	log.Printf("API server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

