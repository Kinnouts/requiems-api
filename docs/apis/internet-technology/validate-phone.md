# Validate Phone API

Validate phone numbers and retrieve country, type, and international format.

## Endpoint

`GET /v1/tech/validate/phone`

## Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `number`  | string | ✓        | Phone number in E.164 format (e.g. `+12015551234`) |

## Response

```json
{
  "data": {
    "number": "+12015551234",
    "valid": true,
    "country": "US",
    "type": "landline_or_mobile",
    "formatted": "+1 201-555-1234"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

When the number is invalid, `valid` is `false` and the optional fields are omitted:

```json
{
  "data": {
    "number": "12345",
    "valid": false
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Number Types

| Value | Description |
|---|---|
| `mobile` | Mobile / cell phone |
| `landline` | Fixed-line telephone |
| `landline_or_mobile` | Cannot be distinguished (common for US numbers) |
| `toll_free` | Toll-free number |
| `premium_rate` | Premium-rate number |
| `shared_cost` | Shared-cost number |
| `voip` | Voice over IP |
| `personal_number` | Personal number |
| `pager` | Pager |
| `uan` | Universal Access Number |
| `voicemail` | Voicemail number |
| `unknown` | Type could not be determined |

## Error Codes

| Code | Status | When |
|---|---|---|
| `bad_request` | 400 | The `number` query parameter is missing |
