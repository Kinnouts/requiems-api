// seed-bins aggregates BIN/IIN data from multiple open-source datasets,
// normalises and merges the records, and upserts them into the bin_data table.
//
// Usage:
//
//	go run ./cmd/seed-bins
//	go run ./cmd/seed-bins --dry-run
//	go run ./cmd/seed-bins --db-url "postgres://requiem:requiem@db:5432/requiem"
//	go run ./cmd/seed-bins --verbose
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
	const (
		defaultIannuttallURL    = "https://raw.githubusercontent.com/iannuttall/binlist-data/master/binlist-data.csv"
		defaultVenelinkochevURL = "https://raw.githubusercontent.com/venelinkochev/bin-list-data/master/bin-list-data.csv"
	)

	var (
		dbURL   string
		dryRun  bool
		verbose bool
		urlA    string
		urlB    string
	)

	cmd := &cobra.Command{
		Use:   "seed-bins",
		Short: "Seed the bin_data table from open-source BIN/IIN datasets",
		Long: `Downloads BIN/IIN data from two open-source CSV datasets, normalises and
merges the records, and upserts them into the bin_data PostgreSQL table.

Run inside the API Docker container so the db hostname resolves:

  docker exec requiem-dev-api-1 go run ./cmd/seed-bins
  docker compose exec api go run ./cmd/seed-bins`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" && !dryRun {
				return fmt.Errorf("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
			}

			sources := []Source{
				{
					Name:       "iannuttall",
					URL:        urlA,
					Confidence: 0.75,
					Parse:      parseIannuttall,
				},
				{
					Name:       "venelinkochev",
					URL:        urlB,
					Confidence: 0.80,
					Parse:      parseVenelinkochev,
				},
			}

			ctx := context.Background()
			start := time.Now()

			all := make([]RawBINRecord, 0, 500_000)

			for _, src := range sources {
				log.Printf("[%s] downloading from %s", src.Name, src.URL)
				records, err := fetchAndParse(src)
				if err != nil {
					log.Printf("[%s] WARN: skipping source: %v", src.Name, err)
					continue
				}
				log.Printf("[%s] parsed %d records", src.Name, len(records))
				all = append(all, records...)
			}

			if len(all) == 0 {
				return fmt.Errorf("no records loaded from any source")
			}

			for i := range all {
				all[i] = normalise(all[i])
				if verbose {
					fmt.Printf("  %s  scheme=%-12s type=%-8s level=%-12s country=%s issuer=%q\n",
						all[i].BINPrefix, all[i].Scheme, all[i].CardType, all[i].CardLevel,
						all[i].CountryCode, all[i].IssuerName)
				}
			}

			merged := mergeRecords(all)
			log.Printf("merged into %d unique BIN prefixes (from %d total records)", len(merged), len(all))

			if dryRun {
				printStats(merged)
				log.Println("dry-run: no data written to database")
				return nil
			}

			conn, err := seedutil.OpenDB(ctx, dbURL)
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			inserted, updated, err := upsertRecords(ctx, conn, merged)
			conn.Close(ctx)
			if err != nil {
				return fmt.Errorf("upsert failed: %w", err)
			}

			elapsed := time.Since(start).Round(time.Millisecond)
			log.Printf("done in %s — inserted=%d updated=%d total=%d", elapsed, inserted, updated, inserted+updated) //nolint:gosec // G706: elapsed is time.Duration, not user input
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string (required unless --dry-run)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Parse and normalise records but do not write to the database")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Log each record as it is processed")
	cmd.Flags().StringVar(&urlA, "url-iannuttall", defaultIannuttallURL, "Override URL for iannuttall/binlist-data CSV")
	cmd.Flags().StringVar(&urlB, "url-venelinkochev", defaultVenelinkochevURL, "Override URL for venelinkochev/bin-list-data CSV")

	return cmd
}
