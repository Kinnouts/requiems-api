// seed-inflation downloads historical CPI inflation data from the World Bank API
// and upserts it into the inflation_data table.
//
// Usage:
//
//	go run ./cmd/seed-inflation
//	go run ./cmd/seed-inflation --dry-run
//	go run ./cmd/seed-inflation --db-url "postgres://requiem:requiem@db:5432/requiem"
//	go run ./cmd/seed-inflation --verbose
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
	const defaultURL = "https://api.worldbank.org/v2/country/all/indicator/FP.CPI.TOTL.ZG?format=json&per_page=20000&mrv=30"

	var (
		dbURL   string
		dryRun  bool
		verbose bool
		url     string
	)

	cmd := &cobra.Command{
		Use:   "seed-inflation",
		Short: "Seed the inflation_data table from the World Bank API",
		Long: `Downloads historical CPI inflation data from the World Bank API and upserts
it into the inflation_data PostgreSQL table (~241 countries, last 30 years).

Run inside the API Docker container so the db hostname resolves:

  docker exec requiem-dev-api-1 go run ./cmd/seed-inflation
  docker compose exec api go run ./cmd/seed-inflation`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			ctx := context.Background()
			start := time.Now()

			log.Printf("downloading from %s", url)
			records, err := fetchAndParse(url)
			if err != nil {
				return fmt.Errorf("fetch failed: %w", err)
			}
			log.Printf("parsed %d records", len(records))

			// Normalise and filter out regional aggregates (empty CountryCode after normalise).
			valid := records[:0]
			for _, r := range records {
				r = normalise(r)
				if r.CountryCode == "" {
					continue
				}
				if verbose {
					fmt.Printf("  %s  %d  rate=%.4f%%\n", r.CountryCode, r.Year, r.Rate)
				}
				valid = append(valid, r)
			}
			log.Printf("normalised to %d records (filtered %d regional aggregates)", len(valid), len(records)-len(valid))

			if dryRun {
				printStats(valid)
				log.Println("dry-run: no data written to database")
				return nil
			}

			conn, err := seedutil.OpenDB(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			inserted, updated, err := upsertRecords(ctx, conn, valid)
			conn.Close(ctx)
			if err != nil {
				return fmt.Errorf("upsert failed: %w", err)
			}

			elapsed := time.Since(start).Round(time.Millisecond)
			log.Printf("done in %s — inserted=%d updated=%d total=%d", elapsed, inserted, updated, inserted+updated)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string (required unless --dry-run)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Parse and normalise records but do not write to the database")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each record as it is processed")
	cmd.Flags().StringVar(&url, "url", defaultURL, "Override World Bank API URL")

	return cmd
}
