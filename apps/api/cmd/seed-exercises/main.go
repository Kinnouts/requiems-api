// seed-exercises reads exercise JSON files from a local directory and upserts
// them into the exercises PostgreSQL table.
//
// The JSON files are not stored in the repository; rsync them to the server
// before running this command in production.
//
// gifUrl and the raw exerciseId are never stored or exposed; exerciseId is kept
// only as an internal upsert key (external_id column).
//
// Usage:
//
//	go run ./cmd/seed-exercises --data-dir /path/to/json
//	go run ./cmd/seed-exercises --data-dir /path/to/json --dry-run
//	go run ./cmd/seed-exercises --data-dir /path/to/json \
//	    --db-url "postgres://requiem:requiem@db:5432/requiem"
//	go run ./cmd/seed-exercises --data-dir /path/to/json --verbose
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"requiems-api/cmd/internal/seedutil"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		dataDir string
		dbURL   string
		dryRun  bool
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "seed-exercises",
		Short: "Seed the exercises table from local JSON files",
		Long: `Reads exercises.json from --data-dir and upserts records into the
exercises PostgreSQL table. The JSON files are not stored in git — rsync them
to the server before running in production.

Run inside the API Docker container so the db hostname resolves:

  docker exec requiem-dev-api-1 go run ./cmd/seed-exercises --data-dir /data
  docker exec requiem-dev-api-1 go run ./cmd/seed-exercises --data-dir /data --dry-run`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dataDir == "" {
				return fmt.Errorf("--data-dir is required")
			}
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			ctx := context.Background()
			start := time.Now()

			log.Printf("loading exercises from %s", dataDir)
			records, err := loadExercises(dataDir)
			if err != nil {
				return fmt.Errorf("load exercises: %w", err)
			}

			if len(records) == 0 {
				return fmt.Errorf("no valid exercise records found in %s", dataDir)
			}

			log.Printf("loaded %d exercises", len(records))

			if verbose {
				for _, r := range records {
					fmt.Printf("  [%s] %s  muscles=%v  parts=%v\n",
						r.ExternalID, r.Name, r.TargetMuscles, r.BodyParts)
				}
			}

			if dryRun {
				log.Println("dry-run: no data written to database")
				return nil
			}

			conn, err := seedutil.OpenDB(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			inserted, updated, err := upsertRecords(ctx, conn, records)
			conn.Close(ctx)
			if err != nil {
				return fmt.Errorf("upsert failed: %w", err)
			}

			elapsed := time.Since(start).Round(time.Millisecond)
			log.Printf("done in %s — inserted=%d updated=%d total=%d",
				elapsed, inserted, updated, inserted+updated)
			return nil
		},
	}

	cmd.Flags().StringVar(&dataDir, "data-dir", "", "Directory containing exercises.json (required)")
	cmd.Flags().StringVar(&dbURL, "db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string (required unless --dry-run)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Parse records but do not write to the database")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each record as it is processed")

	return cmd
}
