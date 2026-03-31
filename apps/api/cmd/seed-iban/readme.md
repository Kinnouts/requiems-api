# seed-iban

CLI tool that downloads the IBAN country registry from the
[php-iban project](https://github.com/globalcitizen/php-iban) — a
maintained, machine-readable mirror of the official SWIFT IBAN Registry —
and upserts it into the `iban_countries` PostgreSQL table.

Run this after the initial deployment and whenever the IBAN Registry is
updated (new countries join the scheme, format changes).

## Usage

**Development** — Go toolchain is available in the dev container:

```bash
docker exec requiem-dev-api-1 go run ./cmd/seed-iban
docker exec requiem-dev-api-1 go run ./cmd/seed-iban --dry-run
docker exec requiem-dev-api-1 go run ./cmd/seed-iban --verbose
```

**Production** — run from `infra/docker/` as a one-off container:

```bash
docker run --rm \
  --network requiem-backend_default \
  --env-file .env \
  -v $(pwd)/../../apps/api:/app \
  -w /app \
  golang:1.26-alpine \
  go run ./cmd/seed-iban
```

## Flags

| Flag        | Default              | Description                               |
| ----------- | -------------------- | ----------------------------------------- |
| `--db-url`  | `$DATABASE_URL`      | PostgreSQL connection string              |
| `--dry-run` | false                | Parse and print without writing to DB     |
| `--verbose` | false                | Log each country record as it is processed |
| `--url`     | php-iban GitHub URL  | Override the registry source URL          |

## Data Source

[globalcitizen/php-iban](https://github.com/globalcitizen/php-iban) —
`registry.txt`

This file is a pipe-separated mirror of the SWIFT IBAN Registry
(ISO 13616). It includes every country that has officially adopted the
IBAN scheme, with:

- Expected total IBAN length
- BBAN format string (SWIFT notation, e.g. `8!n10!n`)
- 0-indexed bank and branch identifier offsets within the BBAN
- SEPA membership flag

## What Gets Seeded

Each row in `iban_countries` represents one country with:

| Column           | Example (DE) | Description                              |
| ---------------- | ------------ | ---------------------------------------- |
| `country_code`   | `DE`         | ISO 3166-1 alpha-2                       |
| `country_name`   | `Germany`    | Full country name                        |
| `iban_length`    | `22`         | Expected total IBAN character count      |
| `bban_format`    | `8!n10!n`    | SWIFT BBAN format string                 |
| `bank_offset`    | `0`          | 0-indexed offset of bank code in BBAN    |
| `bank_length`    | `8`          | Length of bank code (Bankleitzahl)       |
| `account_offset` | `8`          | 0-indexed offset of account number in BBAN |
| `account_length` | `10`         | Length of account number                 |
| `sepa_member`    | `true`       | Whether the country is in SEPA           |

## Pipeline

1. **Fetch** — download registry file (30s timeout)
2. **Parse** — read pipe-separated lines, skip header
3. **Convert** — derive account offset from bank/branch end positions
4. **Upsert** — staging table → COPY bulk insert → `ON CONFLICT DO UPDATE`
