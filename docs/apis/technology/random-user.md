# Random User API

## Status

✅ **Live** - Production-ready

## Overview

Generate random fake user profiles for testing and prototyping. Each call
returns a unique name, email address, phone number, mailing address, and avatar
URL.

## Base URL

All endpoints are mounted under `/v1/technology`

## Endpoints

### Get Random User

Returns a randomly generated fake user profile.

**Endpoint:** `GET /v1/technology/random-user`

**Response:** `200 OK`

```json
{
  "name": "Grace Lopez",
  "email": "grace.lopez@example.org",
  "phone": "555-123-4567",
  "address": {
    "street": "4821 Maple Avenue",
    "city": "North Judyton",
    "state": "California",
    "zip": "94103",
    "country": "United States of America"
  },
  "avatar": "https://api.dicebear.com/9.x/identicon/svg?seed=Grace+Lopez"
}
```

**Example:**

```bash
curl https://api.requiems.xyz/v1/technology/random-user \
  -H "requiems-api-key: YOUR_API_KEY"
```

---

## Use Cases

- Populate test databases with realistic-looking user records
- Demo and prototype applications without real user data
- Load testing with varied user payloads
- UI mockups and design previews

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 875 ms  |
| p95     | 1016 ms |
| p99     | 1057 ms |
| Average | 884 ms  |

_Last updated: 2026-04-16_
