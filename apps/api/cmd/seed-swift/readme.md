# seed-swift

Downloads an open-source SWIFT/BIC code CSV dataset and upserts the records
into the `swift_codes` PostgreSQL table.

## Usage

Run inside the API Docker container so the `db` hostname resolves:

```bash
# Dry run — parse and print stats without writing to the database
docker exec requiem-dev-api-1 go run ./cmd/seed-swift --dry-run

# Seed with default connection string (from DATABASE_URL env var)
docker exec requiem-dev-api-1 go run ./cmd/seed-swift

# Seed with explicit connection string
docker exec requiem-dev-api-1 go run ./cmd/seed-swift \
  --db-url "postgres://requiem:requiem@db:5432/requiem"

# Verbose output — print each record as it is processed
docker exec requiem-dev-api-1 go run ./cmd/seed-swift --verbose
```

## Flags

| Flag        | Default               | Description                                              |
| ----------- | --------------------- | -------------------------------------------------------- |
| `--db-url`  | `$DATABASE_URL`       | PostgreSQL connection string (required unless `--dry-run`) |
| `--dry-run` | `false`               | Parse and print stats without writing to the database    |
| `--verbose` | `false`               | Log each record as it is processed                       |
| `--url`     | See source below      | Override the SWIFT codes CSV URL                         |

## Data Source

The default CSV is fetched from a community-maintained dataset of SWIFT/BIC
codes. Use `--url` to substitute any compatible CSV with these column names
(case-insensitive, any order):

- `swift_code` or `bic` — the SWIFT/BIC code (8 or 11 characters)
- `bank_name` or `name` — institution name
- `city` — branch city
- `country_name` or `country` — country name

8-character codes (primary offices) are automatically expanded to 11 characters
by appending `XXX`.

## Upsert Behaviour

- New codes are inserted.
- Existing codes update `bank_name`, `city`, `country_name`, and `last_updated`.
- Structural fields (`bank_code`, `country_code`, `location_code`, `branch_code`)
  are derived from the BIC itself and are not updated on conflict.
