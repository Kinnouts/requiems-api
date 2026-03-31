package asn

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
		r.Get("/ip/asn/{ip}", func(w http.ResponseWriter, r *http.Request) {
			httpx.Error(w, http.StatusServiceUnavailable, "service_unavailable", "ASN service not available")
		})
		r.Get("/ip/asn", func(w http.ResponseWriter, r *http.Request) {
			httpx.Error(w, http.StatusServiceUnavailable, "service_unavailable", "ASN service not available")
		})
	} else {
		RegisterRoutes(r, testSvc)
	}
	return r
}

func skipIfNoService(t *testing.T) {
	t.Helper()
	if testSvc == nil {
		t.Skip("ASN service not available (database not configured)")
	}
}

func TestASN_HappyPath(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/asn/8.8.8.8", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ASNResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP == "" {
		t.Error("expected non-empty IP in response")
	}
	if resp.Data.ASN == "" {
		t.Error("expected non-empty ASN for public IP")
	}
}

func TestASN_InvalidIPFormat(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/asn/not-an-ip", http.NoBody)
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

func TestASN_PrivateIP(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/asn/192.168.1.1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for private IP, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ASNResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.IP != "192.168.1.1" {
		t.Errorf("expected IP 192.168.1.1, got %s", resp.Data.IP)
	}
	if resp.Data.ASN != "" {
		t.Errorf("expected empty ASN for private IP, got %s", resp.Data.ASN)
	}
}

func TestASN_NoIPParam_UsesRemoteAddr(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/asn", http.NoBody)
	req.RemoteAddr = "8.8.8.8:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestASN_XForwardedFor(t *testing.T) {
	skipIfNoService(t)
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ip/asn", http.NoBody)
	req.Header.Set("X-Forwarded-For", "1.1.1.1, 10.0.0.1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ASNResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Data.IP != "1.1.1.1" {
		t.Errorf("expected IP 1.1.1.1 (first XFF), got %s", resp.Data.IP)
	}
}
