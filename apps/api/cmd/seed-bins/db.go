package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func openDB(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)

	if err != nil {
		return nil, fmt.Errorf("pgx.Connect: %w", err)
	}

	return conn, nil
}

// upsertRecords loads all records into a temporary staging table and then
// merges them into bin_data using INSERT ON CONFLICT. Returns the number of
// rows inserted and updated.
func upsertRecords(ctx context.Context, conn *pgx.Conn, records map[string]RawBINRecord) (inserted, updated int, err error) {
	// 1. Create staging table.
	// Note: no ON COMMIT DROP — that would drop the table at the end of the
	// implicit single-statement transaction, before COPY runs.
	// The temp table is dropped automatically when the session ends.
	_, err = conn.Exec(ctx, `
		CREATE TEMP TABLE IF NOT EXISTS bin_data_stage (LIKE bin_data INCLUDING ALL)
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("create staging table: %w", err)
	}

	// 2. Bulk copy into staging via the PostgreSQL COPY protocol.
	cols := []string{
		"bin_prefix", "prefix_length", "scheme", "card_type", "card_level",
		"issuer_name", "issuer_url", "issuer_phone",
		"country_code", "country_name",
		"prepaid", "source", "confidence",
	}

	rows := make([][]any, 0, len(records))
	for _, r := range records {
		rows = append(rows, []any{
			r.BINPrefix,
			int16(len(r.BINPrefix)),
			r.Scheme,
			r.CardType,
			r.CardLevel,
			r.IssuerName,
			r.IssuerURL,
			r.IssuerPhone,
			r.CountryCode,
			r.CountryName,
			r.Prepaid,
			r.Source,
			r.Confidence,
		})
	}

	n, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"bin_data_stage"},
		cols,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("COPY into staging: %w", err)
	}
	log.Printf("staged %d rows", n)

	// 3. Merge staging → bin_data.
	// We count pre-existing rows to calculate insert vs update split.
	var existingCount int
	if err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM bin_data`).Scan(&existingCount); err != nil {
		return 0, 0, fmt.Errorf("count existing: %w", err)
	}

	_, err = conn.Exec(ctx, `
		INSERT INTO bin_data (
			bin_prefix, prefix_length, scheme, card_type, card_level,
			issuer_name, issuer_url, issuer_phone,
			country_code, country_name,
			prepaid, source, confidence
		)
		SELECT
			bin_prefix, prefix_length, scheme, card_type, card_level,
			issuer_name, issuer_url, issuer_phone,
			country_code, country_name,
			prepaid, source, confidence
		FROM bin_data_stage
		ON CONFLICT (bin_prefix) DO UPDATE SET
			prefix_length = EXCLUDED.prefix_length,
			scheme        = CASE
				WHEN EXCLUDED.confidence >= bin_data.confidence THEN EXCLUDED.scheme
				ELSE bin_data.scheme
			END,
			card_type     = CASE
				WHEN EXCLUDED.confidence >= bin_data.confidence THEN EXCLUDED.card_type
				ELSE bin_data.card_type
			END,
			card_level    = CASE
				WHEN EXCLUDED.confidence >= bin_data.confidence THEN EXCLUDED.card_level
				ELSE bin_data.card_level
			END,
			issuer_name   = CASE WHEN EXCLUDED.issuer_name   <> '' THEN EXCLUDED.issuer_name   ELSE bin_data.issuer_name   END,
			issuer_url    = CASE WHEN EXCLUDED.issuer_url    <> '' THEN EXCLUDED.issuer_url    ELSE bin_data.issuer_url    END,
			issuer_phone  = CASE WHEN EXCLUDED.issuer_phone  <> '' THEN EXCLUDED.issuer_phone  ELSE bin_data.issuer_phone  END,
			country_code  = CASE WHEN EXCLUDED.country_code  <> '' THEN EXCLUDED.country_code  ELSE bin_data.country_code  END,
			country_name  = CASE WHEN EXCLUDED.country_name  <> '' THEN EXCLUDED.country_name  ELSE bin_data.country_name  END,
			prepaid       = EXCLUDED.prepaid,
			source        = EXCLUDED.source,
			confidence    = GREATEST(bin_data.confidence, EXCLUDED.confidence),
			last_updated  = NOW()
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("merge into bin_data: %w", err)
	}

	var totalCount int
	if err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM bin_data`).Scan(&totalCount); err != nil {
		return 0, 0, fmt.Errorf("count final: %w", err)
	}

	inserted = totalCount - existingCount
	if inserted < 0 {
		inserted = 0
	}
	updated = len(records) - inserted

	return inserted, updated, nil
}
