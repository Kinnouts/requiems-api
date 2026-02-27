## Getting Started with Requiem API

---

## Local Development Setup

### Prerequisites

- Docker & Docker Compose

### 1. Start the stack

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

Services available:

| Service        | URL                              | Notes                  |
| -------------- | -------------------------------- | ---------------------- |
| Rails Dashboard | http://localhost:3000            | Hot reload             |
| Go API          | http://localhost:8080/healthz    | Air hot reload         |
| Auth Gateway    | http://localhost:4455/healthz    | Public API entry point |
| API Management  | http://localhost:5544/healthz    | Internal service       |
| PostgreSQL      | localhost:5432                   | user/pass: `requiem`   |
| Redis           | localhost:6379                   |                        |

### 2. Environment variables

All services load `infra/docker/.env.example` by default — **no extra setup
required** for most development work.

For features that require real credentials (payments, email sending), create
`infra/docker/.env.local` (gitignored) and add only the secrets you need:

```bash
# infra/docker/.env.local — DO NOT COMMIT
LEMONSQUEEZY_API_KEY=        # required for testing payment flows
LEMONSQUEEZY_SIGNING_SECRET= # required for webhook verification
SMTP_PASSWORD=               # required for sending real emails
```

Ask a teammate for the actual values. The app starts fine without `.env.local` —
payments and email delivery just won't work end-to-end locally.

> **Note:** Emails in development are logged to stdout (not sent). Check
> `docker compose logs dashboard` to see email content.

---

### 1. Sign up and get an API key

- Visit the Requiem API dashboard,
- Create an account and generate an API key.
- All requests must include your key in the `requiems-api-key` header.

### 2. Make your first request

Example using `curl` against the edge (replace `YOUR_API_KEY` and domain when
live):

```bash
curl https://api.requiems.xyz/v1/text/advice \
  -H "requiems-api-key: YOUR_API_KEY"
```

Development (local Go backend via Docker Compose):

```bash
# Health check (no auth required)
curl http://localhost:8080/healthz

# API endpoints
curl http://localhost:8080/v1/text/advice
curl http://localhost:8080/v1/text/quotes/random
curl http://localhost:8080/v1/text/words/random
```

> **Note:** The local Go backend (port 8080) does not enforce auth. In
> production, the Cloudflare Worker handles auth before forwarding to the
> backend.
