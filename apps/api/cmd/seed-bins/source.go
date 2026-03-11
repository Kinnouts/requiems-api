package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RawBINRecord holds a single BIN entry before normalisation.
type RawBINRecord struct {
	BINPrefix   string
	Scheme      string
	CardType    string
	CardLevel   string
	IssuerName  string
	IssuerURL   string
	IssuerPhone string
	CountryCode string
	CountryName string
	Prepaid     bool
	Source      string
	Confidence  float64
}

// Source describes a BIN dataset source and its parser.
type Source struct {
	Name       string
	URL        string
	Confidence float64
	Parse      func(r io.Reader, sourceName string, baseConfidence float64) ([]RawBINRecord, error)
}

// fetchAndParse downloads a source CSV and returns parsed records.
func fetchAndParse(src Source) ([]RawBINRecord, error) {
	client := &http.Client{Timeout: 120 * time.Second}

	resp, err := client.Get(src.URL)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download: HTTP %d", resp.StatusCode)
	}

	return src.Parse(resp.Body, src.Name, src.Confidence)
}

// parseIannuttall parses the iannuttall/binlist-data CSV.
//
// Columns (header row):
//
//	bin,brand,type,category,issuer,alpha_2,alpha_3,country,latitude,longitude,bank_phone,bank_url
func parseIannuttall(r io.Reader, sourceName string, baseConf float64) ([]RawBINRecord, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// Skip header.
	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	var records []RawBINRecord
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip malformed rows
		}
		if len(row) < 12 {
			continue
		}

		bin := strings.TrimSpace(row[0])
		if !isValidBINPrefix(bin) {
			continue
		}

		rec := RawBINRecord{
			BINPrefix:   bin,
			Scheme:      strings.TrimSpace(row[1]),
			CardType:    strings.TrimSpace(row[2]),
			CardLevel:   strings.TrimSpace(row[3]),
			IssuerName:  strings.TrimSpace(row[4]),
			CountryCode: strings.TrimSpace(row[5]),
			CountryName: strings.TrimSpace(row[7]),
			IssuerPhone: strings.TrimSpace(row[10]),
			IssuerURL:   strings.TrimSpace(row[11]),
			Source:      sourceName,
			Confidence:  baseConf,
		}
		records = append(records, rec)
	}

	return records, nil
}

// parseVenelinkochev parses the venelinkochev/bin-list-data CSV.
//
// Columns (header row):
//
//	BIN,Brand,Type,Category,Issuer,IssuerPhone,IssuerUrl,isoCode2,isoCode3,CountryName
func parseVenelinkochev(r io.Reader, sourceName string, baseConf float64) ([]RawBINRecord, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	cr.LazyQuotes = true

	// Skip header.
	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	var records []RawBINRecord
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip malformed rows
		}
		if len(row) < 10 {
			continue
		}

		bin := strings.TrimSpace(row[0])
		if !isValidBINPrefix(bin) {
			continue
		}

		rec := RawBINRecord{
			BINPrefix:   bin,
			Scheme:      strings.TrimSpace(row[1]),
			CardType:    strings.TrimSpace(row[2]),
			CardLevel:   strings.TrimSpace(row[3]),
			IssuerName:  strings.TrimSpace(row[4]),
			IssuerPhone: strings.TrimSpace(row[5]),
			IssuerURL:   strings.TrimSpace(row[6]),
			CountryCode: strings.TrimSpace(row[7]),
			CountryName: strings.TrimSpace(row[9]),
			Source:      sourceName,
			Confidence:  baseConf,
		}
		records = append(records, rec)
	}

	return records, nil
}

// isValidBINPrefix returns true if s is a 6- or 8-digit numeric string.
func isValidBINPrefix(s string) bool {
	if len(s) != 6 && len(s) != 8 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// printStats prints a summary of the merged BIN dataset to stdout.
func printStats(records map[string]RawBINRecord) {
	schemes := make(map[string]int)
	types := make(map[string]int)
	countries := make(map[string]int)
	withBank := 0

	for _, r := range records {
		schemes[r.Scheme]++
		types[r.CardType]++
		countries[r.CountryCode]++
		if r.IssuerName != "" {
			withBank++
		}
	}

	fmt.Printf("\n=== BIN Seed Stats ===\n")
	fmt.Printf("Total unique BINs: %d\n", len(records))
	fmt.Printf("With bank name:    %d (%.1f%%)\n", withBank, 100*float64(withBank)/float64(len(records)))
	fmt.Printf("\nTop schemes:\n")
	printTopN(schemes, 10)
	fmt.Printf("\nCard types:\n")
	printTopN(types, 10)
	fmt.Printf("\nTop countries:\n")
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
	// simple selection sort for top N
	for i := 0; i < len(pairs) && i < n; i++ {
		max := i
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].v > pairs[max].v {
				max = j
			}
		}
		pairs[i], pairs[max] = pairs[max], pairs[i]
		fmt.Printf("  %-20s %d\n", pairs[i].k, pairs[i].v)
	}
	_ = strconv.Itoa(0) // keep import used
}
