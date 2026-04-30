package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
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

// colIndices holds the resolved column positions from the CSV header.
type colIndices struct {
	swift       int
	bankName    int
	city        int
	countryName int
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
func fetchAndParse(source string) ([]RawSWIFTRecord, error) {
	return fetchAndParseWithCountries(source, "")
}

func fetchAndParseWithCountries(source, countriesSource string) ([]RawSWIFTRecord, error) {
	body, err := openSource(source)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	reader := csv.NewReader(body)
	reader.TrimLeadingSpace = true

	firstRow, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read first row: %w", err)
	}

	if isHeaderRow(firstRow) {
		return parseHeaderBased(reader, firstRow)
	}

	countryNames := map[string]string{}
	if countriesSource != "" {
		countryNames, err = loadCountryNames(countriesSource)
		if err != nil {
			return nil, err
		}
	}

	return parseNoHeaderBased(reader, firstRow, countryNames)
}

func parseHeaderBased(reader *csv.Reader, header []string) ([]RawSWIFTRecord, error) {
	cols, err := parseHeaderCols(header)
	if err != nil {
		return nil, err
	}

	var records []RawSWIFTRecord
	for {
		row, err := reader.Read()
		if err != nil {
			break // EOF or read error — stop
		}
		rec, ok := parseRow(row, cols)
		if ok {
			records = append(records, rec)
		}
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no SWIFT code records found — the dataset format may have changed")
	}

	return records, nil
}

func parseNoHeaderBased(reader *csv.Reader, firstRow []string, countryNames map[string]string) ([]RawSWIFTRecord, error) {
	var records []RawSWIFTRecord

	if rec, ok := parseNoHeaderRow(firstRow, countryNames); ok {
		records = append(records, rec)
	}

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if rec, ok := parseNoHeaderRow(row, countryNames); ok {
			records = append(records, rec)
		}
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no SWIFT code records found — the dataset format may have changed")
	}

	return records, nil
}

func openSource(source string) (io.ReadCloser, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 30 * time.Second}

		resp, err := client.Get(source) //nolint:gosec // source comes from trusted CLI flags, not user HTTP input
		if err != nil {
			return nil, fmt.Errorf("download: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("download: HTTP %d", resp.StatusCode)
		}

		return resp.Body, nil
	}

	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", source, err)
	}

	return f, nil
}

// parseHeaderCols resolves column indexes from a CSV header row.
func parseHeaderCols(header []string) (colIndices, error) {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	swiftCol, ok := resolveCol(idx, "swift_code", "bic", "swiftcode", "swift")
	if !ok {
		return colIndices{}, fmt.Errorf("could not find SWIFT code column in header: %v", header)
	}

	bankNameCol, _ := resolveCol(idx, "bank_name", "name", "institution_name", "institution")
	cityCol, _ := resolveCol(idx, "city", "branch_city", "branch_address")
	countryNameCol, _ := resolveCol(idx, "country_name", "country")

	return colIndices{
		swift:       swiftCol,
		bankName:    bankNameCol,
		city:        cityCol,
		countryName: countryNameCol,
	}, nil
}

func isHeaderRow(row []string) bool {
	for _, cell := range row {
		normalized := strings.ToLower(strings.TrimSpace(cell))
		switch normalized {
		case "swift_code", "bic", "swiftcode", "swift", "bank_name", "name", "city", "country_name", "country":
			return true
		}
	}
	return false
}

func parseNoHeaderRow(row []string, countryNames map[string]string) (RawSWIFTRecord, bool) {
	if len(row) < 4 {
		return RawSWIFTRecord{}, false
	}

	countryCode := strings.ToUpper(strings.TrimSpace(row[0]))
	raw := strings.ToUpper(strings.TrimSpace(row[1]))
	if !isValidBIC(raw) {
		return RawSWIFTRecord{}, false
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
		BankName:     strings.TrimSpace(row[2]),
		City:         strings.TrimSpace(row[3]),
		CountryName:  countryNames[countryCode],
	}

	return rec, true
}

func loadCountryNames(source string) (map[string]string, error) {
	body, err := openSource(source)
	if err != nil {
		return nil, fmt.Errorf("load countries source: %w", err)
	}
	defer body.Close()

	reader := csv.NewReader(body)
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read countries header: %w", err)
	}

	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	codeCol, ok := idx["alpha-2"]
	if !ok {
		return nil, fmt.Errorf("countries CSV missing alpha-2 column")
	}

	nameCol, ok := idx["name"]
	if !ok {
		return nil, fmt.Errorf("countries CSV missing name column")
	}

	countries := make(map[string]string)
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if codeCol >= len(row) || nameCol >= len(row) {
			continue
		}

		code := strings.ToUpper(strings.TrimSpace(row[codeCol]))
		name := strings.TrimSpace(row[nameCol])
		if code == "" || name == "" {
			continue
		}

		countries[code] = name
	}

	return countries, nil
}

// parseRow converts one CSV row into a RawSWIFTRecord. Returns false if the
// row should be skipped (missing or malformed SWIFT code).
func parseRow(row []string, cols colIndices) (RawSWIFTRecord, bool) {
	if cols.swift >= len(row) {
		return RawSWIFTRecord{}, false
	}

	raw := strings.ToUpper(strings.TrimSpace(row[cols.swift]))
	if !isValidBIC(raw) {
		return RawSWIFTRecord{}, false
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
	if cols.bankName >= 0 && cols.bankName < len(row) {
		rec.BankName = strings.TrimSpace(row[cols.bankName])
	}
	if cols.city >= 0 && cols.city < len(row) {
		rec.City = strings.TrimSpace(row[cols.city])
	}
	if cols.countryName >= 0 && cols.countryName < len(row) {
		rec.CountryName = strings.TrimSpace(row[cols.countryName])
	}
	return rec, true
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
	return allAlpha(bic, 0, 4) && // bank code: letters only
		allAlpha(bic, 4, 6) && // country code: letters only
		allAlphanumeric(bic, 6, 8) && // location code: alphanumeric
		(len(bic) == 8 || allAlphanumeric(bic, 8, 11)) // branch code: alphanumeric if present
}

// allAlpha reports whether every byte in bic[from:to] is an uppercase ASCII letter.
func allAlpha(bic string, from, to int) bool {
	for i := from; i < to; i++ {
		if bic[i] < 'A' || bic[i] > 'Z' {
			return false
		}
	}
	return true
}

// allAlphanumeric reports whether every byte in bic[from:to] is an uppercase
// ASCII letter or decimal digit.
func allAlphanumeric(bic string, from, to int) bool {
	for i := from; i < to; i++ {
		b := bic[i]
		if (b < 'A' || b > 'Z') && (b < '0' || b > '9') {
			return false
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
