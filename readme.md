<p align="center">
  <p align="center">
    <a href="https://requiems.xyz/?utm_source=github&utm_medium=logo" target="_blank">
      <img src="https://raw.githubusercontent.com/bobadilla-tech/requiems-api/refs/heads/main/apps/dashboard/app/assets/images/logo.png" alt="Requiems API" width="280" />
    </a>
  </p>
  <p align="center">
    One API key. Multiple enterprise-grade APIs.
  </p>
  <p align="center">
    <i>A product by <a href="https://bobadilla.tech">Bobadilla Technologies</a></i>
  </p>
</p>

# What's Requiems API?

Requiems API is a production-ready API platform providing unified access to multiple enterprise-grade APIs. Eliminate months of data sourcing, validation logic, and infrastructure setup. Start building features today, not infrastructure.

[![CI](https://github.com/bobadilla-tech/requiems-api/actions/workflows/ci.yml/badge.svg)](https://github.com/bobadilla-tech/requiems-api/actions/workflows/ci.yml)
[![Get Started](https://img.shields.io/badge/Get_Started-→-blue)](https://requiems.xyz)
[![Documentation](https://img.shields.io/badge/Documentation-📖-green)](https://requiems.xyz/apis)

## Built for Scale, Designed for Speed

- **Go API** – Lightning-fast backend with domain-driven design
- **Rails Dashboard** – Beautiful UI for users and admins
- **Cloudflare Worker Gateway** – Global edge network handling auth, rate limiting, and credit tracking

## Quick Start

Get your API key at [requiems.xyz](https://requiems.xyz), then try it out:

```bash
# Example: Check if an email is disposable (one of our many APIs)
curl -X POST https://api.requiems.xyz/v1/email/disposable/check \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@tempmail.com"}'
```

Explore our full catalog of APIs including email validation, finances utilities, and more in the [documentation](https://requiems.xyz/apis).
