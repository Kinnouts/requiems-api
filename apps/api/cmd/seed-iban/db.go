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

func upsertRecords(ctx context.Context, conn *pgx.Conn, countries []RawIBANCountry) (inserted, updated int, err error) {
	// 1. Create staging table and clear any rows from a prior run in this session.
	_, err = conn.Exec(ctx, `
		CREATE TEMP TABLE IF NOT EXISTS iban_countries_stage (LIKE iban_countries INCLUDING ALL)
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("create staging table: %w", err)
	}
	if _, err = conn.Exec(ctx, `TRUNCATE TABLE iban_countries_stage`); err != nil {
		return 0, 0, fmt.Errorf("truncate staging table: %w", err)
	}

	// 2. Bulk copy into staging via the PostgreSQL COPY protocol.
	cols := []string{
		"country_code", "country_name", "iban_length", "bban_format",
		"bank_offset", "bank_length", "account_offset", "account_length",
		"sepa_member",
	}

	rows := make([][]any, 0, len(countries))
	for _, c := range countries {
		rows = append(rows, []any{
			c.CountryCode,
			c.CountryName,
			int16(c.IBANLength), //nolint:gosec // IBAN lengths are small positive integers
			c.BBANFormat,
			int16(c.BankOffset()),    //nolint:gosec // bank offsets are small positive integers
			int16(c.BankLength()),    //nolint:gosec // bank lengths are small positive integers
			int16(c.AccountOffset()), //nolint:gosec // account offsets are small positive integers
			int16(c.AccountLength()), //nolint:gosec // account lengths are small positive integers
			c.SEPAMember,
		})
	}

	n, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"iban_countries_stage"},
		cols,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("COPY into staging: %w", err)
	}
	log.Printf("staged %d rows", n)

	// 3. Merge staging → iban_countries. Newer seed always wins on conflict.
	mergeRows, err := conn.Query(ctx, `
		INSERT INTO iban_countries (
			country_code, country_name, iban_length, bban_format,
			bank_offset, bank_length, account_offset, account_length,
			sepa_member
		)
		SELECT
			country_code, country_name, iban_length, bban_format,
			bank_offset, bank_length, account_offset, account_length,
			sepa_member
		FROM iban_countries_stage
		ON CONFLICT (country_code) DO UPDATE SET
			country_name   = EXCLUDED.country_name,
			iban_length    = EXCLUDED.iban_length,
			bban_format    = EXCLUDED.bban_format,
			bank_offset    = EXCLUDED.bank_offset,
			bank_length    = EXCLUDED.bank_length,
			account_offset = EXCLUDED.account_offset,
			account_length = EXCLUDED.account_length,
			sepa_member    = EXCLUDED.sepa_member,
			last_updated   = NOW()
		RETURNING (xmax = 0) AS is_insert
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("merge into iban_countries: %w", err)
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
