package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RawSWIFTRecord holds one parsed entry from the SWIFT/BIC dataset.
type RawSWIFTRecord struct {
	SwiftCode    string // full 11-char BIC, always uppercase
	BankCode     string // chars 1-4
	CountryCode  string // chars 5-6
	LocationCode string // chars 7-8
	BranchCode   string // chars 9-11; "XXX" = primary office
	BankName     string
	City         string
	CountryName  string
}

// fetchAndParse downloads the SWIFT/BIC CSV from url and returns parsed records.
//
// The CSV header row is read to build a column-name-to-index map, making the
// parser robust against column reordering between dataset versions. Expected
// column names (case-insensitive): swift_code (or bic), bank_name, city,
// country_name (or country). The BIC components are derived from the swift_code.
//
// 8-character codes are expanded to 11 characters by appending "XXX".
// Rows with malformed or missing SWIFT codes are skipped.
func fetchAndParse(url string) ([]RawSWIFTRecord, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download: HTTP %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	reader.TrimLeadingSpace = true

	// Read header row and build column index map.
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	colIdx := make(map[string]int, len(header))
	for i, h := range header {
		colIdx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Resolve SWIFT code column — try common names.
	swiftCol, ok := resolveCol(colIdx, "swift_code", "bic", "swiftcode", "swift")
	if !ok {
		return nil, fmt.Errorf("could not find SWIFT code column in header: %v", header)
	}

	bankNameCol, _ := resolveCol(colIdx, "bank_name", "name", "institution_name", "institution")
	cityCol, _ := resolveCol(colIdx, "city", "branch_city", "branch_address")
	countryNameCol, _ := resolveCol(colIdx, "country_name", "country")

	var records []RawSWIFTRecord

	for {
		row, err := reader.Read()
		if err != nil {
			break // EOF or error — stop reading
		}

		if swiftCol >= len(row) {
			continue
		}

		raw := strings.ToUpper(strings.TrimSpace(row[swiftCol]))
		if len(raw) != 8 && len(raw) != 11 {
			continue
		}
		if !isValidBIC(raw) {
			continue
		}

		if len(raw) == 8 {
			raw += "XXX"
		}

		rec := RawSWIFTRecord{
			SwiftCode:    raw,
			BankCode:     raw[0:4],
			CountryCode:  raw[4:6],
			LocationCode: raw[6:8],
			BranchCode:   raw[8:11],
		}

		if bankNameCol >= 0 && bankNameCol < len(row) {
			rec.BankName = strings.TrimSpace(row[bankNameCol])
		}
		if cityCol >= 0 && cityCol < len(row) {
			rec.City = strings.TrimSpace(row[cityCol])
		}
		if countryNameCol >= 0 && countryNameCol < len(row) {
			rec.CountryName = strings.TrimSpace(row[countryNameCol])
		}

		records = append(records, rec)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no SWIFT code records found — the dataset format may have changed")
	}

	return records, nil
}

// resolveCol returns the index of the first matching column name (case-insensitive)
// from the given candidates, or -1 and false if none match.
func resolveCol(colIdx map[string]int, candidates ...string) (int, bool) {
	for _, name := range candidates {
		if idx, ok := colIdx[name]; ok {
			return idx, true
		}
	}
	return -1, false
}

// isValidBIC performs a lightweight character-class check on the uppercased
// BIC string (8 or 11 chars). It does not validate country codes or check
// whether the BIC is registered — it only filters out clearly malformed rows.
func isValidBIC(bic string) bool {
	if len(bic) != 8 && len(bic) != 11 {
		return false
	}
	// Positions 0-3: letters only (bank code)
	for i := range 4 {
		if bic[i] < 'A' || bic[i] > 'Z' {
			return false
		}
	}
	// Positions 4-5: letters only (country code)
	for i := 4; i < 6; i++ {
		if bic[i] < 'A' || bic[i] > 'Z' {
			return false
		}
	}
	// Positions 6-7: alphanumeric (location code)
	for i := 6; i < 8; i++ {
		b := bic[i]
		if !((b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')) {
			return false
		}
	}
	// Positions 8-10: alphanumeric (branch code, only for 11-char)
	if len(bic) == 11 {
		for i := 8; i < 11; i++ {
			b := bic[i]
			if !((b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')) {
				return false
			}
		}
	}
	return true
}

// printStats prints a summary of the parsed dataset to stdout.
func printStats(records []RawSWIFTRecord) {
	primary := 0
	countries := make(map[string]struct{})
	for _, r := range records {
		if r.BranchCode == "XXX" {
			primary++
		}
		countries[r.CountryCode] = struct{}{}
	}

	fmt.Printf("\n=== SWIFT Seed Stats ===\n")
	fmt.Printf("Total codes:      %d\n", len(records))
	fmt.Printf("Primary offices:  %d\n", primary)
	fmt.Printf("Branch offices:   %d\n", len(records)-primary)
	fmt.Printf("Countries:        %d\n", len(countries))

	fmt.Printf("\nSample entries:\n")
	for i, r := range records {
		if i >= 8 {
			break
		}
		fmt.Printf("  %-11s  %-4s  %-2s  %q\n",
			r.SwiftCode, r.BankCode, r.CountryCode, r.BankName)
	}
}
