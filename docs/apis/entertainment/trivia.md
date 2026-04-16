# Trivia API

## Status

✅ **MVP** - Live

## Overview

Get random trivia questions with multiple-choice answers across 10 categories
and 3 difficulty levels. Each response includes the question, answer options,
the correct answer, and metadata about the question.

## Endpoints

### Get Trivia Question

**Endpoint:** `GET /v1/entertainment/trivia`

Get a random trivia question. Use the optional `category` and `difficulty` query
parameters to narrow down the question pool.

#### Query Parameters

| Parameter    | Type   | Required | Description                                                                                                                               |
| ------------ | ------ | -------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| `category`   | string | No       | Filter by category. One of: `science`, `history`, `geography`, `sports`, `music`, `movies`, `literature`, `math`, `technology`, `nature`. |
| `difficulty` | string | No       | Filter by difficulty. One of: `easy`, `medium`, `hard`.                                                                                   |

#### Response Fields

| Field        | Type          | Description                                              |
| ------------ | ------------- | -------------------------------------------------------- |
| `question`   | string        | The trivia question text.                                |
| `options`    | array[string] | Four multiple-choice answer options.                     |
| `answer`     | string        | The correct answer (always one of the `options` values). |
| `category`   | string        | The category the question belongs to.                    |
| `difficulty` | string        | The difficulty level: `easy`, `medium`, or `hard`.       |

#### Example Response

```json
{
  "data": {
    "question": "What is the largest planet in our solar system?",
    "options": ["Earth", "Jupiter", "Saturn", "Mars"],
    "answer": "Jupiter",
    "category": "science",
    "difficulty": "easy"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

#### Example — Filtered Request

```
GET /v1/entertainment/trivia?category=science&difficulty=medium
```

#### Error Responses

| Status | Code           | Description                                               |
| ------ | -------------- | --------------------------------------------------------- |
| 400    | `bad_request`  | An invalid `category` or `difficulty` value was provided. |
| 404    | `not_found`    | No questions match the given filters.                     |
| 401    | `unauthorized` | Missing API key.                                          |
| 403    | `forbidden`    | Invalid or revoked API key.                               |

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 910 ms  |
| p95     | 1169 ms |
| p99     | 1369 ms |
| Average | 966 ms  |

_Last updated: 2026-04-16_ Measured against production
(`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 814 ms  |
| p95     | 1084 ms |
| p99     | 1406 ms |
| Average | 854 ms  |

_Last updated: 2026-04-16_
