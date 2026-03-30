package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"requiems-api/platform/httpx"
)

type dummyService struct{}

func TestGuard_NilService_Returns500(t *testing.T) {
	var svc *dummyService

	handler := httpx.Guard(svc, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "internal_error" {
		t.Errorf("expected error code %q, got %q", "internal_error", resp.Error)
	}
}

func TestGuard_NonNilService_CallsHandler(t *testing.T) {
	svc := &dummyService{}

	handler := httpx.Guard(svc, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
