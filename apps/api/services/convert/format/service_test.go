package convformat

import (
	"strings"
	"testing"
)

func TestService_Convert_SameFormat(t *testing.T) {
	svc := NewService()
	req := Request{From: "json", To: "json", Content: `{"name":"Alice"}`}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Result != req.Content {
		t.Errorf("expected passthrough, got %q", resp.Result)
	}
}

func TestService_Convert_ContentTooLarge(t *testing.T) {
	svc := NewService()
	big := strings.Repeat("a", maxContentSize+1)
	req := Request{From: "json", To: "yaml", Content: big}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error for oversized content")
	}
}

// --- JSON ↔ YAML ---

func TestService_JSONToYAML(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "json",
		To:      "yaml",
		Content: `{"name":"Alice","age":30}`,
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "name: Alice") {
		t.Errorf("expected YAML with 'name: Alice', got %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "age: 30") {
		t.Errorf("expected YAML with 'age: 30', got %q", resp.Result)
	}
}

func TestService_YAMLToJSON(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "yaml",
		To:      "json",
		Content: "name: Alice\nage: 30\n",
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, `"name"`) {
		t.Errorf("expected JSON with 'name' key, got %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected JSON with 'Alice', got %q", resp.Result)
	}
}

func TestService_InvalidJSON(t *testing.T) {
	svc := NewService()
	req := Request{From: "json", To: "yaml", Content: `{invalid`}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestService_InvalidYAML(t *testing.T) {
	svc := NewService()
	req := Request{From: "yaml", To: "json", Content: ":\t:bad yaml\n"}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

// --- JSON ↔ CSV ---

func TestService_JSONToCSV(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "json",
		To:      "csv",
		Content: `[{"name":"Alice","age":"30"},{"name":"Bob","age":"25"}]`,
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "name") {
		t.Errorf("expected CSV with 'name' header, got %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected CSV with 'Alice', got %q", resp.Result)
	}
}

func TestService_CSVToJSON(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "csv",
		To:      "json",
		Content: "name,age\nAlice,30\nBob,25\n",
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, `"name"`) {
		t.Errorf("expected JSON with 'name' key, got %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected JSON with 'Alice', got %q", resp.Result)
	}
}

func TestService_JSONToCSV_NonArray(t *testing.T) {
	svc := NewService()
	req := Request{From: "json", To: "csv", Content: `{"name":"Alice"}`}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error when converting non-array JSON to CSV")
	}
}

// --- JSON ↔ XML ---

func TestService_JSONToXML(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "json",
		To:      "xml",
		Content: `{"name":"Alice","age":30}`,
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "<root>") {
		t.Errorf("expected XML with <root>, got %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected XML with 'Alice', got %q", resp.Result)
	}
}

func TestService_XMLToJSON(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "xml",
		To:      "json",
		Content: `<root><name>Alice</name><age>30</age></root>`,
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected JSON with 'Alice', got %q", resp.Result)
	}
}

func TestService_InvalidXML(t *testing.T) {
	svc := NewService()
	req := Request{From: "xml", To: "json", Content: "<unclosed>"}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error for invalid XML")
	}
}

// --- JSON ↔ TOML ---

func TestService_JSONToTOML(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "json",
		To:      "toml",
		Content: `{"name":"Alice","age":30}`,
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected TOML with 'Alice', got %q", resp.Result)
	}
}

func TestService_TOMLToJSON(t *testing.T) {
	svc := NewService()
	req := Request{
		From:    "toml",
		To:      "json",
		Content: "name = \"Alice\"\nage = 30\n",
	}
	resp, err := svc.Convert(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, `"name"`) {
		t.Errorf("expected JSON with 'name' key, got %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Alice") {
		t.Errorf("expected JSON with 'Alice', got %q", resp.Result)
	}
}

func TestService_InvalidTOML(t *testing.T) {
	svc := NewService()
	req := Request{From: "toml", To: "json", Content: "= invalid toml"}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error for invalid TOML")
	}
}

func TestService_JSONToTOML_Array(t *testing.T) {
	svc := NewService()
	req := Request{From: "json", To: "toml", Content: `[1,2,3]`}
	_, err := svc.Convert(req)
	if err == nil {
		t.Fatal("expected error when converting JSON array to TOML")
	}
}
