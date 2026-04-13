package main

import (
	"encoding/json"
	"testing"
)

func TestParseFREDRow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		record  []string
		wantVal float64
		wantOK  bool
	}{
		{name: "valid daily row", record: []string{"2020-01-15", "55.42"}, wantVal: 55.42, wantOK: true},
		{name: "valid annual row", record: []string{"2010-01-01", "80.00"}, wantVal: 80.00, wantOK: true},
		{name: "dot placeholder", record: []string{"2020-06-01", "."}, wantOK: false},
		{name: "empty value", record: []string{"2020-06-01", ""}, wantOK: false},
		{name: "negative value", record: []string{"2020-06-01", "-5.0"}, wantOK: false},
		{name: "non-numeric value", record: []string{"2020-06-01", "N/A"}, wantOK: false},
		{name: "too few fields", record: []string{"2020-06-01"}, wantOK: false},
		// Years outside [1960, now] are rejected.
		{name: "year too early", record: []string{"1959-01-01", "10.0"}, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			val, _, ok := parseFREDRow(tt.record)
			if ok != tt.wantOK {
				t.Fatalf("parseFREDRow ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && val != tt.wantVal {
				t.Fatalf("parseFREDRow val = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

func TestParseYahooClose(t *testing.T) {
	t.Parallel()

	closes := []interface{}{
		float64(1234.56),
		nil,
		json.Number("99.99"),
		float64(0),    // zero — rejected
		float64(-1.0), // negative — rejected
	}

	tests := []struct {
		name    string
		idx     int
		wantVal float64
		wantOK  bool
	}{
		{name: "float64 value", idx: 0, wantVal: 1234.56, wantOK: true},
		{name: "nil entry", idx: 1, wantOK: false},
		{name: "json.Number", idx: 2, wantVal: 99.99, wantOK: true},
		{name: "zero", idx: 3, wantOK: false},
		{name: "negative", idx: 4, wantOK: false},
		{name: "out of bounds", idx: 99, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			val, ok := parseYahooClose(closes, tt.idx)
			if ok != tt.wantOK {
				t.Fatalf("parseYahooClose ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && val != tt.wantVal {
				t.Fatalf("parseYahooClose val = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

func TestParseYear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{input: "2020-01-15", want: 2020},
		{input: "1985", want: 1985},
		{input: "198", wantErr: true}, // too short
		{input: "abcd", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := parseYear(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseYear(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("parseYear(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildRecords(t *testing.T) {
	t.Parallel()

	cfg := CommodityConfig{
		Slug:       "oil",
		Name:       "Crude Oil",
		Unit:       "barrel",
		Currency:   "USD",
		ConvFactor: 1.0,
	}

	byYear := map[int]*yearAcc{
		2020: {sum: 40.0, count: 2}, // avg = 20.0
	}

	records := buildRecords(cfg, byYear)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]
	if r.Slug != "oil" {
		t.Errorf("Slug = %q, want %q", r.Slug, "oil")
	}
	if r.Year != 2020 {
		t.Errorf("Year = %d, want 2020", r.Year)
	}
	if r.Price != 20.0 {
		t.Errorf("Price = %v, want 20.0", r.Price)
	}
}

func TestAccumulate(t *testing.T) {
	t.Parallel()

	m := make(map[int]*yearAcc)
	accumulate(m, 2021, 10.0)
	accumulate(m, 2021, 20.0)
	accumulate(m, 2022, 5.0)

	if m[2021].count != 2 || m[2021].sum != 30.0 {
		t.Fatalf("year 2021: count=%d sum=%v, want count=2 sum=30.0", m[2021].count, m[2021].sum)
	}
	if m[2022].count != 1 || m[2022].sum != 5.0 {
		t.Fatalf("year 2022: count=%d sum=%v, want count=1 sum=5.0", m[2022].count, m[2022].sum)
	}
}
