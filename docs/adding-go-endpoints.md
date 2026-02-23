# Adding a New Endpoint to the Go Backend

This guide walks through every step required to ship a new endpoint: from writing the Go code to updating the dashboard catalog and API documentation. Follow it in order — the checklist at the end maps to each section.

---

## Architecture Refresher

Before writing any code, understand where the Go backend sits in the request flow:

```
User → Auth Gateway (api.requiems.xyz)
         ↓  validates API key, enforces rate limits
       Go Backend (apps/api, port 8080)
         ↓  business logic, database queries
       PostgreSQL
```

The Go backend **trusts the gateway completely**. It does not perform API key validation or rate limiting — that is the gateway's job. The backend only receives requests with a valid `X-Backend-Secret` header already attached.

### Domain Directory Layout

Every feature lives inside a domain:

```
apps/api/internal/
├── app/
│   ├── app.go            # Bootstrap: DB, Redis, middleware, router
│   └── routes_v1.go      # Mounts all /v1 domain routers
├── {domain}/
│   ├── router.go         # Wires services → chi.Router for this domain
│   └── {feature}/
│       ├── type.go           # Request + response types
│       ├── service.go        # Business logic
│       └── transport_http.go # HTTP handlers
├── platform/
│   ├── httpx/            # JSON, Error, Handle, BindQuery helpers
│   ├── config/           # Env config
│   ├── db/               # PostgreSQL connection + migrations
│   └── middleware/       # BackendSecretAuth
```

Existing top-level domains and their `/v1` prefixes:

| Domain folder | URL prefix |
|---|---|
| `text/` | `/v1/text` |
| `email/` | `/v1/email` |
| `entertainment/` | `/v1/entertainment` |
| `misc/` | `/v1/misc` |
| `places/` | `/v1/places` |

---

## Before Writing Any Code — Check for Existing Libraries

**Do not write a service from scratch if a battle-tested library already solves the problem.**

Go has a rich ecosystem and several `bobadilla-tech` packages already in `go.mod` that were purpose-built for this platform. Writing your own implementation of something that already exists adds maintenance burden, reintroduces bugs that libraries have already fixed, and makes the codebase harder to onboard.

### Check in this order

1. **`go.mod` — existing dependencies first**
   Before anything else, open `apps/api/go.mod` and scan the current dependency list. The problem may already be solved.

   Current domain-specific libraries in use:
   | Library | What it does |
   |---|---|
   | `github.com/bobadilla-tech/is-email-disposable` | Disposable email blocklist (embedded) |
   | `github.com/bobadilla-tech/lorelai` | Lorem ipsum text generation |
   | `github.com/bobadilla-tech/business-days-calculator` | Working-days calculations |
   | `github.com/go-playground/validator/v10` | Struct validation |
   | `github.com/golang-migrate/migrate/v4` | SQL migration runner |

2. **`bobadilla-tech` org on GitHub**
   Check whether a new `bobadilla-tech` package exists for the domain you are building. These are first-party libraries designed to slot directly into this platform.

3. **Well-known Go ecosystem packages**
   For common problems, prefer established packages over custom implementations:
   - Parsing / tokenising: `golang.org/x/text`, `mvdan.cc/xurls`, etc.
   - Cryptography: standard library `crypto/*` — never roll your own
   - Time zones / date math: standard library `time` + `golang.org/x/time`
   - HTTP client retries: `hashicorp/go-retryablehttp`
   - UUID generation: `google/uuid`

4. **Only write it yourself if:**
   - No library exists for the domain
   - Available libraries are unmaintained, have known CVEs, or have prohibitive licenses
   - The logic is so simple (< 20 lines, no edge cases) that a dependency would be overkill
   - You have discussed it with the team and documented the decision

### When you add a new library

```bash
# Inside the api container
docker exec requiem-dev-api-1 go get github.com/some/library@latest
docker exec requiem-dev-api-1 go mod tidy
```

Commit both `go.mod` and `go.sum`. Never edit them by hand.

---

## Step 1 — Write the Go Code

Create four files in this order (each builds on the previous).

### 1a. `type.go` — Define Your Types

Every response struct **must** implement the `IsData()` marker interface — `httpx.JSON` requires it.

> **Rule: always use `snake_case` for JSON field names.**
> Every `json:"..."` tag in this codebase uses lower_snake_case. Never use camelCase or PascalCase.
> ✅ `json:"has_profanity"` `json:"flagged_words"` `json:"browser_version"`
> ❌ `json:"hasProfanity"` `json:"flaggedWords"` `json:"browserVersion"`

