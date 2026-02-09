## Getting Started with Requiem API

### 1. Sign up and get an API key

- Visit the Requiem API dashboard,
- Create an account and generate an API key.
- All requests must include your key in the `requiems-api-key` header.

### 2. Make your first request

Example using `curl` against the edge (replace `YOUR_API_KEY` and domain when
live):

```bash
curl https://api.requiems.xyz/v1/text/advice \
  -H "x-api-key: YOUR_API_KEY"
```

Development (local Go backend via Docker Compose):

```bash
# Health check (no auth required)
curl http://localhost:6969/healthz

# API endpoints
curl http://localhost:6969/v1/text/advice
curl http://localhost:6969/v1/text/quotes/random
curl http://localhost:6969/v1/text/words/random
```

> **Note:** The local Go backend (port 6969) does not enforce auth. In
> production, the Cloudflare Worker handles auth before forwarding to the
> backend.
