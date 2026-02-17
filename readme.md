# Requiems API

[![CI](https://github.com/bobadilla-tech/requiems-api/actions/workflows/ci.yml/badge.svg)](https://github.com/bobadilla-tech/requiems-api/actions/workflows/ci.yml)
[![Get Started](https://img.shields.io/badge/Get_Started-→-blue)](https://requiems.xyz)
[![Documentation](https://img.shields.io/badge/Documentation-📖-green)](https://requiems.xyz/apis)

**One API key. Multiple enterprise-grade APIs.** Eliminate months of data
sourcing, validation logic, and infrastructure setup. Start building features
today, not infrastructure.

## Built for Scale, Designed for Speed

- **Go API** – Lightning-fast backend with domain-driven design
- **Rails Dashboard** – Beautiful UI for users and admins
- **Cloudflare Worker Gateway** – Global edge network handling auth, rate
  limiting, and credit tracking

## Local Development in One Command

**Contributor-friendly setup. Hot reload everything.**

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

**That's it.** All services (API, Dashboard, Database, Caddy) start with hot
reload. Edit any file and see changes instantly.

[Full dev setup guide](./infra/docker/README.md) |
[Full API developer docs](./docs/)

---

## Contributing

**Found a bug? Want to add an API? We'd love your help.**

We welcome contributions of all sizes: bug fixes, performance improvements,
documentation, or entire new API endpoints.

→ [Contributing guidelines](./contributing.md)
