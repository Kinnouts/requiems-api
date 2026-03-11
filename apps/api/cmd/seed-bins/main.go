// seed-bins aggregates BIN/IIN data from multiple open-source datasets,
// normalises and merges the records, and upserts them into the bin_data table.
//
// Usage:
//
//	go run ./cmd/seed-bins --db-url "postgres://requiem:requiem@localhost:5432/requiem"
//	go run ./cmd/seed-bins --dry-run
//	go run ./cmd/seed-bins --db-url "..." --verbose
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	const (
		defaultIannuttallURL    = "https://raw.githubusercontent.com/iannuttall/binlist-data/master/binlist-data.csv"
		defaultVenelinkochevURL = "https://raw.githubusercontent.com/venelinkochev/bin-list-data/master/bin-list-data.csv"
	)

	dbURL := flag.String("db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string (required unless --dry-run)")
	dryRun := flag.Bool("dry-run", false, "Parse and normalise records but do not write to the database")
	verbose := flag.Bool("verbose", false, "Log each record as it is processed")
	urlA := flag.String("url-iannuttall", defaultIannuttallURL, "Override URL for iannuttall/binlist-data CSV")
	urlB := flag.String("url-venelinkochev", defaultVenelinkochevURL, "Override URL for venelinkochev/bin-list-data CSV")

	flag.Parse()

	if *dbURL == "" && !*dryRun {
		log.Fatal("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
	}

	sources := []Source{
		{
			Name:       "iannuttall",
			URL:        *urlA,
			Confidence: 0.75,
			Parse:      parseIannuttall,
		},
		{
			Name:       "venelinkochev",
			URL:        *urlB,
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
		log.Fatal("no records loaded from any source — aborting")
	}

	for i := range all {
		all[i] = normalise(all[i])
		if *verbose {
			fmt.Printf("  %s  scheme=%-12s type=%-8s level=%-12s country=%s issuer=%q\n",
				all[i].BINPrefix, all[i].Scheme, all[i].CardType, all[i].CardLevel,
				all[i].CountryCode, all[i].IssuerName)
		}
	}

	merged := mergeRecords(all)
	log.Printf("merged into %d unique BIN prefixes (from %d total records)", len(merged), len(all))

	if *dryRun {
		printStats(merged)
		log.Println("dry-run: no data written to database")
		return
	}

	conn, err := openDB(ctx, *dbURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	inserted, updated, err := upsertRecords(ctx, conn, merged)
	conn.Close(ctx)

	if err != nil {
		log.Fatalf("upsert failed: %v", err)
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Printf("done in %s — inserted=%d updated=%d total=%d", elapsed, inserted, updated, inserted+updated) //nolint:gosec // G706: elapsed is time.Duration, not user input
}
