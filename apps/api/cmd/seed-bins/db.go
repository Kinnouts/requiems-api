package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"requiems-api/cmd/internal/seedutil"
)

func upsertRecords(ctx context.Context, conn *pgx.Conn, records map[string]RawBINRecord) (inserted, updated int, err error) {
	// 1. Create staging table and clear any rows from a prior run in this session.
	_, err = conn.Exec(ctx, `
		CREATE TEMP TABLE IF NOT EXISTS bin_data_stage (LIKE bin_data INCLUDING ALL)
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("create staging table: %w", err)
	}
	if _, err = conn.Exec(ctx, `TRUNCATE TABLE bin_data_stage`); err != nil {
		return 0, 0, fmt.Errorf("truncate staging table: %w", err)
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
		prefixLen, convErr := seedutil.ToInt16(len(r.BINPrefix), "prefix_length")
		if convErr != nil {
			return 0, 0, convErr
		}

		rows = append(rows, []any{
			r.BINPrefix,
			prefixLen,
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
	log.Printf("staged %d rows", n) //nolint:gosec // G706: n is an int64 from CopyFrom, not user input

	// 3. Merge staging → bin_data.
	mergeRows, err := conn.Query(ctx, `
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
				WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.scheme    <> '' THEN EXCLUDED.scheme
				ELSE bin_data.scheme
			END,
			card_type     = CASE
				WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.card_type <> '' THEN EXCLUDED.card_type
				ELSE bin_data.card_type
			END,
			card_level    = CASE
				WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.card_level <> '' THEN EXCLUDED.card_level
				ELSE bin_data.card_level
			END,
			issuer_name   = CASE WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.issuer_name   <> '' THEN EXCLUDED.issuer_name   ELSE bin_data.issuer_name   END,
			issuer_url    = CASE WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.issuer_url    <> '' THEN EXCLUDED.issuer_url    ELSE bin_data.issuer_url    END,
			issuer_phone  = CASE WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.issuer_phone  <> '' THEN EXCLUDED.issuer_phone  ELSE bin_data.issuer_phone  END,
			country_code  = CASE WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.country_code  <> '' THEN EXCLUDED.country_code  ELSE bin_data.country_code  END,
			country_name  = CASE WHEN EXCLUDED.confidence >= bin_data.confidence AND EXCLUDED.country_name  <> '' THEN EXCLUDED.country_name  ELSE bin_data.country_name  END,
			prepaid       = CASE WHEN EXCLUDED.confidence >= bin_data.confidence THEN EXCLUDED.prepaid ELSE bin_data.prepaid END,
			source        = EXCLUDED.source,
			confidence    = GREATEST(bin_data.confidence, EXCLUDED.confidence),
			last_updated  = NOW()
		RETURNING (xmax = 0) AS is_insert
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("merge into bin_data: %w", err)
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
