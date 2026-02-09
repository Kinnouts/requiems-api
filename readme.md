## ⚰️ Requiem API

Requiem API is a **managed API platform** that gives you a single API key to
access a growing collection of production-ready APIs.

We build, operate, and scale the APIs. You focus on shipping product.

---

## Quick Links

| Resource         | URL                                                                            |
| ---------------- | ------------------------------------------------------------------------------ |
| 🌐 Website       | [requiems.xyz](https://requiems.xyz)                                           |
| 📚 Documentation | [requiems.xyz/docs](https://requiems.xyz/docs)                                 |
| 🎮 Dashboard     | [requiems.xyz/dashboard](https://requiems.xyz/dashboard)                       |
| 🔗 API Base URL  | `https://api.requiems.xyz`                                                     |
| 💼 LinkedIn      | [Requiems API](https://www.linkedin.com/showcase/requiems-api/)                |

---

- **One API key** for many APIs
- **Managed infrastructure**, scaling, and reliability
- **Tier-based billing** with a generous free tier
- **Web playground** to test endpoints before you buy

### Getting Started

1. Visit [requiems.xyz](https://requiems.xyz)
2. Sign up for a free account
3. Try endpoints in the [playground](https://requiems.xyz/playground)
4. Upgrade when you need more

### Example Request

```bash
curl -H "x-api-key: YOUR_KEY" https://api.requiems.xyz/v1/text/advice
```

---

## 🏗️ Architecture

This is a **multi-language monorepo** with clean separation of concerns:

```
apps/
├── api/              # Go backend (internal API)
├── dashboard/        # Rails 8 (landing, dashboard, admin)
└── edge-auth/        # Cloudflare Worker (auth gateway)
```

### URL Structure

| URL | Purpose | Technology |
|-----|---------|------------|
| `requiems.xyz` | Landing page + Dashboard + Admin | Rails 8.1 |
| `api.requiems.xyz` | Public API gateway (auth, rate limiting) | Cloudflare Worker |
| `internal.requiems.xyz` | Internal backend (business logic) | Go 1.23 |

### Request Flow

```
Client
  ↓ x-api-key: rq_live_xxx
Cloudflare Worker (api.requiems.xyz)
  ↓ validate key, check limits
  ↓ X-Backend-Secret: xxx
Go Backend (internal.requiems.xyz)
  ↓ business logic
PostgreSQL
```

## 🚀 Development

### Prerequisites

- Go 1.23+
- Ruby 3.4+
- PostgreSQL 16
- Redis 7
- Docker & Docker Compose

### Quick Start

```bash
# Clone the repository
git clone https://github.com/bobadilla-tech/requiems-api.git
cd requiems-api

# Start all services with Docker Compose
cd infra/docker
docker compose up --build

# Services will be available at:
# - Rails dashboard: http://localhost:3000
# - Go API: http://localhost:6969
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379
```

### Local Development (without Docker)

**Go Backend:**
```bash
cd apps/api
go mod download
go run main.go
# Runs on http://localhost:8080
```

**Rails Dashboard:**
```bash
cd apps/dashboard
bundle install
rails db:create db:migrate
rails server
# Runs on http://localhost:3000
```

**Cloudflare Worker:**
```bash
cd apps/edge-auth
npm install
npm run dev
# Runs on http://localhost:8787
```

### Environment Variables

Create `.env` files in each app directory:

**apps/api/.env:**
```env
PORT=8080
DATABASE_URL=postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable
```

**apps/dashboard/.env:**
```env
DATABASE_URL=postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable
REDIS_URL=redis://localhost:6379
SECRET_KEY_BASE=your_secret_key_base
CLOUDFLARE_ACCOUNT_ID=your_account_id
CLOUDFLARE_KV_NAMESPACE_ID=your_namespace_id
CLOUDFLARE_API_TOKEN=your_api_token
```

**apps/edge-auth/.env:**
```env
BACKEND_URL=http://localhost:8080
BACKEND_SECRET=your_backend_secret
```

## 📁 Repository Structure

```
requiems-api/
├── apps/
│   ├── api/                    # Go backend
│   │   ├── internal/          # Domain-driven design
│   │   │   ├── app/          # Application setup
│   │   │   ├── email/        # Email APIs
│   │   │   ├── text/         # Text APIs
│   │   │   └── platform/     # Shared utilities
│   │   ├── infra/migrations/ # Go database migrations
│   │   ├── go.mod
│   │   └── main.go
│   ├── dashboard/             # Rails 8 dashboard
│   │   ├── app/
│   │   │   ├── controllers/
│   │   │   │   ├── dashboard/  # User dashboard
│   │   │   │   └── admin/      # Admin panel
│   │   │   ├── models/        # User, ApiKey, Subscription, etc.
│   │   │   ├── services/      # Cloudflare KV sync
│   │   │   └── views/
│   │   ├── db/migrate/        # Rails migrations
│   │   ├── Gemfile
│   │   └── config/
│   └── edge-auth/             # Cloudflare Worker
│       ├── src/
│       │   ├── index.ts      # Main handler
│       │   ├── rate-limit.ts # Rate limiting (KV)
│       │   └── credits.ts    # Usage tracking (D1)
│       └── wrangler.toml
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   ├── api.Dockerfile
│   │   └── dashboard.Dockerfile
│   └── caddy/
│       └── Caddyfile         # Reverse proxy config
├── docs/                      # Documentation
└── readme.md
```

## 🗄️ Database Architecture

### Shared PostgreSQL

Both Go and Rails use the same PostgreSQL database (`requiem`):

- **Go manages:** `advice`, `quotes`, `words` (business data)
- **Rails manages:** `users`, `api_keys`, `subscriptions`, `usage_logs` (user data)
- **Migration tracking:** Separate tables (`schema_migrations` for Go, `rails_schema_migrations` for Rails)

### Cloudflare Edge Storage

- **KV Store:** API key lookup (<10ms latency)
- **D1 SQLite:** Usage tracking and aggregations

---

## 📝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
