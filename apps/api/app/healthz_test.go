package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"requiems-api/platform/httpx"
)

type mockPinger struct {
	err error
}

func (m mockPinger) Ping(_ context.Context) error {
	return m.err
}

func TestHealthz_DBAvailable(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", http.NoBody)
	w := httptest.NewRecorder()
	Healthz(mockPinger{}).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[healthzResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Data.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Data.Status)
	}
}

func TestHealthz_DBUnavailable(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", http.NoBody)
	w := httptest.NewRecorder()
	Healthz(mockPinger{err: errors.New("connection refused")}).ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Error != "db_unavailable" {
		t.Errorf("expected error code 'db_unavailable', got %q", resp.Error)
	}
}
