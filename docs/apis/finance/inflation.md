# Inflation API

Annual CPI inflation rates for 241 countries, sourced from the World Bank.

## Endpoint

`GET /v1/finance/inflation?country=US`

## Parameters

| Name      | Type   | Required | Description                                                      |
| --------- | ------ | -------- | ---------------------------------------------------------------- |
| `country` | string | yes      | ISO 3166-1 alpha-2 country code (e.g. US, GB). Case-insensitive. |

## Response

```json
{
  "data": {
    "country": "US",
    "rate": 2.9495,
    "period": "2024",
    "historical": [
      { "period": "2023", "rate": 4.1163 },
      { "period": "2022", "rate": 8.0028 }
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

- `rate` — latest annual CPI % change (e.g. `2.9495` = 2.9495%)
- `period` — year of the latest data point
- `historical` — up to 10 previous years, newest first

## Data Source

World Bank indicator `FP.CPI.TOTL.ZG`. Updated annually; re-run
`cmd/seed-inflation` to pull new figures when the World Bank publishes them.

## Error Codes

| Code             | Status | When                                           |
| ---------------- | ------ | ---------------------------------------------- |
| `bad_request`    | 400    | Missing or invalid country code                |
| `not_found`      | 404    | No data for that country in the World Bank set |
| `internal_error` | 500    | Unexpected failure                             |
