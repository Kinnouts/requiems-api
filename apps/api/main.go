package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	appdb "requiems-api/internal/db"
)

type adviceResponse struct {
	ID     int    `json:"id"`
	Advice string `json:"advice"`
}

func main() {
	ctx := context.Background()

	pool, err := appdb.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	if err := appdb.Migrate(ctx, pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/v1/advice", adviceHandler(pool))

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

func adviceHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		row := pool.QueryRow(r.Context(), `
SELECT id, text
FROM advice
ORDER BY random()
LIMIT 1;
`)
		var id int
		var text string
		if err := row.Scan(&id, &text); err != nil {
			log.Printf("query advice failed: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "no advice available"})
			return
		}

		resp := adviceResponse{
			ID:     id,
			Advice: text,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}



