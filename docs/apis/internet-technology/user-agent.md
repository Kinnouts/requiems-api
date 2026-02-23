# User Agent API

## Status

✅ **Live** — Available at `GET /v1/tech/useragent`

## Overview

Parse and analyze user agent strings to extract browser name, version, operating
system, device type, and bot detection.

## Endpoint

### Parse User Agent

`GET /v1/tech/useragent?ua=<encoded-ua-string>`

**Credit Cost:** 1 credit per request

### Query Parameters

| Parameter | Required | Description |
|---|---|---|
| `ua` | Yes | URL-encoded user agent string |

### Response

```json
{
  "data": {
    "browser": "Chrome",
    "browser_version": "120.0",
    "os": "Windows",
    "os_version": "10/11",
    "device": "desktop",
    "is_bot": false
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

### Device Values

| Value | Description |
|---|---|
| `desktop` | Desktop browser on Windows, macOS, or Linux |
| `mobile` | Mobile browser (iPhone, Android Mobile) |
| `tablet` | Tablet browser (iPad, Android without Mobile) |
| `bot` | Known bot or crawler |
| `unknown` | Empty user agent string |

### Error Codes

| Code | Status | When |
|---|---|---|
| `bad_request` | 400 | `ua` query parameter is missing |
