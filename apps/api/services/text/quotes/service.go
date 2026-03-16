package quotes

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Quote struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author,omitempty"`
}

func (Quote) IsData() {}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Service struct {
	db querier
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Random(ctx context.Context) (Quote, error) {
	row := s.db.QueryRow(ctx, `
SELECT id, text, author
FROM quotes
ORDER BY random()
LIMIT 1;
`)

	var q Quote

	if err := row.Scan(&q.ID, &q.Text, &q.Author); err != nil {
		return Quote{}, err
	}

	return q, nil
}
