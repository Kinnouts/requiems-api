# seed-swift

Loads a vendored SWIFT/BIC code CSV dataset and upserts the records into the
`swift_codes` PostgreSQL table.

## Usage

Run inside the API Docker container so the `db` hostname resolves:

```bash
# Dry run — parse and print stats without writing to the database
docker compose exec api /app/seed-swift --dry-run

# Seed with default connection string (from DATABASE_URL env var)
docker compose exec api /app/seed-swift

# Seed with explicit connection string
docker compose exec api sh -lc '/app/seed-swift --db-url "$DATABASE_URL"'

# Verbose output — print each record as it is processed
docker compose exec api /app/seed-swift --verbose

# Optional remote source override
docker compose exec api /app/seed-swift \
  --url "https://example.com/swift-codes.csv"
```

## Flags

| Flag        | Default               | Description                                                |
| ----------- | --------------------- | ---------------------------------------------------------- |
| `--db-url`  | `$DATABASE_URL`       | PostgreSQL connection string (required unless `--dry-run`) |
| `--dry-run` | `false`               | Parse and print stats without writing to the database      |
| `--verbose` | `false`               | Log each record as it is processed                         |
| `--source`  | `dbs/swift_codes.csv` | Path to a local SWIFT codes CSV file                       |
| `--url`     | `""`                  | Optional SWIFT codes CSV URL (overrides `--source`)        |

## Data Source

The default CSV is loaded from `apps/api/dbs/swift_codes.csv` (copied to
`/app/dbs/swift_codes.csv` in production images). Use `--source` for another
local file or `--url` to substitute a remote CSV with these column names
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
