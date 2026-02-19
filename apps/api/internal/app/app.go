package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/email"
	"requiems-api/internal/entertainment"
	"requiems-api/internal/misc"
	"requiems-api/internal/text"

	"requiems-api/internal/platform/config"
	"requiems-api/internal/platform/db"
	"requiems-api/internal/platform/middleware"
	"requiems-api/internal/platform/reqredis"
)

type App struct {
	cfg     config.Config
	handler http.Handler
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	if err := db.MigrateWithRetry(cfg.DatabaseURL, "infra/migrations"); err != nil {
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
			textRouter := chi.NewRouter()
			text.RegisterRoutes(textRouter, pool)
			v1.Mount("/text", textRouter)

			emailRouter := chi.NewRouter()
			email.RegisterRoutes(emailRouter)
			v1.Mount("/email", emailRouter)

			entertainmentRouter := chi.NewRouter()
			entertainment.RegisterRoutes(entertainmentRouter)
			v1.Mount("/entertainment", entertainmentRouter)

			miscRouter := chi.NewRouter()
			misc.RegisterRoutes(ctx, miscRouter, pool, rdb)
			v1.Mount("/misc", miscRouter)
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
