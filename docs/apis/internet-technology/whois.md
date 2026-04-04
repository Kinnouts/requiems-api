# Whois API

## Status

✅ **Live**

## Overview

Get WHOIS registration details for any domain name. Returns registrar, name
servers, status flags, registration dates, and DNSSEC information.

## Endpoint

### WHOIS Lookup

`GET /v1/tech/whois/{domain}`

Returns WHOIS information for a domain.

| Parameter | Type   | Required | Description                                 |
| --------- | ------ | -------- | ------------------------------------------- |
| `domain`  | string | Yes      | Domain name to look up (e.g. `example.com`) |

## Response

```json
{
  "data": {
    "domain": "example.com",
    "registrar": "RESERVED-Internet Assigned Numbers Authority",
    "name_servers": [
      "A.IANA-SERVERS.NET",
      "B.IANA-SERVERS.NET"
    ],
    "status": [
      "clientDeleteProhibited",
      "clientTransferProhibited",
      "clientUpdateProhibited"
    ],
    "created_date": "1995-08-14T04:00:00Z",
    "updated_date": "2023-08-14T07:01:38Z",
    "expiry_date": "2024-08-13T04:00:00Z",
    "dnssec": true
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Error Codes

| Code             | Status | When                                            |
| ---------------- | ------ | ----------------------------------------------- |
| `bad_request`    | 400    | Domain name format is invalid                   |
| `not_found`      | 404    | Domain is not registered or no WHOIS data found |
| `internal_error` | 500    | Upstream WHOIS query failed                     |
