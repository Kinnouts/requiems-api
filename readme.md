## ⚰️ Requiem API

### Ship Faster With Production-Ready APIs

**One API key. Multiple enterprise-grade APIs.** Eliminate months of data
sourcing, validation logic, and infrastructure setup. Start building features
today, not infrastructure.

**[Get Started →](https://requiems.xyz)** ·
**[Documentation](https://requiems.xyz/docs)** ·
**[Try Live Playground](https://requiems.xyz/playground)**

---

### Stop Building Infrastructure. Start Shipping Features

**Ditch the integration hell:**

- ❌ No more juggling 10+ API providers and payment accounts
- ❌ No more building validation logic from scratch
- ❌ No more maintaining rate limiters, caching, and retry logic
- ❌ No more infrastructure babysitting

**Ship production features in minutes:**

- ✅ **One unified API** – email validation, disposable detection, text
  utilities, and growing
- ✅ **Battle-tested infrastructure** – sub-200ms response times, 99.9% uptime
- ✅ **Zero setup friction** – generous free tier, no credit card required
- ✅ **Test before you commit** – live playground with real responses

### From Zero to Production in 60 Seconds

```bash
# 1. Get your free API key at requiems.xyz (no credit card)

# 2. Make your first call
curl -H "requiems-api-key: YOUR_KEY" https://api.requiems.xyz/v1/email/disposable/check \
  -d '{"email":"test@tempmail.com"}'

# 3. Ship it to production ✅
```

**That's it.** No SDKs to install. No complex authentication flows. No
infrastructure to maintain.

---

## 🏗️ Built for Scale, Designed for Speed

- **Go API** – Lightning-fast backend with domain-driven design
- **Rails Dashboard** – Beautiful UI for users and admins
- **Cloudflare Worker Gateway** – Global edge network handling auth, rate
  limiting, and credit tracking

## 🚀 Local Development in One Command

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

## 🤝 Contributing

**Found a bug? Want to add an API? We'd love your help.**

We welcome contributions of all sizes: bug fixes, performance improvements,
documentation, or entire new API endpoints.

→ [Contributing guidelines](./contributing.md)
