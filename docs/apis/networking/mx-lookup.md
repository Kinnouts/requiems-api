# MX Lookup API

## Status

✅ **MVP** - Implemented

## Overview

Perform MX record lookups. Returns all mail exchange records for a domain,
sorted by priority ascending (lowest numeric value = highest delivery priority
per RFC 5321).

## Endpoints

### MX Lookup

**Endpoint:** `GET /v1/networking/mx/{domain}`

Look up MX records for a domain.

| Parameter | Type   | Required | Description            |
| --------- | ------ | -------- | ---------------------- |
| `domain`  | string | Yes      | Domain name to look up |

### Response

```json
{
  "data": {
    "domain": "gmail.com",
    "records": [
      { "host": "gmail-smtp-in.l.google.com.", "priority": 5 },
      { "host": "alt1.gmail-smtp-in.l.google.com.", "priority": 10 }
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

### Errors

| Code             | Status | Description                        |
| ---------------- | ------ | ---------------------------------- |
| `bad_request`    | 400    | Invalid domain name format         |
| `not_found`      | 404    | No MX records found for the domain |
| `internal_error` | 500    | DNS lookup failed unexpectedly     |
