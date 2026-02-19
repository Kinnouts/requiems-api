package db

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dsn, dir string) error {
	if !filepath.IsAbs(dir) {
		abs, err := filepath.Abs(dir)

		if err != nil {
			return fmt.Errorf("resolve migrations path: %w", err)
		}

		dir = abs
	}

	sourceURL := "file://" + dir

	m, err := migrate.New(sourceURL, dsn)

	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func MigrateWithRetry(dsn, dir string) error {
	var lastErr error

	for range 10 {
		if err := Migrate(dsn, dir); err != nil {
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
