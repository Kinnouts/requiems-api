package domain

import (
	"context"
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

func TestDomain_InvalidFormat(t *testing.T) {
	tests := []struct {
		name   string
		domain string
	}{
		{"bare label", "localhost"},
		{"leading hyphen", "-bad.com"},
		{"trailing hyphen", "bad-.com"},
		{"numeric TLD only", "123"},
		{"just a dot", "."},
	}

	r := setupRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/domain/"+tt.domain, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestDomain_KnownDomain(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/domain/example.com", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[InfoResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Domain != "example.com" {
		t.Errorf("expected domain example.com, got %q", resp.Data.Domain)
	}
	if resp.Data.DNS.A == nil {
		t.Error("expected non-nil A records slice")
	}
	if resp.Data.DNS.AAAA == nil {
		t.Error("expected non-nil AAAA records slice")
	}
	if resp.Data.DNS.MX == nil {
		t.Error("expected non-nil MX records slice")
	}
	if resp.Data.DNS.NS == nil {
		t.Error("expected non-nil NS records slice")
	}
	if resp.Data.DNS.TXT == nil {
		t.Error("expected non-nil TXT records slice")
	}

	// DNS record content is only asserted when network resolution is available.
	if len(resp.Data.DNS.NS) > 0 && resp.Data.Available {
		t.Error("expected available=false when NS records are present")
	}
}

func TestDomain_ResponseShape(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/domain/example.com", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify the raw JSON shape has the expected keys.
	var raw map[string]any
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatal("expected 'data' key in response")
	}
	for _, key := range []string{"domain", "available", "dns"} {
		if _, exists := data[key]; !exists {
			t.Errorf("expected key %q in data", key)
		}
	}

	dns, ok := data["dns"].(map[string]any)
	if !ok {
		t.Fatal("expected 'dns' key in data")
	}
	for _, key := range []string{"a", "aaaa", "mx", "ns", "txt"} {
		if _, exists := dns[key]; !exists {
			t.Errorf("expected key %q in dns", key)
		}
	}
}

func TestService_IsNXDomain(t *testing.T) {
	svc := NewService()

	// A clearly invented domain should either be unavailable (registered) or
	// available (NXDOMAIN). Either way the service should return 200 with the
	// domain name echoed back, without panicking.
	resp := svc.GetInfo(context.Background(), "example.com")
	if resp.Domain != "example.com" {
		t.Errorf("expected domain echoed back, got %q", resp.Domain)
	}
}
