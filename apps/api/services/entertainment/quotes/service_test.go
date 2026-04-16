package quotes

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
			*dest[0].(*int) = 7
			*dest[1].(*string) = "Be yourself; everyone else is already taken."
			*dest[2].(*string) = "Oscar Wilde"
			return nil
		},
	})

	got, err := svc.Random(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 7 {
		t.Errorf("expected ID 7, got %d", got.ID)
	}
	if got.Text != "Be yourself; everyone else is already taken." {
		t.Errorf("unexpected text: %q", got.Text)
	}
	if got.Author != "Oscar Wilde" {
		t.Errorf("unexpected author: %q", got.Author)
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
	if got.ID != 0 || got.Text != "" || got.Author != "" {
		t.Errorf("expected zero Quote on error, got %+v", got)
	}
}

func TestRandom_EmptyAuthorAllowed(t *testing.T) {
	svc := newTestService(&mockRow{
		scanFn: func(dest ...any) error {
			*dest[0].(*int) = 3
			*dest[1].(*string) = "Anonymous wisdom."
			*dest[2].(*string) = ""
			return nil
		},
	})

	got, err := svc.Random(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Author != "" {
		t.Errorf("expected empty author, got %q", got.Author)
	}
}
