## ⚰️ Requiem API

A **managed API platform** that gives you one API key to access a growing collection of production-ready APIs. We handle the infrastructure, you ship faster.

**[Get Started →](https://requiems.xyz)** · **[Documentation](https://requiems.xyz/docs)** · **[Playground](https://requiems.xyz/playground)**

[![LinkedIn](https://img.shields.io/badge/LinkedIn-Requiems%20API-0077B5?logo=linkedin)](https://www.linkedin.com/showcase/requiems-api/)

---

### Why Requiems?

- **One API key** for many APIs – no juggling multiple accounts
- **Managed infrastructure** – scaling, monitoring, and uptime handled for you
- **Generous free tier** – start building without a credit card
- **Test in the playground** – try endpoints before committing

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

Multi-language monorepo with three main apps:

- **Dashboard** (Rails 8) – Landing page, user dashboard, admin panel
- **API** (Go 1.23) – Core backend with business logic
- **Auth Gateway** (Cloudflare Worker) – Authentication, rate limiting, credit tracking

Requests flow through the Cloudflare Worker for authentication, then route to the Go backend for processing.

## 🚀 Development

Want to run this locally or contribute? We've made it easy.

### Quick Start with Docker

Everything runs with hot reload out of the box:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

All services start automatically. Edit any code and see changes instantly.

See [infra/docker/README.md](infra/docker/README.md) for more details.

## 📁 Repository Structure

```
apps/
├── api/           # Go backend (domain-driven design)
├── dashboard/     # Rails 8 dashboard + admin
└── edge-auth/     # Cloudflare Worker (auth gateway)

infra/
├── docker/        # Docker Compose setup
└── caddy/         # Reverse proxy config
```

See full directory structure in [docs/](docs/).

---

## 📝 Contributing

We welcome contributions! Whether it's bug fixes, performance improvements, documentation, or new features, we'd love your help.

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to get started.

---

**Questions?** Open an issue or reach out on [LinkedIn](https://www.linkedin.com/showcase/requiems-api/). We're here to help.
