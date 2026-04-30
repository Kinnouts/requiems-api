package words

import (
	"context"
	"fmt"
	"strings"
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

// Definition represents a single definition entry for a word.
type Definition struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definition   string `json:"definition"`
	Example      string `json:"example,omitempty"`
}

// DictionaryEntry is the response payload for the dictionary endpoint.
type DictionaryEntry struct {
	Word        string       `json:"word"`
	Phonetic    string       `json:"phonetic,omitempty"`
	Definitions []Definition `json:"definitions"`
	Synonyms    []string     `json:"synonyms"`
}

func (DictionaryEntry) IsData() {}

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

// Define returns the dictionary entry for the given word.
// Returns an error if the word is not found in the dataset.
func (s *Service) Define(word string) (DictionaryEntry, error) {
	normalized := strings.ToLower(strings.TrimSpace(word))

	e, ok := dictionaryData[normalized]
	if !ok {
		return DictionaryEntry{}, fmt.Errorf("word not found: %s", normalized)
	}

	defs := make([]Definition, 0, len(e.definitions))
	for _, d := range e.definitions {
		defs = append(defs, Definition{
			PartOfSpeech: d.partOfSpeech,
			Definition:   d.definition,
			Example:      d.example,
		})
	}

	synonyms := e.synonyms
	if synonyms == nil {
		synonyms = []string{}
	}

	return DictionaryEntry{
		Word:        normalized,
		Phonetic:    e.phonetic,
		Definitions: defs,
		Synonyms:    synonyms,
	}, nil
}
