package words

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
)

// mockRow implements pgx.Row.
type mockRow struct {
	scanFn func(dest ...any) error
}

func (m *mockRow) Scan(dest ...any) error { return m.scanFn(dest...) }

// mockQuerier implements querier.
type mockQuerier struct {
	row pgx.Row
}

func (m *mockQuerier) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return m.row
}

func newTestService(row pgx.Row) *Service {
	return &Service{db: &mockQuerier{row: row}}
}

func TestRandom_EmptyTable(t *testing.T) {
	svc := newTestService(&mockRow{
		scanFn: func(_ ...any) error { return pgx.ErrNoRows },
	})

	_, err := svc.Random(context.Background())
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}
}

func TestRandom_SingleRow(t *testing.T) {
	svc := newTestService(&mockRow{
		scanFn: func(dest ...any) error {
			*dest[0].(*int) = 42
			*dest[1].(*string) = "ephemeral"
			*dest[2].(*string) = "Lasting for a very short time."
			*dest[3].(*string) = "adjective"
			return nil
		},
	})

	got, err := svc.Random(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("expected ID 42, got %d", got.ID)
	}
	if got.Word != "ephemeral" {
		t.Errorf("unexpected word: %q", got.Word)
	}
	if got.Definition != "Lasting for a very short time." {
		t.Errorf("unexpected definition: %q", got.Definition)
	}
	if got.PartOfSpeech != "adjective" {
		t.Errorf("unexpected part_of_speech: %q", got.PartOfSpeech)
	}
}

func TestRandom_ScanError(t *testing.T) {
	scanErr := errors.New("scan failed")
	svc := newTestService(&mockRow{
		scanFn: func(_ ...any) error { return scanErr },
	})

	_, err := svc.Random(context.Background())
	if !errors.Is(err, scanErr) {
		t.Errorf("expected scan error, got %v", err)
	}
}

func TestRandom_ReturnsZeroValueOnError(t *testing.T) {
	svc := newTestService(&mockRow{
		scanFn: func(_ ...any) error { return pgx.ErrNoRows },
	})

	got, _ := svc.Random(context.Background())
	if got.ID != 0 || got.Word != "" || got.Definition != "" || got.PartOfSpeech != "" {
		t.Errorf("expected zero Word on error, got %+v", got)
	}
}

func TestRandom_EmptyPartOfSpeechAllowed(t *testing.T) {
	svc := newTestService(&mockRow{
		scanFn: func(dest ...any) error {
			*dest[0].(*int) = 1
			*dest[1].(*string) = "run"
			*dest[2].(*string) = "Move at a speed faster than walking."
			*dest[3].(*string) = ""
			return nil
		},
	})

	got, err := svc.Random(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.PartOfSpeech != "" {
		t.Errorf("expected empty part_of_speech, got %q", got.PartOfSpeech)
	}
}
