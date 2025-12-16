# Credits & Pricing

This document explains the credit-based pricing model for Requiem API.

---

## Overview

Requiem API uses a **credit-based system** instead of simple request counts:

- Every API key has a **daily credit allowance**
- Different endpoints cost **different amounts of credits**
- More complex/expensive operations cost more credits
- Simple operations cost fewer credits

This lets us price fairly based on actual resource usage.

---

## Free Tier

| Plan | Daily Credits | Resets                    |
| ---- | ------------- | ------------------------- |
| Free | 50 credits    | Every 24h at midnight UTC |

Free tier is designed for:

- Testing and evaluation
- Hobby projects
- Low-volume personal use

> **Note:** Free credits do NOT accumulate. You get 50 fresh credits each
> day—use them or lose them.

---

## Paid Plans

| Plan       | Monthly Credits | Monthly Price | Per-Credit Overage |
| ---------- | --------------- | ------------- | ------------------ |
| Starter    | 30,000          | $9            | $0.001             |
| Pro        | 150,000         | $49           | $0.0008            |
| Business   | 500,000         | $299          | $0.0005            |
| Enterprise | Custom          | Custom        | Custom             |

### How Paid Credits Work

- **Monthly pool**: Use your credits whenever you want during the billing cycle
- **No daily limits**: Burst 30k requests in one day if you need to
- **Resets monthly**: Unused credits do NOT roll over to next month
- **Overage**: If you exceed your plan, overage is billed per-credit

### Why Monthly Instead of Daily?

Paid users need flexibility:

- Marketing launch? Burst traffic on day 1
- Batch processing? Run everything overnight
- Seasonal business? Heavy usage some weeks, light others

Rate limiting (second) still protects infrastructure from abuse.

### Overage Behavior

- **Free**: Hard limit, requests rejected after 50 daily credits
- **Paid**: Soft limit, overage billed at per-credit rate at end of month
- **Enterprise**: Custom limits and pricing

### Text Domain (`/v1/text/*`)

| Endpoint                               | Credits | Notes                |
| -------------------------------------- | ------- | -------------------- |
| `GET /v1/text/advice`                  | 1       | Simple random lookup |
| `GET /v1/text/quotes/random`           | 1       | Simple random lookup |
| `GET /v1/text/words/random`            | 1       | Simple random lookup |
| `GET /v1/text/words/define?word=...`   | 2       | Dictionary lookup    |
| `GET /v1/text/words/synonyms?word=...` | 2       | Thesaurus lookup     |

### Finance Domain (`/v1/finance/*`) — _Coming Soon_

| Endpoint                              | Credits | Notes                |
| ------------------------------------- | ------- | -------------------- |
| `GET /v1/finance/commodities/:symbol` | 3       | Real-time price data |
| `GET /v1/finance/stocks/:symbol`      | 3       | Real-time price data |
| `GET /v1/finance/crypto/:symbol`      | 3       | Real-time price data |
| `GET /v1/finance/exchange-rates`      | 2       | Currency conversion  |

### Places Domain (`/v1/places/*`) — _Coming Soon_

| Endpoint                         | Credits | Notes                 |
| -------------------------------- | ------- | --------------------- |
| `GET /v1/places/geocode`         | 5       | Address → coordinates |
| `GET /v1/places/reverse-geocode` | 5       | Coordinates → address |
| `GET /v1/places/timezone`        | 2       | Timezone lookup       |

---

## Credit Cost Factors

Credits are priced based on:

| Factor             | Impact                                   |
| ------------------ | ---------------------------------------- |
| **Compute**        | CPU/memory intensive = more credits      |
| **External APIs**  | If we pay upstream, you pay more credits |
| **Data freshness** | Real-time data costs more than cached    |
| **Response size**  | Large payloads cost more                 |
| **Complexity**     | AI/ML endpoints cost significantly more  |

## Credit Costs by Endpoint

### Response Headers

Every API response includes:

```
X-Credits-Used: 1
X-Credits-Remaining: 49
X-Credits-Reset: 2025-12-16T00:00:00Z
```

### Dashboard

View detailed usage at
[requiems-api.xyz/dashboard](https://requiems-api.xyz/dashboard):

- Real-time credit balance
- Usage by endpoint
- Historical usage charts
- Billing and invoices

### API Endpoint

```bash
curl -H "x-api-key: YOUR_KEY" https://api.requiems-api.xyz/v1/account/usage
```

```json
{
  "plan": "free",
  "credits": {
    "used": 12,
    "remaining": 38,
    "daily_limit": 50,
    "resets_at": "2025-12-16T00:00:00Z"
  }
}
```

---

## Rate Limiting

In addition to credits, we apply rate limits to prevent abuse:

| Plan       | Requests/second | Requests/minute |
| ---------- | --------------- | --------------- |
| Free       | 1               | 30              |
| Starter    | 10              | 300             |
| Pro        | 30              | 1,000           |
| Business   | 100             | 5,000           |
| Enterprise | 1,000           | 50,000          |

Rate limit headers:

```
X-RateLimit-Limit: 30
X-RateLimit-Remaining: 29
X-RateLimit-Reset: 1702684800
```

---

## Implementation Notes (Internal)

### Database Schema

```sql
CREATE TABLE api_keys (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  key_hash TEXT UNIQUE NOT NULL,
  plan TEXT DEFAULT 'free',
  daily_credit_limit INT DEFAULT 50,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE credit_usage (
  id UUID PRIMARY KEY,
  api_key_id UUID REFERENCES api_keys(id),
  endpoint TEXT NOT NULL,
  credits_used INT NOT NULL,
  used_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_credit_usage_daily ON credit_usage (api_key_id, used_at);
```

### Credit Check Flow (Cloudflare Worker)

```
1. Extract x-api-key header
2. Lookup key in DB/KV
3. Get today's usage: SUM(credits_used) WHERE used_at >= today
4. Check: used + endpoint_cost <= daily_limit
5. If over limit: return 429 with X-Credits-Remaining: 0
6. If OK: forward request, then record usage
```

### Endpoint Cost Config

```json
{
  "GET /v1/text/advice": 1,
  "GET /v1/text/quotes/random": 1,
  "GET /v1/text/words/random": 1,
  "GET /v1/text/words/define": 2,
  "GET /v1/finance/commodities/:symbol": 3,
  "GET /v1/places/geocode": 5
}
```

Store in KV for fast lookups at the edge.

---

## FAQ

**Q: Do unused credits roll over?**\
A: No, credits reset daily at midnight UTC.

**Q: What happens if I hit the limit mid-request?**\
A: The request is rejected with `429 Too Many Requests`.

**Q: Can I buy more credits without upgrading?**\
A: Not yet. Upgrade to a paid plan for higher limits + overage billing.

**Q: Why credits instead of request counts?**\
A: Fair pricing. A simple text lookup shouldn't cost the same as a geocoding
request that calls external APIs.
