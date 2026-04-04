Please refer to [requiems.xyz/docs](https://requiems.xyz/docs) for the full
documentation. This document is mainly for feature tracking and statuses.

For developers:

## Directory Structure

```
docs/apis/
├── readme.md                   # This file — index and sync guide
├── status.md                   # Overall status dashboard (all APIs, all categories)
└── {category}/
    ├── readme.md               # Per-category status list with routes and credit costs
    └── {api}.md                # Individual API documentation
```

## API Status — Where to Look

There are 5 places that track API status. They must be kept in sync:

| Location                                                                                       | What it tracks                                          | Vocabulary                                               |
| ---------------------------------------------------------------------------------------------- | ------------------------------------------------------- | -------------------------------------------------------- |
| [`docs/apis/status.md`](./status.md)                                                           | All APIs across all categories                          | `planned / partial / mvp / complete`                     |
| [`docs/apis/{category}/readme.md`](./)                                                         | Per-category endpoint list with routes                  | `planned / mvp`                                          |
| [`apps/dashboard/config/api_catalog.yml`](../../apps/dashboard/config/api_catalog.yml)         | Rails UI catalog — what appears in the public directory | `status: live` per API; `coming_soon: true` per category |
| [`apps/dashboard/app/helpers/apis_helper.rb`](../../apps/dashboard/app/helpers/apis_helper.rb) | Status badge display logic                              | `live / beta / deprecated`                               |
| Go router files (see below)                                                                    | **Source of truth** — what is actually running          | Route registered = live                                  |

### Go Routers (source of truth)

| File                                        | Registered services                     |
| ------------------------------------------- | --------------------------------------- |
| `apps/api/app/routes_v1.go`                 | Top-level domain mounts                 |
| `apps/api/services/convert/router.go`       | base64 encode/decode                    |
| `apps/api/services/text/router.go`          | advice, lorem, profanity, quotes, words |
| `apps/api/services/entertainment/router.go` | horoscope                               |
| `apps/api/services/misc/router.go`          | counter, unit conversion                |
| `apps/api/services/places/router.go`        | timezone, working-days, holidays        |
| `apps/api/services/tech/router.go`          | phone, password, useragent              |
| `apps/api/services/email/router.go`         | disposable email                        |

## Sync Checklist — Adding a New Live API

When an API goes from planned → live, update all 5 locations in order:

1. **Go router** — register the route in `apps/api/services/{domain}/router.go`
2. **`api_catalog.yml`** — add entry with `status: live` (and remove
   `coming_soon: true` from the category if needed)
3. **`docs/apis/status.md`** — update status emoji and progress counts
4. **`docs/apis/{category}/readme.md`** — update the individual entry and
   category statistics
5. **`docs/apis/{category}/{api}.md`** — update the individual API doc status
   field
