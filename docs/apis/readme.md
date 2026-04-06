Please refer to [requiems.xyz/docs](https://requiems.xyz/docs) for the full
documentation.

This document is mainly for feature tracking and statuses.

For developers:

## Directory Structure

```
docs/apis/
├── readme.md                   # This file — index and sync guide
└── {category}/
    ├── readme.md               # Per-category status list with routes and credit costs
    └── {api}.md                # Individual API documentation
```

Note: The physical directory layout in `docs/apis/` predates the current
category structure. The canonical category groupings are defined in
`apps/dashboard/config/api_catalog.yml`.

## Current Categories

| Category              | APIs | Directory        |
| --------------------- | ---- | ---------------- |
| Finance               | 8    | `finance/`       |
| Validation            | 3    | `validation/`    |
| Networking & Internet | 7    | `networking/`    |
| Places                | 7    | `places/`        |
| Text & Language       | 9    | `text/`          |
| Developer Tools       | 11   | `technology/`    |
| Entertainment         | 9    | `entertainment/` |
| Health                | 1    | `health/`        |

## API Status — Where to Look

There are 5 places that track API status. They must be kept in sync:

| Location                                                                                       | What it tracks                                          | Vocabulary                                               |
| ---------------------------------------------------------------------------------------------- | ------------------------------------------------------- | -------------------------------------------------------- |
| [`docs/apis/status.md`](./status.md)                                                           | All APIs across all categories                          | `mvp / complete`                                         |
| [`docs/apis/{category}/readme.md`](./)                                                         | Per-category endpoint list with routes                  | `mvp / complete`                                         |
| [`apps/dashboard/config/api_catalog.yml`](../../apps/dashboard/config/api_catalog.yml)         | Rails UI catalog — what appears in the public directory | `status: live` per API; `coming_soon: true` per category |
| [`apps/dashboard/app/helpers/apis_helper.rb`](../../apps/dashboard/app/helpers/apis_helper.rb) | Status badge display logic                              | `live / beta / deprecated`                               |
| Go router files (see below)                                                                    | **Source of truth** — what is actually running          | Route registered = live                                  |

### Go Routers (source of truth)

| File                                        | Registered services                                                                                                         |
| ------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `apps/api/app/routes_v1.go`                 | Top-level domain mounts                                                                                                     |
| `apps/api/services/convert/router.go`       | base64, number base conversion, markdown, data format                                                                       |
| `apps/api/services/text/router.go`          | advice, lorem, profanity, quotes, words, spell-check, thesaurus, dictionary, sentiment, language detection, text similarity |
| `apps/api/services/entertainment/router.go` | horoscope, trivia, facts, emoji, sudoku, dad-jokes, chuck-norris                                                            |
| `apps/api/services/misc/router.go`          | counter, unit conversion                                                                                                    |
| `apps/api/services/places/router.go`        | timezone, working-days, holidays, world-time, geocode, postal-code, cities                                                  |
| `apps/api/services/tech/router.go`          | phone, password, useragent, qr, barcode, random-user, ip, asn, vpn, whois, domain, mx, color                                |
| `apps/api/services/email/router.go`         | disposable email, email validate, email normalize                                                                           |
| `apps/api/services/finance/router.go`       | exchange-rate, crypto, bin, iban, swift, mortgage, commodities, inflation                                                   |
| `apps/api/services/health/router.go`        | fitness exercises                                                                                                           |

## Sync Checklist — Adding a New Live API

When an API goes live, update all 5 locations in order:

1. **Go router** — register the route in `apps/api/services/{domain}/router.go`
2. **`api_catalog.yml`** — add entry with `status: live`
3. **`docs/apis/status.md`** — add to the relevant category section
4. **`docs/apis/{category}/readme.md`** — add entry with route and credit cost
5. **`docs/apis/{category}/{api}.md`** — update the individual API doc status
   field
