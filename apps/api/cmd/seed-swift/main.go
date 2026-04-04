// seed-swift downloads an open-source SWIFT/BIC code dataset and upserts the
// bank records into the swift_codes PostgreSQL table.
//
// Usage:
//
//	go run ./cmd/seed-swift
//	go run ./cmd/seed-swift --dry-run
//	go run ./cmd/seed-swift --db-url "postgres://requiem:requiem@db:5432/requiem"
//	go run ./cmd/seed-swift --verbose
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// defaultDataURL is a community-maintained CSV of SWIFT/BIC codes.
// The --url flag lets operators substitute any compatible CSV source.
const defaultDataURL = "https://raw.githubusercontent.com/ardislu/swift-codes/main/swift-codes.csv"

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
		url     string
	)

	cmd := &cobra.Command{
		Use:   "seed-swift",
		Short: "Seed the swift_codes table from an open-source SWIFT/BIC dataset",
		Long: `Downloads a CSV of SWIFT/BIC codes and upserts the bank records into
the swift_codes PostgreSQL table.

Each record includes the full 11-character BIC, component fields (bank code,
country code, location code, branch code), bank name, city, and country name.

Run inside the API Docker container so the db hostname resolves:

  docker exec requiem-dev-api-1 go run ./cmd/seed-swift
  docker exec requiem-dev-api-1 go run ./cmd/seed-swift --dry-run
  docker exec requiem-dev-api-1 go run ./cmd/seed-swift --verbose`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			ctx := context.Background()
			start := time.Now()

			log.Printf("downloading SWIFT codes from %s", url)
			records, err := fetchAndParse(url)
			if err != nil {
				return fmt.Errorf("fetch failed: %w", err)
			}
			log.Printf("parsed %d SWIFT code records", len(records))

			if verbose {
				for _, r := range records {
					fmt.Printf("  %-11s  bank=%-4s  country=%-2s  location=%-2s  branch=%-3s  %q\n",
						r.SwiftCode, r.BankCode, r.CountryCode, r.LocationCode, r.BranchCode, r.BankName)
				}
			}

			if dryRun {
				printStats(records)
				log.Println("dry-run: no data written to database")
				return nil
			}

			conn, err := openDB(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			inserted, updated, err := upsertRecords(ctx, conn, records)
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
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Parse and print records without writing to the database")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each SWIFT code record as it is processed")
	cmd.Flags().StringVar(&url, "url", defaultDataURL, "Override the SWIFT codes CSV URL")

	return cmd
}
