package words

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Word struct {
	ID           int    `json:"id"`
	Word         string `json:"word"`
	Definition   string `json:"definition"`
	PartOfSpeech string `json:"part_of_speech,omitempty"`
}

func (Word) IsData() {}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Service struct {
	db querier
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Random(ctx context.Context) (Word, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := s.db.QueryRow(ctx, `
SELECT id, word, definition, part_of_speech
FROM words
ORDER BY random()
LIMIT 1;
`)

	var w Word
	if err := row.Scan(&w.ID, &w.Word, &w.Definition, &w.PartOfSpeech); err != nil {
		return Word{}, fmt.Errorf("scan word: %w", err)
	}

	return w, nil
}
