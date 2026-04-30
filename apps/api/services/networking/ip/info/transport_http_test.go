package info

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
		r.Get("/ip/{ip}", func(w http.ResponseWriter, r *http.Request) {
			httpx.Error(w, http.StatusServiceUnavailable, "service_unavailable", "IP info service not available")
		})
		r.Get("/ip", func(w http.ResponseWriter, r *http.Request) {
			httpx.Error(w, http.StatusServiceUnavailable, "service_unavailable", "IP info service not available")
		})
	} else {
		RegisterRoutes(r, testSvc)
	}
	return r
}

func skipIfNoService(t *testing.T) {
	t.Helper()
	if testSvc == nil {
		t.Skip("IP info service not available (database not configured)")
	}
}

func TestInfo_HappyPath_PathParam(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/8.8.8.8", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[LookupResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP != "8.8.8.8" {
		t.Errorf("expected IP 8.8.8.8, got %s", resp.Data.IP)
	}
}

func TestInfo_HappyPath_NoPathParam(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip", http.NoBody)
	req.RemoteAddr = "1.1.1.1:54321"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[LookupResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP == "" {
		t.Error("expected non-empty IP in response")
	}
}

func TestInfo_InvalidIPFormat(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/not-an-ip", http.NoBody)
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

func TestInfo_PrivateIP_ReturnsEmpty(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/10.0.0.1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for private IP, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[LookupResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP == "" {
		t.Error("expected IP field to be set even for private addresses")
	}
	// Country/city should be empty for private IPs
	if resp.Data.Country != "" {
		t.Errorf("expected empty country for private IP, got %s", resp.Data.Country)
	}
}

func TestInfo_XRealIP(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip", http.NoBody)
	req.Header.Set("X-Real-IP", "8.8.4.4")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[LookupResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Data.IP != "8.8.4.4" {
		t.Errorf("expected IP 8.8.4.4 from X-Real-IP, got %s", resp.Data.IP)
	}
}

func TestInfo_IPv6Address(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/2001:4860:4860::8888", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for IPv6, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[LookupResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP == "" {
		t.Error("expected non-empty IP for IPv6")
	}
}

func TestInfo_ResponseFields(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/1.1.1.1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.Bytes()

	var resp httpx.Response[LookupResponse]
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP != "1.1.1.1" {
		t.Errorf("expected IP 1.1.1.1, got %s", resp.Data.IP)
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("failed to decode raw response: %v", err)
	}
	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected object at data")
	}
	if _, ok := data["is_vpn"]; !ok {
		t.Error("expected data.is_vpn to be present in response JSON")
	}
}
