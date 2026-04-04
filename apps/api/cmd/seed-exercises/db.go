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

// upsertRecords inserts or updates all records in the exercises table.
// Returns the number of inserted and updated rows.
func upsertRecords(ctx context.Context, conn *pgx.Conn, records []ExerciseRecord) (inserted, updated int, err error) {
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
		INSERT INTO exercises (external_id, name, body_parts, equipment, target_muscles, secondary_muscles, instructions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (external_id) DO UPDATE SET
		    name              = EXCLUDED.name,
		    body_parts        = EXCLUDED.body_parts,
		    equipment         = EXCLUDED.equipment,
		    target_muscles    = EXCLUDED.target_muscles,
		    secondary_muscles = EXCLUDED.secondary_muscles,
		    instructions      = EXCLUDED.instructions
		RETURNING (xmax = 0) AS is_insert
	`

	for _, r := range records {
		var isInsert bool
		row := tx.QueryRow(ctx, query,
			r.ExternalID,
			r.Name,
			r.BodyParts,
			r.Equipment,
			r.TargetMuscles,
			r.SecondaryMuscles,
			r.Instructions,
		)
		if scanErr := row.Scan(&isInsert); scanErr != nil {
			return 0, 0, fmt.Errorf("upsert (%s): %w", r.ExternalID, scanErr)
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
