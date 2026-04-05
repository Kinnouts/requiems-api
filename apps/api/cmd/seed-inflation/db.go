package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"requiems-api/cmd/internal/seedutil"
)


func upsertRecords(ctx context.Context, conn *pgx.Conn, records []RawInflationRecord) (inserted, updated int, err error) {
	// 1. Create staging table and clear any rows from a prior run in this session.
	_, err = conn.Exec(ctx, `
		CREATE TEMP TABLE IF NOT EXISTS inflation_data_stage (LIKE inflation_data INCLUDING ALL)
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("create staging table: %w", err)
	}
	if _, err = conn.Exec(ctx, `TRUNCATE TABLE inflation_data_stage`); err != nil {
		return 0, 0, fmt.Errorf("truncate staging table: %w", err)
	}

	// 2. Bulk copy into staging via the PostgreSQL COPY protocol.
	cols := []string{"country_code", "country_name", "year", "rate", "source"}

	rows := make([][]any, 0, len(records))
	for _, r := range records {
		year, convErr := seedutil.ToInt16(r.Year, "year")
		if convErr != nil {
			return 0, 0, convErr
		}

		rows = append(rows, []any{
			r.CountryCode,
			r.CountryName,
			year,
			r.Rate,
			r.Source,
		})
	}

	n, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"inflation_data_stage"},
		cols,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("COPY into staging: %w", err)
	}
	log.Printf("staged %d rows", n)

	// 3. Merge staging → inflation_data. Newer seed always wins.
	mergeRows, err := conn.Query(ctx, `
		INSERT INTO inflation_data (country_code, country_name, year, rate, source)
		SELECT country_code, country_name, year, rate, source
		FROM inflation_data_stage
		ON CONFLICT (country_code, year) DO UPDATE SET
			country_name = EXCLUDED.country_name,
			rate         = EXCLUDED.rate,
			source       = EXCLUDED.source,
			last_updated = NOW()
		RETURNING (xmax = 0) AS is_insert
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("merge into inflation_data: %w", err)
	}
	defer mergeRows.Close()

	for mergeRows.Next() {
		var isInsert bool
		if err = mergeRows.Scan(&isInsert); err != nil {
			return 0, 0, fmt.Errorf("scan returning: %w", err)
		}
		if isInsert {
			inserted++
		} else {
			updated++
		}
	}
	if err = mergeRows.Err(); err != nil {
		return 0, 0, fmt.Errorf("rows: %w", err)
	}

	return inserted, updated, nil
}
