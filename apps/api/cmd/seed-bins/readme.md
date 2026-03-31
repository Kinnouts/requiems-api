# seed-bins

CLI tool that downloads BIN/IIN data from two open-source datasets, normalises
and merges the records, and upserts them into the `bin_data` PostgreSQL table.

Run this after the initial deployment and whenever you want to refresh the data.
The Go migration that creates the table runs automatically on API startup.

## Usage

Run inside the API Docker container (required — uses the `db` hostname):

```bash
# Basic run (uses DATABASE_URL from container env)
docker exec requiem-dev-api-1 go run ./cmd/seed-bins

# Dry-run: fetch and normalise but do not write to the database
docker exec requiem-dev-api-1 go run ./cmd/seed-bins --dry-run

# Verbose: print each record as it is processed
docker exec requiem-dev-api-1 go run ./cmd/seed-bins --verbose

# Custom database URL
docker exec requiem-dev-api-1 go run ./cmd/seed-bins --db-url "postgres://..."
```

In production (container is named differently):

```bash
docker compose exec api go run ./cmd/seed-bins
```

## Flags

| Flag                  | Default         | Description                               |
| --------------------- | --------------- | ----------------------------------------- |
| `--db-url`            | `$DATABASE_URL` | PostgreSQL connection string              |
| `--dry-run`           | false           | Parse and normalise without writing to DB |
| `--verbose`           | false           | Log each record as it is processed        |
| `--url-iannuttall`    | GitHub CSV URL  | Override source A                         |
| `--url-venelinkochev` | GitHub CSV URL  | Override source B                         |

## Data Sources

- [iannuttall/binlist-data](https://github.com/iannuttall/binlist-data) —
  confidence 0.75
- [venelinkochev/bin-list-data](https://github.com/venelinkochev/bin-list-data)
  — confidence 0.80

Records from both sources are merged by BIN prefix. When both sources have the
same prefix, fields are resolved by confidence score with a +0.10 agreement
bonus for multi-source matches.

## Pipeline

1. **Fetch** — download CSV from each source (120s timeout)
2. **Parse** — convert rows to `RawBINRecord` structs
3. **Normalise** — canonical scheme/type/level names, validate country codes,
   compute confidence scores (`normalize.go`)
4. **Merge** — deduplicate by BIN prefix, field-level conflict resolution
   (`normalize.go`)
5. **Upsert** — staging table → `COPY` bulk insert → `ON CONFLICT DO UPDATE`
   with confidence-weighted field merging (`db.go`)
