package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RawIBANCountry holds one parsed entry from the IBAN registry.
//
// Bank and account positions are 0-indexed offsets within the BBAN (not the
// full IBAN string). The registry stores explicit bank/branch positions; the
// account offset is derived as the portion of the BBAN that follows the bank
// and branch identifiers.
type RawIBANCountry struct {
	CountryCode string
	CountryName string
	IBANLength  int
	BBANLength  int
	BBANFormat  string
	// 0-indexed positions within the BBAN (per SWIFT IBAN Registry).
	BankIDStart    int
	BankIDEnd      int
	BranchIDStart  int
	BranchIDEnd    int
	SEPAMember     bool
}

// BankOffset returns the 0-indexed start position of the bank code within the
// BBAN. Positions are already BBAN-relative in this registry format.
func (r RawIBANCountry) BankOffset() int {
	return r.BankIDStart
}

// BankLength returns the number of characters in the bank code.
func (r RawIBANCountry) BankLength() int {
	if r.BankIDEnd < r.BankIDStart {
		return 0
	}
	return r.BankIDEnd - r.BankIDStart + 1
}

// AccountOffset returns the 0-indexed start position of the account number
// within the BBAN. The account is defined as the portion of the BBAN that
// follows the bank and branch identifiers.
func (r RawIBANCountry) AccountOffset() int {
	// Branch comes after bank in the BBAN. If a valid branch range is defined
	// and it extends beyond the bank end, account starts after the branch.
	if r.BranchIDEnd >= r.BranchIDStart && r.BranchIDEnd > r.BankIDEnd {
		return r.BranchIDEnd + 1
	}
	return r.BankIDEnd + 1
}

// AccountLength returns the number of characters in the account number.
func (r RawIBANCountry) AccountLength() int {
	acctEnd := r.BBANLength - 1
	offset := r.AccountOffset()
	if offset > acctEnd {
		return 0
	}
	return acctEnd - offset + 1
}

// fetchAndParse downloads the php-iban IBAN registry and returns parsed
// country records.
//
// The registry uses a pipe-separated format. The first line is a header; all
// subsequent lines are one country per line.
//
// Column layout (0-indexed after splitting on "|"):
//
//	0:  ISO 3166-1 alpha-2 country code
//	1:  country name
//	4:  BBAN format in SWIFT notation (e.g. "8!n10!n")
//	6:  BBAN length
//	10: IBAN length
//	11: BBAN bank identifier start offset (0-indexed in BBAN; empty = not defined)
//	12: BBAN bank identifier stop offset (0-indexed in BBAN, inclusive)
//	13: BBAN branch identifier start offset
//	14: BBAN branch identifier stop offset
//	16: SEPA member ("1" = yes, "0" = no)
func fetchAndParse(url string) ([]RawIBANCountry, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download: HTTP %d", resp.StatusCode)
	}

	var countries []RawIBANCountry
	scanner := bufio.NewScanner(resp.Body)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// First line is the header; skip it.
		if lineNum == 1 {
			continue
		}

		fields := strings.Split(line, "|")
		if len(fields) < 17 {
			continue
		}

		ibanLen, err := strconv.Atoi(strings.TrimSpace(fields[10]))
		if err != nil || ibanLen < 5 {
			continue
		}

		bbanLen, err := strconv.Atoi(strings.TrimSpace(fields[6]))
		if err != nil || bbanLen < 1 {
			continue
		}

		countries = append(countries, RawIBANCountry{
			CountryCode:   strings.ToUpper(strings.TrimSpace(fields[0])),
			CountryName:   strings.TrimSpace(fields[1]),
			IBANLength:    ibanLen,
			BBANLength:    bbanLen,
			BBANFormat:    strings.TrimSpace(fields[4]),
			BankIDStart:   parseOptInt(fields[11]),
			BankIDEnd:     parseOptInt(fields[12]),
			BranchIDStart: parseOptInt(fields[13]),
			BranchIDEnd:   parseOptInt(fields[14]),
			SEPAMember:    strings.TrimSpace(fields[16]) == "1",
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	if len(countries) == 0 {
		return nil, fmt.Errorf("no country records found — the registry format may have changed")
	}

	return countries, nil
}

// parseOptInt parses s as an integer, returning 0 for empty, "N/A", or
// unparseable values.
func parseOptInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" || s == "N/A" || s == "-" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// printStats prints a summary of the parsed IBAN dataset to stdout.
func printStats(countries []RawIBANCountry) {
	sepa := 0
	withBank := 0
	for _, c := range countries {
		if c.SEPAMember {
			sepa++
		}
		if c.BankLength() > 0 {
			withBank++
		}
	}

	fmt.Printf("\n=== IBAN Seed Stats ===\n")
	fmt.Printf("Total countries:      %d\n", len(countries))
	fmt.Printf("SEPA members:         %d\n", sepa)
	fmt.Printf("With bank ID defined: %d\n", withBank)
	fmt.Printf("\nSample entries:\n")
	for i, c := range countries {
		if i >= 8 {
			break
		}
		fmt.Printf("  %-2s  %-32s  iban_len=%-3d  bban=%-12s  bank=%d+%d  acct=%d+%d\n",
			c.CountryCode, c.CountryName, c.IBANLength, c.BBANFormat,
			c.BankOffset(), c.BankLength(),
			c.AccountOffset(), c.AccountLength(),
		)
	}
}
