# seed-commodities

Seeds the `commodity_price_history` table from **FRED** (Federal Reserve Economic Data).

## Data Source

[FRED](https://fred.stlouisfed.org) provides free, public CSV downloads for commodity price series — no API key required. Each commodity is mapped to a specific FRED series ID. Daily and monthly observations are averaged per calendar year to produce annual averages stored in the database.

## Supported Commodities

| Slug          | Name              | FRED Series           | Unit        |
|---------------|-------------------|-----------------------|-------------|
| `gold`        | Gold              | GOLDAMGBD228NLBM      | oz          |
| `silver`      | Silver            | SLVPRUSD              | oz          |
| `platinum`    | Platinum          | PLTNUMGBD228NLBM      | oz          |
| `palladium`   | Palladium         | PALLADIUMGBD228NLBM   | oz          |
| `oil`         | Crude Oil (WTI)   | DCOILWTICO            | barrel      |
| `brent`       | Brent Crude       | DCOILBRENTEU          | barrel      |
| `natural-gas` | Natural Gas       | MHHNGSP               | mmbtu       |
| `copper`      | Copper            | PCOPPUSDM             | lb          |
| `aluminum`    | Aluminum          | PALUMUSDM             | metric_ton  |
| `wheat`       | Wheat             | PWHEAMTUSDM           | metric_ton  |
| `corn`        | Corn              | PMAIZEUSDM            | metric_ton  |
| `soybeans`    | Soybeans          | PSOYBUSDM             | metric_ton  |
| `coffee`      | Coffee            | PCOFFOTMUSDM          | lb          |
| `sugar`       | Sugar             | PSUGAISAUSDM          | lb          |
| `cotton`      | Cotton            | PCOTTINDUSDM          | lb          |
| `cocoa`       | Cocoa             | PCOCOAUSDM            | metric_ton  |

## Usage

Run inside the API Docker container so the `db` hostname resolves:

```bash
# Preview what would be seeded (no DB writes)
docker exec requiem-dev-api-1 go run ./cmd/seed-commodities --dry-run

# Seed the database
docker exec requiem-dev-api-1 go run ./cmd/seed-commodities \
  --db-url "postgres://requiem:requiem@db:5432/requiem"

# Verbose output (print each record)
docker exec requiem-dev-api-1 go run ./cmd/seed-commodities \
  --db-url "postgres://requiem:requiem@db:5432/requiem" \
  --verbose
```

## Re-seeding

Re-running the command is safe — it uses `INSERT ... ON CONFLICT DO UPDATE` so existing rows are updated in place. Run annually (or whenever FRED updates their series) to refresh prices.

## Unit Conversions

Some FRED series report in different units than the API's display unit:

| Commodity | FRED Unit       | API Unit   | Conversion     |
|-----------|----------------|------------|----------------|
| Copper    | USD/metric ton | USD/lb     | ÷ 2204.623     |
| Coffee    | USD/kg         | USD/lb     | ÷ 2.20462      |
| Sugar     | US cents/lb    | USD/lb     | ÷ 100          |
| Cotton    | US cents/lb    | USD/lb     | ÷ 100          |
