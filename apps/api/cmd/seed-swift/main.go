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

	"requiems-api/cmd/internal/seedutil"
)

const (
	defaultDataSource = "https://raw.githubusercontent.com/maranemil/swift-bic-codes-allinone/master/csv/AllCountries_v1.csv"
	defaultCountries  = "https://raw.githubusercontent.com/maranemil/swift-bic-codes-allinone/master/csv/countries.csv"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		dbURL      string
		dryRun     bool
		verbose    bool
		sourcePath string
		countries  string
		url        string
	)

	cmd := &cobra.Command{
		Use:   "seed-swift",
		Short: "Seed the swift_codes table from an upstream SWIFT/BIC dataset",
		Long: `Downloads and normalises SWIFT/BIC code data, then upserts the bank records into
the swift_codes PostgreSQL table.

Each record includes the full 11-character BIC, component fields (bank code,
country code, location code, branch code), bank name, city, and country name.

Run inside the API Docker container so the db hostname resolves:

	docker compose exec api /app/seed-swift
	docker compose exec api /app/seed-swift --dry-run
	docker compose exec api /app/seed-swift --verbose`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			ctx := context.Background()
			start := time.Now()

			effectiveSource := sourcePath
			if url != "" {
				effectiveSource = url
			}

			log.Printf("loading SWIFT codes from %s", effectiveSource)
			records, err := fetchAndParseWithCountries(effectiveSource, countries)
			if err != nil {
				return fmt.Errorf("load failed: %w", err)
			}
			log.Printf("parsed %d SWIFT code records", len(records))
			records = dedupeBySwiftCode(records)
			log.Printf("deduped to %d unique SWIFT code records", len(records))

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
			log.Printf("done in %s — inserted=%d updated=%d total=%d", elapsed, inserted, updated, inserted+updated)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string (required unless --dry-run)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Parse and print records without writing to the database")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each SWIFT code record as it is processed")
	cmd.Flags().StringVar(&sourcePath, "source", defaultDataSource, "Path or URL for SWIFT codes CSV source")
	cmd.Flags().StringVar(&countries, "countries", defaultCountries, "Path or URL for country code mapping CSV")
	cmd.Flags().StringVar(&url, "url", "", "Optional SWIFT codes CSV URL (overrides --source)")

	return cmd
}

func dedupeBySwiftCode(records []RawSWIFTRecord) []RawSWIFTRecord {
	seen := make(map[string]int, len(records))
	result := make([]RawSWIFTRecord, 0, len(records))

	for _, r := range records {
		if idx, ok := seen[r.SwiftCode]; ok {
			existing := result[idx]
			if existing.BankName == "" && r.BankName != "" {
				existing.BankName = r.BankName
			}
			if existing.City == "" && r.City != "" {
				existing.City = r.City
			}
			if existing.CountryName == "" && r.CountryName != "" {
				existing.CountryName = r.CountryName
			}
			result[idx] = existing
			continue
		}

		seen[r.SwiftCode] = len(result)
		result = append(result, r)
	}

	return result
}
