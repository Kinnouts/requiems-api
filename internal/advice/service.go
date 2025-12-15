package advice

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Advice struct {
	ID   int    `json:"id"`
	Text string `json:"advice"`
}

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Random(ctx context.Context) (Advice, error) {
	row := s.db.QueryRow(ctx, `
SELECT id, text
FROM advice
ORDER BY random()
LIMIT 1;
`)

	var a Advice
	if err := row.Scan(&a.ID, &a.Text); err != nil {
		return Advice{}, err
	}
	return a, nil
}


