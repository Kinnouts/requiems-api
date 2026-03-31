# Commodity Prices API

Returns historical and current annual average prices for 16 major commodities — precious metals, energy, and agricultural goods. Data is sourced from [FRED](https://fred.stlouisfed.org) (Federal Reserve Economic Data).

## Endpoint

`GET /v1/finance/commodities/{commodity}`

## Supported Commodities

| Slug          | Name              | Unit        | FRED Series           |
|---------------|-------------------|-------------|-----------------------|
| `gold`        | Gold              | oz          | GOLDAMGBD228NLBM      |
| `silver`      | Silver            | oz          | SLVPRUSD              |
| `platinum`    | Platinum          | oz          | PLTNUMGBD228NLBM      |
| `palladium`   | Palladium         | oz          | PALLADIUMGBD228NLBM   |
| `oil`         | Crude Oil (WTI)   | barrel      | DCOILWTICO            |
| `brent`       | Brent Crude       | barrel      | DCOILBRENTEU          |
| `natural-gas` | Natural Gas       | mmbtu       | MHHNGSP               |
| `copper`      | Copper            | lb          | PCOPPUSDM             |
| `aluminum`    | Aluminum          | metric_ton  | PALUMUSDM             |
| `wheat`       | Wheat             | metric_ton  | PWHEAMTUSDM           |
| `corn`        | Corn              | metric_ton  | PMAIZEUSDM            |
| `soybeans`    | Soybeans          | metric_ton  | PSOYBUSDM             |
| `coffee`      | Coffee            | lb          | PCOFFOTMUSDM          |
| `sugar`       | Sugar             | lb          | PSUGAISAUSDM          |
| `cotton`      | Cotton            | lb          | PCOTTINDUSDM          |
| `cocoa`       | Cocoa             | metric_ton  | PCOCOAUSDM            |

## Response Envelope

All responses use the standard envelope:

```json
{
  "data": {
    "commodity": "gold",
    "name": "Gold",
    "price": 2386.3300,
    "unit": "oz",
    "currency": "USD",
    "change_24h": 23.01,
    "historical": [
      { "period": "2023", "price": 1940.5400 },
      { "period": "2022", "price": 1800.1200 }
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Notes

- **Prices are annual averages**, not real-time spot prices. Daily and monthly FRED observations are averaged per calendar year.
- **`change_24h`** is the year-over-year percentage change between the latest annual average and the prior year's average, despite the field name.
- **Up to 10 years** of historical data are returned alongside the current value.
- **Precious metals** (gold, silver, platinum, palladium) are priced per troy ounce.
- **Copper and coffee** are converted from FRED's USD/metric ton and USD/kg to USD/lb.
- **Sugar and cotton** are converted from FRED's US cents/lb to USD/lb.

## Error Codes

| Code             | Status | When                                       |
|------------------|--------|--------------------------------------------|
| `not_found`      | 404    | Slug is not in the database (run seed CLI) |
| `internal_error` | 500    | Unexpected server error                    |

## Seeding

The table must be populated with the seed CLI before the endpoint returns data:

```bash
docker exec requiem-dev-api-1 go run ./cmd/seed-commodities \
  --db-url "postgres://requiem:requiem@db:5432/requiem"
```

See [`cmd/seed-commodities/readme.md`](../../../apps/api/cmd/seed-commodities/readme.md) for full options.
