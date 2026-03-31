package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CommodityConfig defines one commodity and its FRED data source.
type CommodityConfig struct {
	Slug       string  // primary key slug used in the API (e.g. "gold")
	Name       string  // human-readable name (e.g. "Gold")
	SeriesID   string  // FRED series identifier
	Unit       string  // display unit for the API response (e.g. "oz", "barrel")
	Currency   string  // always "USD"
	ConvFactor float64 // multiply raw FRED value by this to get the target unit's price
}

// CommodityRecord is a single annual average price ready for upsert.
type CommodityRecord struct {
	Slug     string
	Name     string
	Unit     string
	Currency string
	Year     int
	Price    float64
}

// commodities is the authoritative list of supported commodity slugs.
//
// FRED series notes:
//   - Daily series (e.g. gold, silver, oil): averaged per calendar year.
//   - Monthly series (e.g. natural gas, copper, agricultural): averaged per year.
//   - ConvFactor converts FRED's native unit to the API's display unit.
//     FRED Copper/Aluminum/etc. is USD/metric ton; Copper display unit is lb → ÷ 2204.623
//     FRED Sugar/Cotton is US cents/lb → ÷ 100 for USD/lb
//     FRED Coffee is USD/kg → ÷ 2.20462 for USD/lb
//     FRED Wheat/Corn/Soybeans is USD/metric ton → display is USD/metric ton (no conversion)
var commodities = []CommodityConfig{
	{Slug: "gold", Name: "Gold", SeriesID: "GOLDAMGBD228NLBM", Unit: "oz", Currency: "USD", ConvFactor: 1.0},
	{Slug: "silver", Name: "Silver", SeriesID: "SLVPRUSD", Unit: "oz", Currency: "USD", ConvFactor: 1.0},
	{Slug: "platinum", Name: "Platinum", SeriesID: "PLTNUMGBD228NLBM", Unit: "oz", Currency: "USD", ConvFactor: 1.0},
	{Slug: "palladium", Name: "Palladium", SeriesID: "PALLADIUMGBD228NLBM", Unit: "oz", Currency: "USD", ConvFactor: 1.0},
	{Slug: "oil", Name: "Crude Oil (WTI)", SeriesID: "DCOILWTICO", Unit: "barrel", Currency: "USD", ConvFactor: 1.0},
	{Slug: "brent", Name: "Brent Crude", SeriesID: "DCOILBRENTEU", Unit: "barrel", Currency: "USD", ConvFactor: 1.0},
	{Slug: "natural-gas", Name: "Natural Gas", SeriesID: "MHHNGSP", Unit: "mmbtu", Currency: "USD", ConvFactor: 1.0},
	{Slug: "copper", Name: "Copper", SeriesID: "PCOPPUSDM", Unit: "lb", Currency: "USD", ConvFactor: 1.0 / 2204.623},
	{Slug: "aluminum", Name: "Aluminum", SeriesID: "PALUMUSDM", Unit: "metric_ton", Currency: "USD", ConvFactor: 1.0},
	{Slug: "wheat", Name: "Wheat", SeriesID: "PWHEAMTUSDM", Unit: "metric_ton", Currency: "USD", ConvFactor: 1.0},
	{Slug: "corn", Name: "Corn", SeriesID: "PMAIZEUSDM", Unit: "metric_ton", Currency: "USD", ConvFactor: 1.0},
	{Slug: "soybeans", Name: "Soybeans", SeriesID: "PSOYBUSDM", Unit: "metric_ton", Currency: "USD", ConvFactor: 1.0},
	{Slug: "coffee", Name: "Coffee", SeriesID: "PCOFFOTMUSDM", Unit: "lb", Currency: "USD", ConvFactor: 1.0 / 2.20462},
	{Slug: "sugar", Name: "Sugar", SeriesID: "PSUGAISAUSDM", Unit: "lb", Currency: "USD", ConvFactor: 1.0 / 100},
	{Slug: "cotton", Name: "Cotton", SeriesID: "PCOTTINDUSDM", Unit: "lb", Currency: "USD", ConvFactor: 1.0 / 100},
	{Slug: "cocoa", Name: "Cocoa", SeriesID: "PCOCOAUSDM", Unit: "metric_ton", Currency: "USD", ConvFactor: 1.0},
}

// fredBaseURL is the FRED CSV download endpoint (no API key required).
const fredBaseURL = "https://fred.stlouisfed.org/graph/fredgraph.csv"

// fetchAndAggregate downloads the FRED CSV series for cfg, parses date-value
// pairs, skips missing observations ("."), groups by calendar year, and returns
// the annual averages as CommodityRecord slices.
func fetchAndAggregate(cfg CommodityConfig) ([]CommodityRecord, error) {
	url := fredBaseURL + "?id=" + cfg.SeriesID

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("download: HTTP %d — %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// yearAcc accumulates (sum, count) per year for averaging.
	type acc struct {
		sum   float64
		count int
	}
	byYear := make(map[int]*acc)

	r := csv.NewReader(resp.Body)

	// Skip the header row (DATE,VALUE).
	if _, err = r.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv: %w", err)
		}
		if len(record) < 2 {
			continue
		}

		dateStr := strings.TrimSpace(record[0])
		valStr := strings.TrimSpace(record[1])

		// FRED encodes missing observations as "."
		if valStr == "." || valStr == "" {
			continue
		}

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil || math.IsNaN(val) || math.IsInf(val, 0) || val <= 0 {
			continue
		}

		// dateStr is either YYYY-MM-DD (daily/monthly) or YYYY (annual).
		year, err := parseYear(dateStr)
		if err != nil || year < 1960 || year > time.Now().Year() {
			continue
		}

		if byYear[year] == nil {
			byYear[year] = &acc{}
		}
		byYear[year].sum += val
		byYear[year].count++
	}

	records := make([]CommodityRecord, 0, len(byYear))
	for year, a := range byYear {
		if a.count == 0 {
			continue
		}
		avg := a.sum / float64(a.count)
		price := math.Round(avg*cfg.ConvFactor*10000) / 10000 // 4 decimal places

		records = append(records, CommodityRecord{
			Slug:     cfg.Slug,
			Name:     cfg.Name,
			Unit:     cfg.Unit,
			Currency: cfg.Currency,
			Year:     year,
			Price:    price,
		})
	}

	return records, nil
}

// parseYear extracts the 4-digit year from a FRED date string.
// Handles both "YYYY-MM-DD" and "YYYY" formats.
func parseYear(s string) (int, error) {
	if len(s) < 4 {
		return 0, fmt.Errorf("date too short: %q", s)
	}
	return strconv.Atoi(s[:4])
}

// printStats prints a summary of the loaded records to stdout.
func printStats(records []CommodityRecord) {
	bySlug := make(map[string]int)
	for _, r := range records {
		bySlug[r.Slug]++
	}
	fmt.Printf("\n=== Commodity Seed Stats ===\n")
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Commodities:   %d\n", len(bySlug))
	fmt.Printf("\nYears per commodity:\n")
	for slug, count := range bySlug {
		fmt.Printf("  %-20s %d years\n", slug, count)
	}
}
