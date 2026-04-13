package httpx_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"requiems-api/platform/httpx"
)

// bindReq is a sample struct used to exercise BindAndValidate.
type bindReq struct {
	Email string `json:"email" validate:"required,email"`
	Count int    `json:"count" validate:"min=1"`
}

func TestBindAndValidate_ValidJSON(t *testing.T) {
	body := strings.NewReader(`{"email":"user@example.com","count":3}`)
	r := httptest.NewRequest("POST", "/", body)

	var req bindReq
	if err := httpx.BindAndValidate(r, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Email != "user@example.com" {
		t.Errorf("email: want user@example.com, got %q", req.Email)
	}
	if req.Count != 3 {
		t.Errorf("count: want 3, got %d", req.Count)
	}
}

func TestBindAndValidate_MalformedJSON(t *testing.T) {
	body := strings.NewReader(`not-json`)
	r := httptest.NewRequest("POST", "/", body)

	var req bindReq
	if err := httpx.BindAndValidate(r, &req); err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestBindAndValidate_UnknownFields(t *testing.T) {
	body := strings.NewReader(`{"email":"user@example.com","count":1,"unknown":"field"}`)
	r := httptest.NewRequest("POST", "/", body)

	var req bindReq
	if err := httpx.BindAndValidate(r, &req); err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestBindAndValidate_ValidationFailure(t *testing.T) {
	// count is below min=1
	body := strings.NewReader(`{"email":"user@example.com","count":0}`)
	r := httptest.NewRequest("POST", "/", body)

	var req bindReq
	err := httpx.BindAndValidate(r, &req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	vf, ok := err.(*httpx.ValidationFailure)
	if !ok {
		t.Fatalf("expected *httpx.ValidationFailure, got %T", err)
	}
	if len(vf.Fields) == 0 {
		t.Error("expected at least one field error")
	}
}

// queryReq is used to exercise BindQuery with several field types.
type queryReq struct {
	Name    string    `query:"name"    validate:"required"`
	Age     int       `query:"age"`
	Score   float64   `query:"score"`
	Active  bool      `query:"active"`
	Since   time.Time `query:"since"`
	Ignored string    // no query tag – must be skipped
}

func TestBindQuery_StringField(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=alice", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Name != "alice" {
		t.Errorf("name: want alice, got %q", req.Name)
	}
}

func TestBindQuery_IntField(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&age=30", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Age != 30 {
		t.Errorf("age: want 30, got %d", req.Age)
	}
}

func TestBindQuery_FloatField(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&score=9.5", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Score != 9.5 {
		t.Errorf("score: want 9.5, got %v", req.Score)
	}
}

func TestBindQuery_BoolField(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&active=true", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !req.Active {
		t.Error("active: want true, got false")
	}
}

func TestBindQuery_TimeField(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&since=2024-06-15", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !req.Since.Equal(want) {
		t.Errorf("since: want %v, got %v", want, req.Since)
	}
}

func TestBindQuery_InvalidInt(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&age=notanint", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestBindQuery_InvalidFloat(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&score=notafloat", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err == nil {
		t.Fatal("expected error for invalid float, got nil")
	}
}

func TestBindQuery_InvalidBool(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&active=notabool", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err == nil {
		t.Fatal("expected error for invalid bool, got nil")
	}
}

func TestBindQuery_InvalidTime(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=x&since=not-a-date", nil)

	var req queryReq
	if err := httpx.BindQuery(r, &req); err == nil {
		t.Fatal("expected error for invalid date, got nil")
	}
}

func TestBindQuery_ValidationFailure(t *testing.T) {
	// name is required but missing
	r := httptest.NewRequest("GET", "/", nil)

	var req queryReq
	err := httpx.BindQuery(r, &req)
	if err == nil {
		t.Fatal("expected validation error for missing required field, got nil")
	}

	if _, ok := err.(*httpx.ValidationFailure); !ok {
		t.Fatalf("expected *httpx.ValidationFailure, got %T", err)
	}
}

func TestBindQuery_NonPointerDst(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	var req queryReq
	if err := httpx.BindQuery(r, req); err == nil {
		t.Fatal("expected error when dst is not a pointer, got nil")
	}
}
