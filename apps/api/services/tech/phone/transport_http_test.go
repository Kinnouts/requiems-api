package phone

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

func TestPhone_ValidUSNumber(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/validate/phone?number=%2B12015551234", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ValidateResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Data.Valid {
		t.Error("expected valid=true for a valid US number")
	}
	if resp.Data.Country != "US" {
		t.Errorf("expected country US, got %q", resp.Data.Country)
	}
	if resp.Data.Formatted == "" {
		t.Error("expected non-empty formatted number")
	}
	if resp.Data.Type == "" {
		t.Error("expected non-empty type")
	}
}

func TestPhone_InvalidNumber(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/validate/phone?number=12345", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ValidateResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Valid {
		t.Error("expected valid=false for an invalid number")
	}
	if resp.Data.Country != "" {
		t.Errorf("expected empty country for invalid number, got %q", resp.Data.Country)
	}
}

func TestPhone_MissingNumber(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/validate/phone", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPhone_UKMobile(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/validate/phone?number=%2B447400123456", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ValidateResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Data.Valid {
		t.Error("expected valid=true for a valid UK number")
	}
	if resp.Data.Country != "GB" {
		t.Errorf("expected country GB, got %q", resp.Data.Country)
	}
	if resp.Data.Type != "mobile" {
		t.Errorf("expected type mobile, got %q", resp.Data.Type)
	}
}

func TestPhone_CarrierPresent(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/validate/phone?number=%2B51923531893", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ValidateResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Data.Valid {
		t.Fatal("expected valid=true")
	}
	if resp.Data.Carrier == nil {
		t.Fatal("expected carrier to be present")
	}
	if resp.Data.Carrier.Name == "" {
		t.Error("expected non-empty carrier name")
	}
	if resp.Data.Carrier.Source != "metadata" {
		t.Errorf("expected carrier source %q, got %q", "metadata", resp.Data.Carrier.Source)
	}
}

func TestPhone_RiskVoIP(t *testing.T) {
	svc := NewService()
	// Google Voice numbers are VOIP type in the US (area code 202 VOIP range)
	// Use a number whose type we know via the service
	tests := []struct {
		name   string
		number string
	}{
		// +1-500 numbers are personal/VOIP in the US
		{"US personal number", "+15005550006"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Validate(tt.number)
			if result.Risk == nil {
				t.Fatal("expected risk to be present on valid number")
			}
		})
	}
}

func TestPhone_RiskMobile(t *testing.T) {
	svc := NewService()
	result := svc.Validate("+447400123456")

	if !result.Valid {
		t.Fatal("expected valid=true")
	}
	if result.Risk == nil {
		t.Fatal("expected risk to be present")
	}
	if result.Risk.IsVoIP {
		t.Error("expected is_voip=false for mobile number")
	}
	if result.Risk.IsVirtual {
		t.Error("expected is_virtual=false for mobile number")
	}
}

func TestPhone_InvalidHasNoCarrierOrRisk(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/validate/phone?number=12345", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp httpx.Response[ValidateResponse]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Valid {
		t.Fatal("expected valid=false")
	}
	if resp.Data.Carrier != nil {
		t.Error("expected carrier to be absent for invalid number")
	}
	if resp.Data.Risk != nil {
		t.Error("expected risk to be absent for invalid number")
	}
}

func TestService_NumberType(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		wantType string
	}{
		{"UK landline", "+441613281234", "landline"},
		{"UK mobile", "+447400123456", "mobile"},
		{"US toll-free", "+18005551234", "toll_free"},
	}

	svc := NewService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Validate(tt.number)
			if !result.Valid {
				t.Fatalf("expected valid=true for %q", tt.number)
			}
			if result.Type != tt.wantType {
				t.Errorf("expected type %q, got %q", tt.wantType, result.Type)
			}
		})
	}
}
