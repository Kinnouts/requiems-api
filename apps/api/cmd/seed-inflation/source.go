package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RawInflationRecord holds a single inflation entry before normalisation.
type RawInflationRecord struct {
	CountryCode string
	CountryName string
	Year        int
	Rate        float64
	Source      string
}

// worldBankEntry is one element in the World Bank API data array.
type worldBankEntry struct {
	Country struct {
		ID    string `json:"id"`    // 2-letter ISO code (e.g. "US")
		Value string `json:"value"` // full country name
	} `json:"country"`
	Date  string   `json:"date"`  // year as string, e.g. "2023"
	Value *float64 `json:"value"` // nullable — null when data is unavailable
}

// fetchAndParse downloads the World Bank inflation JSON and returns parsed records.
// The API returns a 2-element array: [metadata_object, data_array].
func fetchAndParse(url string) ([]RawInflationRecord, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download: HTTP %d", resp.StatusCode)
	}

	// Decode the outer 2-element array without loading it all into memory at once.
	var envelope []json.RawMessage
	if err = json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode envelope: %w", err)
	}
	if len(envelope) < 2 {
		return nil, fmt.Errorf("unexpected response: expected 2-element array, got %d", len(envelope))
	}

	var entries []worldBankEntry
	if err = json.Unmarshal(envelope[1], &entries); err != nil {
		return nil, fmt.Errorf("decode data array: %w", err)
	}

	records := make([]RawInflationRecord, 0, len(entries))
	for _, e := range entries {
		if e.Value == nil {
			continue // no data for this country/year
		}

		year, err := strconv.ParseInt(e.Date, 10, 16)
		if err != nil || year == 0 {
			continue // malformed date field
		}

		records = append(records, RawInflationRecord{
			CountryCode: e.Country.ID,
			CountryName: e.Country.Value,
			Year:        int(year),
			Rate:        *e.Value,
			Source:      "world_bank",
		})
	}

	return records, nil
}

// printStats prints a summary of the parsed inflation dataset to stdout.
func printStats(records []RawInflationRecord) {
	countries := make(map[string]int)
	minRate, maxRate := records[0].Rate, records[0].Rate

	for _, r := range records {
		countries[r.CountryCode]++
		if r.Rate < minRate {
			minRate = r.Rate
		}
		if r.Rate > maxRate {
			maxRate = r.Rate
		}
	}

	fmt.Printf("\n=== Inflation Seed Stats ===\n")
	fmt.Printf("Total records:    %d\n", len(records))
	fmt.Printf("Unique countries: %d\n", len(countries))
	fmt.Printf("Rate range:       %.4f%% to %.4f%%\n", minRate, maxRate)
	fmt.Printf("\nTop countries by record count:\n")
	printTopN(countries, 10)
}

func printTopN(m map[string]int, n int) {
	type kv struct {
		k string
		v int
	}
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	for i := 0; i < len(pairs) && i < n; i++ {
		maxIdx := i
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].v > pairs[maxIdx].v {
				maxIdx = j
			}
		}
		pairs[i], pairs[maxIdx] = pairs[maxIdx], pairs[i]
		fmt.Printf("  %-4s  %d\n", pairs[i].k, pairs[i].v)
	}
}