```go
package riddle

// Request for POST endpoints with JSON body.
// Use validate tags for automatic validation.
type GenerateRequest struct {
    Category string `json:"category" validate:"required,oneof=general science history"`
}

// Response types — every one needs IsData().
type Riddle struct {
    ID       int    `json:"id"`
    Question string `json:"question"`
    Answer   string `json:"answer"`
    Category string `json:"category"`
}

func (Riddle) IsData() {}

// For collections or richer responses:
type RiddleList struct {
    Items []Riddle `json:"items"`
    Total int      `json:"total"`
}

func (RiddleList) IsData() {}
```

**Validation tag reference** (from `go-playground/validator`):

| Tag | Meaning |
|---|---|
| `required` | Field must be present and non-zero |
| `email` | Must be a valid email address |
| `oneof=a b c` | Must be one of the listed values |
| `min=1,max=100` | Numeric range (or string/slice length) |
| `dive,email` | Validate each element of a slice |
| `url` | Must be a valid URL |

> **Rule: always declare validation in struct tags — never write it by hand in handlers.**
> `validate` tags work on both `json` (POST body via `httpx.Handle`) and `query` (GET params via `httpx.BindQuery`) structs.
> The framework enforces them automatically and returns consistent error responses.
>
> ✅ Correct — validation declared in the struct:
> ```go
> type ParseRequest struct {
>     UA string `query:"ua" validate:"required"`
> }
> ```
>
> ❌ Wrong — manual inline check that bypasses the validation pipeline:
> ```go
> ua := r.URL.Query().Get("ua")
> if ua == "" {
>     httpx.Error(w, http.StatusBadRequest, "bad_request", "ua parameter is required")
>     return
> }
> ```

### 1b. `service.go` — Business Logic

Services receive dependencies through their constructor. Keep this layer free of HTTP concerns.

```go
package riddle

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
    db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
    return &Service{db: db}
}

func (s *Service) Random(ctx context.Context, category string) (Riddle, error) {
    row := s.db.QueryRow(ctx, `
        SELECT id, question, answer, category
        FROM riddles
        WHERE category = $1
        ORDER BY random()
        LIMIT 1`, category)

    var r Riddle
    if err := row.Scan(&r.ID, &r.Question, &r.Answer, &r.Category); err != nil {
        return Riddle{}, err
    }
    return r, nil
}
```

For services with **no database dependency** (e.g., in-memory computation, embedded data):

```go
type Service struct{}

func NewService() *Service { return &Service{} }
```

### 1c. `transport_http.go` — HTTP Handlers

Choose the right pattern based on the endpoint's input method.

---

#### Pattern A: `httpx.Handle` — POST with JSON body (recommended for mutations)

Use this when the request has a JSON body. The wrapper handles decoding, validation, and `AppError` unwrapping automatically.

```go
package riddle

import (
    "context"
    "net/http"

    "github.com/go-chi/chi/v5"
    "your-module/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
    r.Post("/riddle/generate", httpx.Handle(
        func(ctx context.Context, req GenerateRequest) (Riddle, error) {
            return svc.Random(ctx, req.Category)
        },
    ))
}
```

`httpx.Handle` automatically:
- Caps the request body at 1 MB
- Decodes JSON (strict mode — unknown fields are rejected)
- Runs struct validation; returns `422 Unprocessable Entity` with field-level errors on failure
- Returns `500 Internal Server Error` for unexpected errors
- Unwraps `*httpx.AppError` and responds with the specified status/code

---

#### Pattern B: `httpx.BindQuery` — GET with query parameters

Use this when input comes from query string parameters. Set defaults **before** calling `BindQuery`. Always declare required fields and constraints using `validate` tags on the struct — do not add manual checks in the handler body.

```go
func RegisterRoutes(r chi.Router, svc *Service) {
    r.Get("/riddles", func(w http.ResponseWriter, r *http.Request) {
        // Always set defaults before binding.
        req := ListRequest{Page: 1, PerPage: 20}

        if err := httpx.BindQuery(r, &req); err != nil {
            httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
            return
        }

        result, err := svc.List(r.Context(), req.Page, req.PerPage)
        if err != nil {
            httpx.Error(w, http.StatusInternalServerError, "internal_error", "failed to fetch riddles")
            return
        }

        httpx.JSON(w, http.StatusOK, result)
    })
}
```

---

#### Pattern C: `chi.URLParam` — GET with URL path parameters

Use this when the identifier is part of the path (e.g., `/riddles/{id}`).

