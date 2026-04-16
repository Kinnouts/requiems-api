# Quotes API

## Status

✅ **MVP** - Basic implementation live, production-ready

## Overview

Get random quotes and inspirational messages. This endpoint provides
motivational and thought-provoking quotes from various authors and sources.

## Live Endpoints

### Get Random Quote

**Endpoint:** `GET /v1/entertainment/quotes/random`

Get a random inspirational quote.

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 865 ms  |
| p95     | 1092 ms |
| p99     | 1540 ms |
| Average | 914 ms  |

_Last updated: 2026-04-16_ Measured against production
(`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 817 ms  |
| p95     | 941 ms  |
| p99     | 1150 ms |
| Average | 868 ms  |

_Last updated: 2026-04-16_
