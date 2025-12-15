package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Sensible default for local dev if not provided.
		dsn = "postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// Migrate applies minimal bootstrapping migrations needed for the MVP.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	// Advice table for the MVP endpoint.
	_, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS advice (
  id   SERIAL PRIMARY KEY,
  text TEXT NOT NULL
);
`)
	if err != nil {
		return fmt.Errorf("create advice table: %w", err)
	}

	// Seed initial advice rows if table is empty.
	_, err = pool.Exec(ctx, `
INSERT INTO advice (text)
SELECT v FROM (VALUES
  ('Ship small, ship often.'),
  ('Talk to users before you over-engineer.'),
  ('Good logs today save you hours tomorrow.'),
  ('Automate anything you do more than twice.'),
  ('Optimize for readability over cleverness.'),
  ('Start with the simplest data model that could work.')
) AS seed(v)
WHERE NOT EXISTS (SELECT 1 FROM advice);
`)
	if err != nil {
		return fmt.Errorf("seed advice: %w", err)
	}

	return nil
}


