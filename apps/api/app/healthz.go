package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"requiems-api/platform/httpx"
)

type healthzResponse struct {
	Status string `json:"status"`
}

func (healthzResponse) IsData() {}

func Healthz(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			httpx.Error(w, http.StatusServiceUnavailable, "db_unavailable", "database ping failed")
			return
		}

		httpx.JSON(w, http.StatusOK, healthzResponse{Status: "ok"})
	}
}
