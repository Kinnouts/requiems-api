package db

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate runs all up migrations using golang-migrate.
// dsn should match the DATABASE_URL, and dir is a filesystem path like "infra/migrations".
func Migrate(dsn, dir string) error {
	// Ensure we always pass an absolute path to the file:// source.
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
