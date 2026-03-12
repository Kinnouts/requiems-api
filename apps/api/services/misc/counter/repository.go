package counter

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Durable long-term storage for counters.
type Repository interface {
	Upsert(ctx context.Context, namespace string, total int64) error
	UpsertBatch(ctx context.Context, counters map[string]int64) error
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

func (r *postgresRepository) UpsertBatch(ctx context.Context, counters map[string]int64) error {
	if len(counters) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for namespace, total := range counters {
		batch.Queue(`
INSERT INTO counters(namespace, total, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT(namespace)
DO UPDATE SET
  total      = EXCLUDED.total,
  updated_at = NOW()
`, namespace, total)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		if _, err := results.Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (r *postgresRepository) Get(ctx context.Context, namespace string) (int64, error) {
	var total int64

	err := r.db.QueryRow(ctx, `SELECT total FROM counters WHERE namespace = $1`, namespace).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
