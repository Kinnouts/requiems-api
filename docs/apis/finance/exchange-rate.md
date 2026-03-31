# Exchange Rate API

## Status

✅ **Live**

## Overview

Get live currency exchange rates and convert amounts between currencies. Rates
are sourced from the European Central Bank via
[Frankfurter](https://www.frankfurter.app/) and cached in Redis for up to one
hour. Supports ~170 ISO 4217 currency codes.

## Endpoints

### Get Exchange Rate

`GET /v1/finance/exchange-rate`

Returns the current exchange rate between two currencies.

**Query parameters:**

| Parameter | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `from`    | string | yes      | ISO 4217 source currency (3 chars) |
| `to`      | string | yes      | ISO 4217 target currency (3 chars) |

**Example:**

```
GET /v1/finance/exchange-rate?from=USD&to=EUR
```

**Response:**

```json
{
  "data": {
    "from": "USD",
    "to": "EUR",
    "rate": 0.92,
    "timestamp": "2024-12-15T00:00:00Z"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

---

### Convert Currency

`GET /v1/finance/convert`

Converts an amount from one currency to another. Returns the rate alongside the
converted value.

**Query parameters:**

| Parameter | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `from`    | string | yes      | ISO 4217 source currency (3 chars) |
| `to`      | string | yes      | ISO 4217 target currency (3 chars) |
| `amount`  | number | yes      | Amount to convert (must be > 0)    |

**Example:**

```
GET /v1/finance/convert?from=USD&to=EUR&amount=100
```

**Response:**

```json
{
  "data": {
    "from": "USD",
    "to": "EUR",
    "rate": 0.92,
    "amount": 100,
    "converted": 92.00,
    "timestamp": "2024-12-15T00:00:00Z"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

The `converted` field is `amount × rate`, rounded to 2 decimal places.

---

## Caching

Each currency pair is cached in Redis with a 1-hour TTL. A cache miss triggers
a single HTTP call to `api.frankfurter.app`. Subsequent requests within the
hour are served instantly from cache.

## Error Codes

| Code               | Status | When                                                  |
| ------------------ | ------ | ----------------------------------------------------- |
| `bad_request`      | 400    | Missing parameter or currency code is not 3 letters   |
| `invalid_currency` | 422    | Currency code not recognised by the ECB data source   |
| `upstream_error`   | 503    | Frankfurter upstream temporarily unavailable          |

## Data Source

Rates are provided by [Frankfurter](https://www.frankfurter.app/), an
open-source API backed by the European Central Bank. Rates are updated on ECB
business days, typically around 16:00 CET.
