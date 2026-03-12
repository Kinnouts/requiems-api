package counter

import (
	"context"
	"fmt"
	"sort"
	"strings"

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

	namespaces := make([]string, 0, len(counters))
	for namespace := range counters {
		namespaces = append(namespaces, namespace)
	}
	sort.Strings(namespaces)

	placeholders := make([]string, 0, len(namespaces))
	args := make([]any, 0, len(namespaces)*2)
	for i, namespace := range namespaces {
		argOffset := i*2 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, NOW())", argOffset, argOffset+1))
		args = append(args, namespace, counters[namespace])
	}

	query := `
INSERT INTO counters(namespace, total, updated_at)
VALUES ` + strings.Join(placeholders, ", ") + `
ON CONFLICT(namespace)
DO UPDATE SET
  total      = EXCLUDED.total,
  updated_at = NOW()
`

	_, err := r.db.Exec(ctx, query, args...)
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
