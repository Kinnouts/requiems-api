# ⚰️ Requiem API

Requiem API is a **managed API service** that provides a single API key to access a growing collection of production-ready APIs.

We build, operate, and scale the APIs. Customers only consume them.

---

## What Requiem API Is

* A unified API platform
* One API key for many APIs
* Managed infrastructure, scaling, and reliability
* Usage-based billing with free tiers
* Built for developers and product teams

## What Requiem API Is Not

* Not an API framework
* Not a marketplace for third-party APIs
* Not a no-code or low-code builder

---

## High-Level Architecture

```
Client
  │
  ▼
Edge Auth Gateway (Cloudflare Workers)
- API key validation
- Rate limits & quotas
- Abuse protection
- Request signing
  │
  ▼
Backend API (Go)
- API routing
- Usage accounting
- Billing integration
- Job orchestration
  │
  ▼
Internal Workers (Go)
- Async processing
- Heavy computation
- Data pipelines
```

All APIs are owned, maintained, and operated by Requiem API.

---

## Core Components

### Edge Auth Gateway

Runs at the edge using Cloudflare Workers.

Responsibilities:

* Validate API keys
* Enforce rate limits and plans
* Block abuse early
* Forward trusted requests to the backend

The edge never runs business logic or long-running tasks.

---

### Backend API (Go)

A single monolithic Go service.

Responsibilities:

* Expose all public APIs
* Track usage per customer
* Integrate payments and plans
* Generate signed upload URLs when required
* Orchestrate asynchronous jobs

All backend traffic comes through the edge gateway.

---

### Workers (Go)

Internal asynchronous workers.

Responsibilities:

* Execute heavy or slow tasks
* Process files and data
* Fetch external sources
* Update job state

Workers are private and never exposed publicly.

---

### Documentation & Dashboard

Responsibilities:

* API documentation
* API key management
* Usage and billing visibility
* Examples and SDK references

Initially focused on documentation.

---

## Typical Request Flow

1. Client sends a request with an API key
2. Edge gateway validates plan and limits
3. Backend API processes the request
4. Optional async job is created
5. Worker executes the job
6. Client polls status or receives the result

Files are uploaded directly to object storage using signed URLs. The API never receives raw files.

---

## Technology Stack

* Edge: Cloudflare Workers (TypeScript)
* Backend API: Go
* Workers: Go
* Queue: Redis or managed queue service
* Storage: S3-compatible object storage
* Payments: Lemon Squeezy
* Infrastructure: Docker, Caddy, VPS

---

## Repository Structure

```
apps/
  edge-auth/        Edge authentication gateway
  api/              Backend API (Go)
  worker/           Async workers (Go)
  web/              Docs and dashboard

infra/
  docker/           Dockerfiles and compose
  caddy/            HTTPS and routing
  scripts/          Deployment helpers
```

---

## Deployment Philosophy

* Start on a single VPS
* Everything runs in Docker
* HTTPS via Caddy
* Edge deployed independently
* Scale workers horizontally
* Migrate to managed services when needed

---

## Open Source Policy

This repository is source-available.

The code is public to encourage transparency, feedback, and contributions.

The hosted service, APIs, and brand are operated exclusively by Requiem API.

---

## Project Status

Early stage. APIs will be added incrementally. Breaking changes may occur.

---

## Philosophy

APIs are leverage.
Infrastructure is power.
Distribution decides the winner.

Requiem API focuses on owning the boring, hard parts so customers don’t have to.
