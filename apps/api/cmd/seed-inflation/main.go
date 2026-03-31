// seed-inflation downloads historical CPI inflation data from the World Bank API
// and upserts it into the inflation_data table.
//
// Usage:
//
//	go run ./cmd/seed-inflation --db-url "postgres://requiem:requiem@localhost:5432/requiem"
//	go run ./cmd/seed-inflation --dry-run
//	go run ./cmd/seed-inflation --db-url "..." --verbose
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
	const defaultURL = "https://api.worldbank.org/v2/country/all/indicator/FP.CPI.TOTL.ZG?format=json&per_page=20000&mrv=30"

	dbURL := flag.String("db-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection string (required unless --dry-run)")
	dryRun := flag.Bool("dry-run", false, "Parse and normalise records but do not write to the database")
	verbose := flag.Bool("verbose", false, "Log each record as it is processed")
	url := flag.String("url", defaultURL, "Override World Bank API URL")

	flag.Parse()

	if *dbURL == "" && !*dryRun {
		log.Fatal("--db-url is required when not using --dry-run (or set DATABASE_URL env var)")
	}

	ctx := context.Background()
	start := time.Now()

	log.Printf("downloading from %s", *url)
	records, err := fetchAndParse(*url)
	if err != nil {
		log.Fatalf("fetch failed: %v", err)
	}
	log.Printf("parsed %d records", len(records))

	// Normalise and filter out regional aggregates (empty CountryCode after normalise).
	valid := records[:0]
	for _, r := range records {
		r = normalise(r)
		if r.CountryCode == "" {
			continue
		}
		if *verbose {
			fmt.Printf("  %s  %d  rate=%.4f%%\n", r.CountryCode, r.Year, r.Rate)
		}
		valid = append(valid, r)
	}
	log.Printf("normalised to %d records (filtered %d regional aggregates)", len(valid), len(records)-len(valid))

	if *dryRun {
		printStats(valid)
		log.Println("dry-run: no data written to database")
		return
	}

	conn, err := openDB(ctx, *dbURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	inserted, updated, err := upsertRecords(ctx, conn, valid)
	conn.Close(ctx)

	if err != nil {
		log.Fatalf("upsert failed: %v", err)
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Printf("done in %s — inserted=%d updated=%d total=%d", elapsed, inserted, updated, inserted+updated)
}
