# seed-bins

CLI tool that downloads BIN/IIN data from two open-source datasets, normalises
and merges the records, and upserts them into the `bin_data` PostgreSQL table.

Run this after the initial deployment and whenever you want to refresh the data.
The Go migration that creates the table runs automatically on API startup.

## Usage

**Production** — run from `infra/docker/` as a one-off container on the compose network:

```bash
docker run --rm \
  --network requiem-backend_default \
  --env-file .env \
  -v $(pwd)/../../apps/api:/app \
  -w /app \
  golang:1.26-alpine \
  go run ./cmd/seed-bins
```

Add `--dry-run` or `--verbose` at the end as needed.

**Development** — Go toolchain is available in the dev container:

```bash
docker exec requiem-dev-api-1 go run ./cmd/seed-bins
docker exec requiem-dev-api-1 go run ./cmd/seed-bins --dry-run
docker exec requiem-dev-api-1 go run ./cmd/seed-bins --verbose
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
