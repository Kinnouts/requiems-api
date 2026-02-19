package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/email"
	"requiems-api/internal/misc"
	"requiems-api/internal/entertainment"
	"requiems-api/internal/platform/config"
	"requiems-api/internal/platform/db"
	"requiems-api/internal/platform/middleware"
	"requiems-api/internal/text"
)

type App struct {
	cfg config.Config
	h   http.Handler
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	if err := migrateWithRetry(cfg.DatabaseURL, "infra/migrations"); err != nil {
		return nil, err
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)

	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	// Public routes (no auth required)
	r.Get("/healthz", Healthz)

	// Protected routes (require X-Backend-Secret header)
	r.Group(func(protected chi.Router) {
		// Apply backend secret authentication middleware
		protected.Use(middleware.BackendSecretAuth(cfg.BackendSecret))

		textRouter := chi.NewRouter()
		text.RegisterRoutes(textRouter, pool)
		protected.Mount("/v1/text", textRouter)

		emailRouter := chi.NewRouter()
		email.RegisterRoutes(emailRouter)
		protected.Mount("/v1/email", emailRouter)

		miscRouter := chi.NewRouter()
		misc.RegisterRoutes(miscRouter)
		protected.Mount("/v1/misc", miscRouter)
    
		entertainmentRouter := chi.NewRouter()
		entertainment.RegisterRoutes(entertainmentRouter)
		protected.Mount("/v1/entertainment", entertainmentRouter)
	})

	return &App{
		cfg: cfg,
		h:   r,
	}, nil
}

func (a *App) Handler() http.Handler {
	return a.h
}

func migrateWithRetry(dsn, dir string) error {
	var lastErr error

	for range 10 {
		if err := db.Migrate(dsn, dir); err != nil {
			lastErr = err

			msg := err.Error()

			if strings.Contains(msg, "the database system is starting up") ||
				strings.Contains(msg, "connection refused") {
				time.Sleep(1 * time.Second)
				continue
			}

			return err
		}

		return nil
	}

	if lastErr != nil {
		return fmt.Errorf("migrations failed after retries: %w", lastErr)
	}

	return fmt.Errorf("migrations failed after retries with unknown error")
}
