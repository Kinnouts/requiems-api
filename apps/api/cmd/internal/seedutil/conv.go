package seedutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// OpenDB opens a PostgreSQL connection using the given DSN.
func OpenDB(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)

	if err != nil {
		return nil, fmt.Errorf("pgx.Connect: %w", err)
	}
	
	return conn, nil
}

// ToInt16 converts v to int16, returning an error if the value is outside the
// valid int16 range [-32768, 32767].
func ToInt16(v int, field string) (int16, error) {
	if v < -32768 || v > 32767 {
		return 0, fmt.Errorf("%s out of int16 range: %d", field, v)
	}

	return int16(v), nil
}
