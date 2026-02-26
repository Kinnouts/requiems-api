package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/config"
	"requiems-api/platform/db"
	"requiems-api/platform/middleware"
	"requiems-api/platform/reqredis"
)

type App struct {
	cfg     config.Config
	handler http.Handler
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	if err := db.MigrateWithRetry(cfg.DatabaseURL, "migrations"); err != nil {
		return nil, err
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	rdb, err := reqredis.Connect(ctx, cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	router.Get("/healthz", Healthz)

	router.Group(func(protected chi.Router) {
		protected.Use(middleware.BackendSecretAuth(cfg.BackendSecret))

		protected.Route("/v1", func(v1 chi.Router) {
			registerV1Routes(ctx, v1, pool, rdb)
		})
	})

	return &App{
		cfg:     cfg,
		handler: router,
	}, nil
}

func (a *App) Handler() http.Handler {
	return a.handler
}
