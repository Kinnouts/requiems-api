package mx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}

func TestMXLookup_InvalidDomain(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name   string
		domain string
	}{
		{"empty-like path", "not_a_domain"},
		{"plain label", "localhost"},
		{"starts with dash", "-bad.com"},
		{"IP address", "1.2.3.4"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/mx/"+tc.domain, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400 for %q, got %d", tc.domain, w.Code)
			}

			var resp httpx.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}
			if resp.Error != "bad_request" {
				t.Errorf("expected error code 'bad_request', got %q", resp.Error)
			}
		})
	}
}

func TestMXLookup_NonExistentDomain(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/mx/nonexistent-domain-that-does-not-exist.invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return 404 (no MX records / NXDOMAIN) or 500 (network unavailable in CI)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 404 or 500 for non-existent domain, got %d", w.Code)
	}
}

func TestMXLookup_HappyPath(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/mx/gmail.com", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Network may not be available in all CI environments; accept 200 or 500/404
	if w.Code == http.StatusInternalServerError || w.Code == http.StatusNotFound {
		t.Skip("DNS not available in this environment")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[LookupResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Domain != "gmail.com" {
		t.Errorf("expected domain 'gmail.com', got %q", resp.Data.Domain)
	}
	if len(resp.Data.Records) == 0 {
		t.Error("expected at least one MX record for gmail.com")
	}

	// Verify priority ordering (ascending)
	for i := 1; i < len(resp.Data.Records); i++ {
		if resp.Data.Records[i].Priority < resp.Data.Records[i-1].Priority {
			t.Errorf("records not sorted by priority: record[%d].Priority=%d < record[%d].Priority=%d",
				i, resp.Data.Records[i].Priority, i-1, resp.Data.Records[i-1].Priority)
		}
	}

	// Each record should have a non-empty host
	for i, rec := range resp.Data.Records {
		if rec.Host == "" {
			t.Errorf("record[%d] has empty host", i)
		}
	}
}