```go
func RegisterRoutes(r chi.Router, svc *Service) {
    r.Get("/riddles/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        if id == "" {
            httpx.Error(w, http.StatusBadRequest, "bad_request", "id is required")
            return
        }

        riddle, err := svc.GetByID(r.Context(), id)
        if err != nil {
            httpx.Error(w, http.StatusNotFound, "not_found", "riddle not found")
            return
        }

        httpx.JSON(w, http.StatusOK, riddle)
    })
}
```

---

**Error response reference:**

| Situation | Status | Code (snake_case) |
|---|---|---|
| Missing/invalid JSON body | 400 | `bad_request` |
| Failed struct validation | 422 | `validation_failed` (automatic via Handle) |
| Resource not found | 404 | `not_found` |
| Caller not authorised | 403 | `forbidden` |
| Upstream/DB unavailable | 503 | `service_unavailable` |
| Unexpected failure | 500 | `internal_error` |

### 1d. `router.go` — Wire the Domain

This file instantiates the service and calls the feature's `RegisterRoutes`.

```go
package text  // parent domain package

import (
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "your-module/internal/text/riddle"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool) {
    // ... existing features
    svc := riddle.NewService(pool)
    riddle.RegisterRoutes(r, svc)
}
```

---

## Step 2 — Mount the Router

**Adding to an existing domain** (most common): just add the lines above to the existing `router.go` for that domain. No other files need changing.

**Creating a brand-new top-level domain**: add two lines to `apps/api/internal/app/routes_v1.go`:

```go
func registerV1Routes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client) {
    // ... existing mounts

    puzzlesRouter := chi.NewRouter()
    puzzles.RegisterRoutes(puzzlesRouter, pool)
    r.Mount("/puzzles", puzzlesRouter)
}
```

---

## Step 3 — Database Migrations (if needed)

If the feature needs new tables or columns, create a SQL migration file:

```
infra/migrations/YYYYMMDDHHMMSS_add_riddles_table.up.sql
```

Example:
```sql
CREATE TABLE riddles (
    id       SERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer   TEXT NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'general'
);
```

**No Go registration is needed.** The app calls `db.MigrateWithRetry()` on startup and discovers all `*.up.sql` files in that directory automatically.

Do **not** create a `.down.sql` file unless you have a tested rollback procedure — partial rollbacks are worse than no rollback.

---

## Step 4 — Tests

Tests are **required before merge**. There are two layers to cover.

### Service Unit Tests (`service_test.go`)

Use table-driven tests. Keep each case small and named clearly.

```go
package riddle

import (
    "testing"
)

func TestService_Random(t *testing.T) {
    // For services with no DB: svc := NewService()
    // For DB-backed services: use a test DB or interface/mock.

    tests := []struct {
        name     string
        category string
        wantErr  bool
    }{
        {name: "valid category", category: "general", wantErr: false},
        {name: "empty category", category: "", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // assertion logic here
        })
    }
}
```

### HTTP Handler Tests (`transport_http_test.go`)

Test the full HTTP layer using `httptest`. Always include sad paths.

```go
package riddle

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/go-chi/chi/v5"
    "your-module/internal/platform/httpx"
)

func setupRouter() chi.Router {
    r := chi.NewRouter()
    svc := NewService() // inject test deps here if needed
    RegisterRoutes(r, svc)
    return r
}

func TestRiddle_Generate_HappyPath(t *testing.T) {
    r := setupRouter()

    body := `{"category":"general"}`
    req := httptest.NewRequest(http.MethodPost, "/riddle/generate", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    var resp httpx.Response[Riddle]
    if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
        t.Fatalf("failed to decode: %v", err)
    }

    if resp.Data.Question == "" {
        t.Error("expected a non-empty question")
    }
}

func TestRiddle_Generate_MissingCategory(t *testing.T) {
    r := setupRouter()

    req := httptest.NewRequest(http.MethodPost, "/riddle/generate", strings.NewReader(`{}`))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    r.ServeHTTP(w, req)

    if w.Code != http.StatusUnprocessableEntity {
        t.Errorf("expected 422, got %d", w.Code)
    }
}
```

**Run tests:**

```bash
# Domain-scoped (fast feedback during development)
docker exec requiem-dev-api-1 go test ./internal/text/riddle/...

# Full suite (required before pushing)
docker exec requiem-dev-api-1 go test ./...

# With race detection and coverage (CI equivalent)
docker exec requiem-dev-api-1 go test -race -coverprofile=coverage.out ./...
```

---

## Step 5 — Update the Dashboard API Catalog

> Skip this step if you are adding a new endpoint to an **existing** API (the catalog entry already exists).

File: `apps/dashboard/config/api_catalog.yml`

