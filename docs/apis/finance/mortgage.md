# Mortgage Calculator API

Monthly payment and full amortization schedule for any fixed-rate loan.

## Endpoint

`GET /v1/finance/mortgage?principal=300000&rate=6.5&years=30`

## Parameters

| Name        | Type    | Required | Description                                                          |
| ----------- | ------- | -------- | -------------------------------------------------------------------- |
| `principal` | number  | yes      | Loan amount (currency-agnostic, e.g. 300000)                         |
| `rate`      | number  | yes      | Annual interest rate as a percentage (e.g. 6.5 = 6.5%). Must be > 0. |
| `years`     | integer | yes      | Loan term in years (1–50)                                            |

## Response

```json
{
  "data": {
    "principal": 300000,
    "rate": 6.5,
    "years": 30,
    "monthly_payment": 1896.2,
    "total_payment": 682632.0,
    "total_interest": 382632.0,
    "schedule": [
      {
        "month": 1,
        "payment": 1896.2,
        "principal": 271.2,
        "interest": 1625.0,
        "balance": 299728.8
      },
      {
        "month": 2,
        "payment": 1896.2,
        "principal": 272.67,
        "interest": 1623.53,
        "balance": 299456.13
      }
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

- `monthly_payment` — fixed monthly payment (rounded to 2 decimal places)
- `total_payment` — monthly_payment × (years × 12)
- `total_interest` — total_payment − principal
- `schedule` — one entry per month; length is always years × 12

## Formula

Standard fixed-rate amortization:

```
M = P × (r(1+r)^n) / ((1+r)^n − 1)
```

Where `P` is the principal, `r` is the monthly rate (`annual_rate / 100 / 12`),
and `n` is the total number of payments (`years × 12`).

## Error Codes

| Code             | Status | When                                                                  |
| ---------------- | ------ | --------------------------------------------------------------------- |
| `bad_request`    | 400    | Missing parameter, non-numeric value, rate ≤ 0, or years outside 1–50 |
| `internal_error` | 500    | Unexpected failure                                                    |

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 818 ms  |
| p95     | 921 ms  |
| p99     | 1300 ms |
| Average | 833 ms  |

_Last updated: 2026-04-16_
