package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFetchAndParse_LocalFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "swift.csv")

	csv := "swift_code,bank_name,city,country_name\nDEUTDEDB,Deutsche Bank,Frankfurt,Germany\nCHASUS33XXX,JPMorgan Chase,New York,United States\n"
	if err := os.WriteFile(filePath, []byte(csv), 0o600); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	records, err := fetchAndParse(filePath)
	if err != nil {
		t.Fatalf("fetchAndParse: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	if records[0].SwiftCode != "DEUTDEDBXXX" {
		t.Fatalf("expected DEUTDEDBXXX, got %q", records[0].SwiftCode)
	}

	if records[1].BranchCode != "XXX" {
		t.Fatalf("expected branch code XXX, got %q", records[1].BranchCode)
	}
}

func TestFetchAndParse_NoValidRecords(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "swift-empty.csv")

	csv := "swift_code,bank_name,city,country_name\nINVALID,No Bank,Nowhere,No Country\n"
	if err := os.WriteFile(filePath, []byte(csv), 0o600); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	_, err := fetchAndParse(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseRow_Accepts8And11CharCodes(t *testing.T) {
	t.Parallel()

	cols := colIndices{swift: 0, bankName: 1, city: 2, countryName: 3}

	r8, ok := parseRow([]string{"DEUTDEDB", "Deutsche Bank", "Frankfurt", "Germany"}, cols)
	if !ok {
		t.Fatal("expected 8-char SWIFT row to parse")
	}
	if r8.SwiftCode != "DEUTDEDBXXX" {
		t.Fatalf("expected expanded code DEUTDEDBXXX, got %q", r8.SwiftCode)
	}

	r11, ok := parseRow([]string{"CHASUS33XXX", "JPMorgan Chase", "New York", "United States"}, cols)
	if !ok {
		t.Fatal("expected 11-char SWIFT row to parse")
	}
	if r11.SwiftCode != "CHASUS33XXX" {
		t.Fatalf("expected CHASUS33XXX, got %q", r11.SwiftCode)
	}
}