Add an entry under the appropriate category:

```yaml
- id: riddles
  name: Riddles
  category: text
  description: Generate random riddles across categories like general, science, and history.
  endpoints: 2
  status: live          # or coming_soon
  popular: false
  documentation_url: /apis/riddles
  tags:
    - riddles
    - trivia
    - fun
```

The `documentation_url` must match the route the Rails app will serve for the API detail page.

---

## Step 6 — Add the API Documentation YAML

File: `apps/dashboard/config/api_docs/riddles.yml`

This YAML powers the interactive API documentation page in the dashboard. Every field is rendered in the UI — do not leave any required section empty.

```yaml
api_id: riddles
api_name: Riddles
description: Generate random riddles across multiple categories for trivia apps, brain teasers, and educational content.
base_url: https://api.requiems.xyz

overview:
  use_cases:
    - Trivia apps and quiz games
    - Educational platforms
    - Ice-breaker bots
    - Daily challenge widgets

  features:
    - Multiple categories (general, science, history)
    - Deterministic question/answer pairs
    - Fast, low-latency responses

endpoints:
  - name: Generate Random Riddle
    method: POST
    path: /v1/text/riddle/generate
    description: Returns a random riddle from the specified category.

    parameters:
      - name: category
        type: string
        required: true
        location: body
        description: Category of riddle to return.
        example: general

    request_example: |
      {
        "category": "general"
      }

    response_example: |
      {
        "data": {
          "id": 42,
          "question": "What has keys but no locks?",
          "answer": "A keyboard",
          "category": "general"
        },
        "metadata": {
          "timestamp": "2026-01-01T00:00:00Z"
        }
      }

    response_fields:
      - name: id
        type: integer
        description: Unique riddle identifier
      - name: question
        type: string
        description: The riddle question
      - name: answer
        type: string
        description: The answer to the riddle
      - name: category
        type: string
        description: Category the riddle belongs to

    errors:
      - code: validation_failed
        status: 422
        description: The category field is missing or contains an invalid value.
      - code: internal_error
        status: 500
        description: Unexpected server error.

    code_examples:
      curl: |
        curl -X POST https://api.requiems.xyz/v1/text/riddle/generate \
          -H "requiems-api-key: YOUR_API_KEY" \
          -H "Content-Type: application/json" \
          -d '{"category": "general"}'

      python: |
        import requests

        url = "https://api.requiems.xyz/v1/text/riddle/generate"
        headers = {
            "requiems-api-key": "YOUR_API_KEY",
            "Content-Type": "application/json"
        }
        payload = {"category": "general"}

        response = requests.post(url, headers=headers, json=payload)
        print(response.json())

      javascript: |
        const response = await fetch('https://api.requiems.xyz/v1/text/riddle/generate', {
          method: 'POST',
          headers: {
            'requiems-api-key': 'YOUR_API_KEY',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ category: 'general' })
        });

        const data = await response.json();
        console.log(data.data.question);

      ruby: |
        require 'net/http'
        require 'json'

        uri = URI('https://api.requiems.xyz/v1/text/riddle/generate')
        request = Net::HTTP::Post.new(uri)
        request['requiems-api-key'] = 'YOUR_API_KEY'
        request['Content-Type'] = 'application/json'
        request.body = { category: 'general' }.to_json

        response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
          http.request(request)
        end

        data = JSON.parse(response.body)
        puts data['data']['question']

pricing:
  free_tier: Included in all plans
  rate_limits:
    - plan: Free
      limit: 30 requests/minute
    - plan: Developer
      limit: 5,000 requests/minute
    - plan: Business
      limit: 10,000 requests/minute
    - plan: Professional
      limit: 50,000 requests/minute

faq:
  - question: Can I request riddles from multiple categories at once?
    answer: Not currently — each request returns one riddle from one category. Batch support is planned.

  - question: Are riddle IDs stable across requests?
    answer: Yes, IDs are stable database identifiers. You can use them to avoid showing a riddle twice.
```

Use `apps/dashboard/config/api_docs/disposable-email.yml` as the reference template for complex multi-endpoint APIs and `advice.yml` for simple single-endpoint APIs.

---

## Step 7 — Add Markdown API Docs

File: `docs/apis/{category}/{api-name}.md`

This is the developer-facing plain-text companion in the repo. It should explain the *why* and *edge cases* that the YAML cannot express clearly.

