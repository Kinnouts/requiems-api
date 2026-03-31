// seed-iban downloads the IBAN country registry from the php-iban project
// (a maintained mirror of the official SWIFT IBAN Registry) and upserts the
// country format data — including bank identifier and account number positions
// — into the iban_countries PostgreSQL table.
//
// Usage:
//
//	go run ./cmd/seed-iban
//	go run ./cmd/seed-iban --dry-run
//	go run ./cmd/seed-iban --db-url "postgres://requiem:requiem@db:5432/requiem"
//	go run ./cmd/seed-iban --verbose
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// defaultRegistryURL is the php-iban project's IBAN registry file, which is
// a maintained machine-readable mirror of the official SWIFT IBAN Registry.
// Source: https://github.com/globalcitizen/php-iban
const defaultRegistryURL = "https://raw.githubusercontent.com/globalcitizen/php-iban/master/registry.txt"

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
		Use:   "seed-iban",
		Short: "Seed the iban_countries table from the SWIFT IBAN Registry",
		Long: `Downloads the IBAN country registry from the php-iban project
(a maintained mirror of the official SWIFT IBAN Registry) and upserts
country format data into the iban_countries PostgreSQL table.

Each record includes the country name, expected IBAN length, BBAN format
string, and the 0-indexed offsets for extracting the bank identifier and
account number from the BBAN.

Run inside the API Docker container so the db hostname resolves:

  docker exec requiem-dev-api-1 go run ./cmd/seed-iban
  docker exec requiem-dev-api-1 go run ./cmd/seed-iban --dry-run
  docker exec requiem-dev-api-1 go run ./cmd/seed-iban --verbose`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			ctx := context.Background()
			start := time.Now()

			log.Printf("downloading registry from %s", url)
			countries, err := fetchAndParse(url)
			if err != nil {
				return fmt.Errorf("fetch failed: %w", err)
			}
			log.Printf("parsed %d country records", len(countries))

			if verbose {
				for _, c := range countries {
					fmt.Printf("  %-2s  %-32s  iban_len=%-3d  bank=%d+%d  acct=%d+%d  sepa=%v\n",
						c.CountryCode, c.CountryName, c.IBANLength,
						c.BankOffset(), c.BankLength(),
						c.AccountOffset(), c.AccountLength(),
						c.SEPAMember,
					)
				}
			}

			if dryRun {
				printStats(countries)
				log.Println("dry-run: no data written to database")
				return nil
			}

			conn, err := openDB(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			inserted, updated, err := upsertRecords(ctx, conn, countries)
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
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each country record as it is processed")
	cmd.Flags().StringVar(&url, "url", defaultRegistryURL, "Override the IBAN registry URL")

	return cmd
}
