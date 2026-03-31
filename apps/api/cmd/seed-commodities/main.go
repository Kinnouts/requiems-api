// seed-commodities downloads historical annual commodity price data from FRED
// (Federal Reserve Economic Data) and upserts it into the commodity_price_history table.
//
// Data source: https://fred.stlouisfed.org — free public API, no key required.
// Each commodity maps to a FRED series. Daily/monthly series are averaged per year.
//
// Usage:
//
//	go run ./cmd/seed-commodities
//	go run ./cmd/seed-commodities --dry-run
//	go run ./cmd/seed-commodities --db-url "postgres://requiem:requiem@db:5432/requiem"
//	go run ./cmd/seed-commodities --verbose
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		dbURL   string
		dryRun  bool
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "seed-commodities",
		Short: "Seed the commodity_price_history table from FRED",
		Long: `Downloads historical annual commodity price data from FRED (Federal Reserve
Economic Data) and upserts it into the commodity_price_history PostgreSQL table.

Each commodity is mapped to a FRED series ID. Daily and monthly series are
averaged per calendar year. Annual averages going back up to 60 years are stored.

Run inside the API Docker container so the db hostname resolves:

  docker exec requiem-dev-api-1 go run ./cmd/seed-commodities
  docker compose exec api go run ./cmd/seed-commodities`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			ctx := context.Background()
			start := time.Now()

			var allRecords []CommodityRecord

			for _, cfg := range commodities {
				log.Printf("[%s] fetching FRED series %s", cfg.Slug, cfg.SeriesID)
				recs, err := fetchAndAggregate(cfg)
				if err != nil {
					log.Printf("[%s] WARN: skipping — %v", cfg.Slug, err)
					continue
				}
				log.Printf("[%s] parsed %d annual averages", cfg.Slug, len(recs))
				if verbose {
					for _, r := range recs {
						fmt.Printf("  %-20s %d  %.4f %s\n", r.Slug, r.Year, r.Price, r.Unit)
					}
				}
				allRecords = append(allRecords, recs...)
			}

			if len(allRecords) == 0 {
				return fmt.Errorf("no records loaded from any source")
			}

			log.Printf("total: %d records across %d commodities", len(allRecords), len(commodities))

			if dryRun {
				printStats(allRecords)
				log.Println("dry-run: no data written to database")
				return nil
			}

			conn, err := openDB(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			inserted, updated, err := upsertRecords(ctx, conn, allRecords)
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
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Parse and aggregate records but do not write to the database")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each record as it is processed")

	return cmd
}
