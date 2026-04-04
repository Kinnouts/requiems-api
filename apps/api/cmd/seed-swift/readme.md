# seed-swift

Downloads and normalises a full SWIFT/BIC dataset, then upserts records into
the `swift_codes` PostgreSQL table.

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

# Optional country mapping override (ISO alpha-2 to country name)
docker compose exec api /app/seed-swift \
  --countries "https://example.com/countries.csv"
```

## Flags

| Flag          | Default         | Description                                                |
| ------------- | --------------- | ---------------------------------------------------------- |
| `--db-url`    | `$DATABASE_URL` | PostgreSQL connection string (required unless `--dry-run`) |
| `--dry-run`   | `false`         | Parse and print stats without writing to the database      |
| `--verbose`   | `false`         | Log each record as it is processed                         |
| `--source`    | upstream URL    | Path or URL for SWIFT codes CSV source                     |
| `--countries` | upstream URL    | Path or URL for country mapping CSV (`name`, `alpha-2`)    |
| `--url`       | `""`            | Optional SWIFT codes CSV URL (overrides `--source`)        |

## Data Source

By default, the seeder pulls:

- a full SWIFT source CSV
- a country mapping CSV (to populate `country_name`)

Supported SWIFT source formats:

- Header-based CSV with flexible columns (case-insensitive, any order):

- `swift_code` or `bic` — the SWIFT/BIC code (8 or 11 characters)
- `bank_name` or `name` — institution name
- `city` — branch city
- `country_name` or `country` — country name

- No-header 4-column format:
  - `country_code`,`swift_code`,`bank_name`,`city`

8-character codes (primary offices) are automatically expanded to 11 characters
by appending `XXX`.

## Upsert Behaviour

- New codes are inserted.
- Existing codes update `bank_name`, `city`, `country_name`, and `last_updated`.
- Structural fields (`bank_code`, `country_code`, `location_code`, `branch_code`)
  are derived from the BIC itself and are not updated on conflict.
