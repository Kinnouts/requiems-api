# Email APIs

## Overview

Email validation and services endpoints providing tools for email verification
and management.

## Endpoints

### [Disposable Email](./disposable.md) - ✅ MVP

Check if an email address or domain is from a disposable email service

- **Status:** mvp
- **Endpoints:**
  - `POST /v1/email/disposable/check` — Check single email
  - `POST /v1/email/disposable/check-batch` — Check up to 100 emails
  - `GET /v1/email/disposable/domain/{domain}` — Check domain
  - `GET /v1/email/disposable/domains` — List all disposable domains (paginated)
  - `GET /v1/email/disposable/stats` — Blocklist statistics
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 1
- Live: 1
- Planned: 0
