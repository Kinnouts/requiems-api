package main

import (
	"math"
	"strings"
)

// normalise applies canonical transformations to a raw record.
// Records with non-2-letter country codes (e.g. World Bank regional aggregates
// like "1W" or "EAP") are flagged by setting CountryCode to "" so the caller
// can filter them out.
func normalise(r RawInflationRecord) RawInflationRecord {
	r.CountryCode = strings.ToUpper(strings.TrimSpace(r.CountryCode))
	r.CountryName = strings.TrimSpace(r.CountryName)

	// Discard World Bank aggregate regions (non-2-letter codes).
	if len(r.CountryCode) != 2 {
		r.CountryCode = ""
	}

	// Round to 4 decimal places to match NUMERIC(8,4) storage.
	r.Rate = math.Round(r.Rate*10000) / 10000

	return r
}
