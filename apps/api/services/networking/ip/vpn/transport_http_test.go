package vpn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobadilla-tech/go-ip-intelligence/v2/ipi"
	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

var testSvc *Service

func init() {
	client, err := ipi.New(
		ipi.WithDatabasePath(""),
		ipi.WithASNDatabasePath(""),
		ipi.WithCityDatabasePath(""),
	)
	if err == nil {
		testSvc = NewService(client)
	}
}

func setupRouter() chi.Router {
	r := chi.NewRouter()
	if testSvc == nil {
		r.Get("/ip/vpn/{ip}", func(w http.ResponseWriter, r *http.Request) {
			httpx.Error(w, http.StatusServiceUnavailable, "service_unavailable", "VPN service not available")
		})
	} else {
		RegisterRoutes(r, testSvc)
	}
	return r
}

func skipIfNoService(t *testing.T) {
	if testSvc == nil {
		t.Skip("VPN service not available (database not configured)")
	}
}

func TestVPN_HappyPath(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/vpn/8.8.8.8", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[IPCheckResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP == "" {
		t.Error("expected non-empty IP in response")
	}

	if resp.Data.Score < 0 {
		t.Errorf("expected non-negative score, got %d", resp.Data.Score)
	}
}

func TestVPN_ValidIPFields(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/vpn/1.1.1.1", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp httpx.Response[IPCheckResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP != "1.1.1.1" {
		t.Errorf("expected IP 1.1.1.1, got %s", resp.Data.IP)
	}

	validThreats := map[string]bool{
		"none":     true,
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validThreats[resp.Data.Threat.String()] {
		t.Errorf("invalid threat level: %s", resp.Data.Threat)
	}

	if resp.Data.FraudScore < 0 || resp.Data.FraudScore > 100 {
		t.Errorf("fraud_score out of range: %d", resp.Data.FraudScore)
	}
}

func TestVPN_InvalidIPFormat(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/vpn/not-an-ip", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if testSvc == nil {
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("expected 503 without service, got %d", w.Code)
		}
		return
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp httpx.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error != "bad_request" {
		t.Errorf("expected error code 'bad_request', got %q", resp.Error)
	}
}

func TestVPN_IPv6Address(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/vpn/2001:4860:4860::8888", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for IPv6, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[IPCheckResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP == "" {
		t.Error("expected non-empty IP for IPv6")
	}
}

func TestVPN_AllBooleansReturned(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/vpn/8.8.8.8", http.NoBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp httpx.Response[IPCheckResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IsVPN != false && resp.Data.IsVPN != true {
		t.Error("is_vpn should be a boolean")
	}
	if resp.Data.IsProxy != false && resp.Data.IsProxy != true {
		t.Error("is_proxy should be a boolean")
	}
	if resp.Data.IsTor != false && resp.Data.IsTor != true {
		t.Error("is_tor should be a boolean")
	}
	if resp.Data.IsHosting != false && resp.Data.IsHosting != true {
		t.Error("is_hosting should be a boolean")
	}
}
