package counter

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles durable long-term storage for counters.
type Repository interface {
	Upsert(ctx context.Context, namespace string, total int64) error
	Get(ctx context.Context, namespace string) (int64, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Upsert(ctx context.Context, namespace string, total int64) error {
	_, err := r.db.Exec(ctx, `
INSERT INTO counters(namespace, total, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT(namespace)
DO UPDATE SET
  total      = EXCLUDED.total,
  updated_at = NOW()
`, namespace, total)
	return err
}

func (r *postgresRepository) Get(ctx context.Context, namespace string) (int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `SELECT total FROM counters WHERE namespace = $1`, namespace).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
