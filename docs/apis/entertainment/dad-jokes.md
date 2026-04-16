# Dad Jokes API

## Status

✅ **Available**

## Overview

Get classic dad jokes. Groan-worthy puns and wholesome humor that only dads seem
to enjoy.

## Endpoints

### Get a Random Dad Joke

**Endpoint:** `GET /v1/entertainment/jokes/dad`

Get a random dad joke.

#### Response

```json
{
  "data": {
    "id": "joke_7",
    "joke": "Why don't scientists trust atoms? Because they make up everything!"
  },
  "metadata": {
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

#### Fields

| Field  | Type   | Description                    |
| ------ | ------ | ------------------------------ |
| `id`   | string | Unique identifier for the joke |
| `joke` | string | The dad joke text              |

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 921 ms  |
| p95     | 1174 ms |
| p99     | 1232 ms |
| Average | 955 ms  |

_Last updated: 2026-04-16_ Measured against production
(`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 815 ms  |
| p95     | 1153 ms |
| p99     | 1268 ms |
| Average | 888 ms  |

_Last updated: 2026-04-16_
