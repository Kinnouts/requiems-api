# seed-inflation

CLI tool that downloads historical CPI inflation data from the World Bank API
and upserts it into the `inflation_data` PostgreSQL table.

Run this after the initial deployment and whenever you want to refresh the data
(e.g. when the World Bank publishes new annual figures). The Go migration that
creates the table runs automatically on API startup.

## Usage

**Production** — run from `infra/docker/` as a one-off container on the compose network:

```bash
docker run --rm \
  --network requiem-backend_default \
  --env-file .env \
  -v $(pwd)/../../apps/api:/app \
  -w /app \
  golang:1.26-alpine \
  go run ./cmd/seed-inflation
```

Add `--dry-run` or `--verbose` at the end as needed.

**Development** — Go toolchain is available in the dev container:

```bash
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation --dry-run
docker exec requiem-dev-api-1 go run ./cmd/seed-inflation --verbose
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
