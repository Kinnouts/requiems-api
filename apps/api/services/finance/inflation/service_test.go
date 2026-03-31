package inflation

import (
	"testing"

	"requiems-api/platform/httpx"
)

// ---- Request validation ----

func TestRequest_Valid_US(t *testing.T) {
	req := Request{Country: "US"}
	if err := httpx.Validate.Struct(&req); err != nil {
		t.Errorf("expected no error for valid country US, got %v", err)
	}
}

func TestRequest_Valid_GB(t *testing.T) {
	req := Request{Country: "GB"}
	if err := httpx.Validate.Struct(&req); err != nil {
		t.Errorf("expected no error for valid country GB, got %v", err)
	}
}

func TestRequest_Empty_Country_Fails(t *testing.T) {
	req := Request{}
	if err := httpx.Validate.Struct(&req); err == nil {
		t.Error("expected error for empty country, got nil")
	}
}

func TestRequest_Invalid_Country_ZZZ_Fails(t *testing.T) {
	req := Request{Country: "ZZZ"}
	if err := httpx.Validate.Struct(&req); err == nil {
		t.Error("expected error for invalid country code ZZZ, got nil")
	}
}

func TestRequest_Lowercase_Fails(t *testing.T) {
	// iso3166_1_alpha2 requires uppercase; the transport layer uppercases before binding.
	req := Request{Country: "us"}
	if err := httpx.Validate.Struct(&req); err == nil {
		t.Error("expected error for lowercase country code, got nil")
	}
}
