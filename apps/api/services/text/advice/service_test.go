package advice

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
			*dest[0].(*int) = 1
			*dest[1].(*string) = "Do one thing every day that scares you."
			return nil
		},
	})

	got, err := svc.Random(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}
	if got.Text != "Do one thing every day that scares you." {
		t.Errorf("unexpected text: %q", got.Text)
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
	if got.ID != 0 || got.Text != "" {
		t.Errorf("expected zero Advice on error, got %+v", got)
	}
}
