package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// openDB opens a single PostgreSQL connection for the seed operation.
func openDB(ctx context.Context, dbURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return conn, nil
}

// upsertRecords inserts or updates all records in commodity_price_history.
// Returns the number of inserted and updated rows.
func upsertRecords(ctx context.Context, conn *pgx.Conn, records []CommodityRecord) (inserted, updated int, err error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const query = `
		INSERT INTO commodity_price_history (slug, name, unit, currency, year, price, source, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'fred', NOW())
		ON CONFLICT (slug, year) DO UPDATE
		    SET name       = EXCLUDED.name,
		        unit       = EXCLUDED.unit,
		        currency   = EXCLUDED.currency,
		        price      = EXCLUDED.price,
		        source     = EXCLUDED.source,
		        updated_at = EXCLUDED.updated_at
		RETURNING (xmax = 0) AS is_insert
	`

	for _, r := range records {
		var isInsert bool
		row := tx.QueryRow(ctx, query, r.Slug, r.Name, r.Unit, r.Currency, r.Year, r.Price)
		if scanErr := row.Scan(&isInsert); scanErr != nil {
			return 0, 0, fmt.Errorf("upsert (%s, %d): %w", r.Slug, r.Year, scanErr)
		}
		if isInsert {
			inserted++
		} else {
			updated++
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("commit: %w", err)
	}
	return inserted, updated, nil
}
