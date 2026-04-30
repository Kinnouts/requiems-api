package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func upsertRecords(ctx context.Context, conn *pgx.Conn, records []RawSWIFTRecord) (inserted, updated int, err error) {
	// 1. Create staging table and clear any rows from a prior run in this session.
	_, err = conn.Exec(ctx, `
		CREATE TEMP TABLE IF NOT EXISTS swift_codes_stage (LIKE swift_codes INCLUDING ALL)
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("create staging table: %w", err)
	}
	if _, err = conn.Exec(ctx, `TRUNCATE TABLE swift_codes_stage`); err != nil {
		return 0, 0, fmt.Errorf("truncate staging table: %w", err)
	}

	// 2. Bulk copy into staging via the PostgreSQL COPY protocol.
	cols := []string{
		"swift_code", "bank_code", "country_code", "location_code", "branch_code",
		"bank_name", "city", "country_name",
	}

	rows := make([][]any, 0, len(records))
	for _, r := range records {
		rows = append(rows, []any{
			r.SwiftCode,
			r.BankCode,
			r.CountryCode,
			r.LocationCode,
			r.BranchCode,
			r.BankName,
			r.City,
			r.CountryName,
		})
	}

	n, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"swift_codes_stage"},
		cols,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("COPY into staging: %w", err)
	}
	log.Printf("staged %d rows", n)

	// 3. Merge staging → swift_codes. Newer seed always wins on conflict.
	mergeRows, err := conn.Query(ctx, `
		INSERT INTO swift_codes (
			swift_code, bank_code, country_code, location_code, branch_code,
			bank_name, city, country_name
		)
		SELECT
			swift_code, bank_code, country_code, location_code, branch_code,
			bank_name, city, country_name
		FROM swift_codes_stage
		ON CONFLICT (swift_code) DO UPDATE SET
			bank_name    = EXCLUDED.bank_name,
			city         = EXCLUDED.city,
			country_name = EXCLUDED.country_name,
			last_updated = NOW()
		RETURNING (xmax = 0) AS is_insert
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("merge into swift_codes: %w", err)
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
