package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/advice"
	"requiems-api/internal/config"
	"requiems-api/internal/db"
	"requiems-api/internal/quotes"
	"requiems-api/internal/words"
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

		r.Get("/healthz", Healthz)

	adviceSvc := advice.NewService(pool)
	advice.RegisterRoutes(r, adviceSvc)

	quotesSvc := quotes.NewService(pool)
	quotes.RegisterRoutes(r, quotesSvc)

	wordsSvc := words.NewService(pool)	
	words.RegisterRoutes(r, wordsSvc)

	return &App{
		cfg: cfg,
		h:   r,
	}, nil
}

func (a *App) Handler() http.Handler {
	return a.h
}

// migrateWithRetry wraps db.Migrate with a small retry loop so that
// we don't fail just because Postgres is still starting up.
func migrateWithRetry(dsn, dir string) error {
	var lastErr error

	for i := 0; i < 10; i++ {
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

