## Getting Started with Requiem API

### 1. Sign up and get an API key

- Visit the Requiem API dashboard (coming soon).
- Create an account and generate an API key.
- All requests must include your key in the `x-api-key` header.

### 2. Make your first request

Example using `curl` against the edge (replace `YOUR_API_KEY` and domain when
live):

```bash
curl https://api.requiems-api.xyz/v1/advice \
  -H "x-api-key: YOUR_API_KEY"
```

Development (local stack):

```bash
curl http://localhost:6969/v1/advice \
  -H "x-api-key: dev-key"  # if enforced by the Worker in local dev
```

### 3. Explore available APIs

- See `docs/api-catalogue.md` (or `apis.md` in the repo) for the full list of
  APIs and their status.
