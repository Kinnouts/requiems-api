# seed-inflation

CLI tool that downloads historical CPI inflation data from the World Bank API
and upserts it into the `inflation_data` PostgreSQL table.

Run this after the initial deployment and whenever you want to refresh the data
(e.g. when the World Bank publishes new annual figures). The Go migration that
creates the table runs automatically on API startup.

## Usage

Run inside the API Docker container (required — uses the `db` hostname):

```bash
# Basic run (uses DATABASE_URL from container env)
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation

# Dry-run: fetch and normalise but do not write to the database
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation --dry-run

# Verbose: print each record as it is processed
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation --verbose

# Custom database URL
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation --db-url "postgres://..."
```

In production (container is named differently):

```bash
docker compose exec api go run ./cmd/seed-inflation
```

## Flags

| Flag        | Default            | Description                               |
| ----------- | ------------------ | ----------------------------------------- |
| `--db-url`  | `$DATABASE_URL`    | PostgreSQL connection string              |
| `--dry-run` | false              | Parse and normalise without writing to DB |
| `--verbose` | false              | Log each record as it is processed        |
| `--url`     | World Bank API URL | Override the data source URL              |

## Data Source

[World Bank — Inflation, consumer prices (annual %)](https://data.worldbank.org/indicator/FP.CPI.TOTL.ZG)

- ~241 countries, last 30 years of annual CPI data
- No API key required
- Regional aggregates (e.g. `1W`, `EAP`) are filtered out automatically

## Pipeline

1. **Fetch** — HTTP GET to World Bank API (30s timeout), decodes 2-element JSON
   envelope `[metadata, data_array]`
2. **Parse** — converts entries to `RawInflationRecord` structs, skips null
   values and malformed dates (`source.go`)
3. **Normalise** — uppercase country codes, filter non-2-letter regional
   aggregates, round rate to 4 decimal places (`normalize.go`)
4. **Upsert** — staging table → `COPY` bulk insert → `ON CONFLICT DO UPDATE`,
   newer seed always wins (`db.go`)