```markdown
# Riddles API

Returns random riddles across multiple categories.

## Endpoint

`POST /v1/text/riddle/generate`

## Categories

| Value | Description |
|---|---|
| `general` | Everyday wordplay and logic |
| `science` | Biology, physics, chemistry |
| `history` | Historical facts framed as riddles |

## Response Envelope

All responses are wrapped in the standard envelope:

\`\`\`json
{
  "data": { ... },
  "metadata": { "timestamp": "2026-01-01T00:00:00Z" }
}
\`\`\`

## Error Codes

| Code | Status | When |
|---|---|---|
| `validation_failed` | 422 | Invalid or missing category |
| `internal_error` | 500 | Unexpected failure |
```

---

## Step 8 — Update the Credit Multiplier (if non-default)

File: `apps/workers/shared/src/config.ts`

Most endpoints cost **1 credit** per request and do not need a config entry. Only add an entry if your endpoint is more expensive to run:

```ts
// In ENDPOINT_MULTIPLIERS map:
["/v1/text/riddle/generate", 1],   // default — omit this line unless non-1
["/v1/ai/summarize", 5],           // expensive AI call
["/v1/translate/text", 3],         // ML translation
```

The gateway uses `getRequestMultiplier(method, pathname)` to look up the multiplier before deducting credits. If no entry exists, it defaults to `1`.

---

## Docker Considerations

**No Docker changes are needed** for adding Go endpoints. The dev container (`requiem-dev-api-1`) runs with Air hot reload — new `.go` files are compiled and reloaded automatically when saved.

If your feature adds a **new external service dependency** (e.g., a new third-party HTTP client, a new Redis data structure), check `infra/docker/docker-compose.dev.yml` for any environment variables or service additions required, but this is uncommon.

---

## Pre-Merge Verification Checklist

Work through these in order before opening a PR.

```
Go code
  [ ] go test ./... passes in the container (zero failures)
  [ ] go test -race -coverprofile=coverage.out ./... passes (no races)
  [ ] golangci-lint run passes or only advisory warnings remain

Manual smoke test
  [ ] curl -X POST http://localhost:8080/v1/{domain}/{endpoint} \
        -H "X-Backend-Secret: your_local_secret" \
        -H "Content-Type: application/json" \
        -d '{...}' returns expected response
  [ ] Invalid input returns correct 4xx with descriptive error

Database (if applicable)
  [ ] Migration file follows naming convention: YYYYMMDDHHMMSS_description.up.sql
  [ ] App starts cleanly after migration (no startup errors in docker logs api)

Documentation
  [ ] apps/dashboard/config/api_docs/{name}.yml created with all sections
  [ ] apps/dashboard/config/api_catalog.yml updated (new API only)
  [ ] docs/apis/{category}/{name}.md created

Workers
  [ ] apps/workers/shared/src/config.ts updated if credit cost != 1
  [ ] bun run typecheck passes in apps/workers/shared/ (if config.ts was touched)
```

---

## Worked Example: `GET /v1/text/riddle/random`

Below is a complete minimal implementation for a riddle endpoint backed by an in-memory list (no database).

**`apps/api/internal/text/riddle/type.go`**

```go
package riddle

type Riddle struct {
    Question string `json:"question"`
    Answer   string `json:"answer"`
}

func (Riddle) IsData() {}
```

**`apps/api/internal/text/riddle/service.go`**

```go
package riddle

import "math/rand"

var riddles = []Riddle{
    {Question: "What has keys but no locks?", Answer: "A keyboard"},
    {Question: "What gets wetter the more it dries?", Answer: "A towel"},
}

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Random() Riddle {
    return riddles[rand.Intn(len(riddles))]
}
```

**`apps/api/internal/text/riddle/transport_http.go`**

```go
package riddle

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "your-module/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
    r.Get("/riddle/random", func(w http.ResponseWriter, r *http.Request) {
        httpx.JSON(w, http.StatusOK, svc.Random())
    })
}
```

**`apps/api/internal/text/riddle/transport_http_test.go`**

```go
package riddle

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "your-module/internal/platform/httpx"
)

func setupRouter() chi.Router {
    r := chi.NewRouter()
    RegisterRoutes(r, NewService())
    return r
}

func TestRiddle_Random(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/riddle/random", http.NoBody)
    w := httptest.NewRecorder()
    setupRouter().ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    var resp httpx.Response[Riddle]
    if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
        t.Fatalf("decode error: %v", err)
    }

    if resp.Data.Question == "" {
        t.Error("expected non-empty question")
    }
    if resp.Data.Answer == "" {
        t.Error("expected non-empty answer")
    }
}
```

**Add to `apps/api/internal/text/router.go`:**

```go
svc := riddle.NewService()
riddle.RegisterRoutes(r, svc)
```

That is everything needed to ship a working, tested, documented endpoint.
