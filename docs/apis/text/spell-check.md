# Spell Check API

## Status

✅ **Live**

## Overview

Check spelling and get correction suggestions. This endpoint identifies
misspelled words in the input text, provides the best correction per word, and
returns a rebuilt version of the text with all corrections applied.

## Endpoint

`POST /v1/text/spellcheck`

## Request

```json
{
  "text": "Ths is a smiple tset"
}
```

| Field  | Type   | Required | Description          |
| ------ | ------ | -------- | -------------------- |
| `text` | string | ✅       | The text to analyse. |

## Response

```json
{
  "data": {
    "corrected": "This is a simple test",
    "corrections": [
      { "original": "Ths", "suggested": "This", "position": 0 },
      { "original": "smiple", "suggested": "simple", "position": 9 },
      { "original": "tset", "suggested": "test", "position": 16 }
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field         | Type             | Description                                                   |
| ------------- | ---------------- | ------------------------------------------------------------- |
| `corrected`   | string           | Input text with misspelled words replaced by their correction |
| `corrections` | array of objects | One entry per misspelled word (see below)                     |

Each `corrections` entry:

| Field       | Type   | Description                                            |
| ----------- | ------ | ------------------------------------------------------ |
| `original`  | string | The misspelled word as it appears in the original text |
| `suggested` | string | The best correction found                              |
| `position`  | int    | 0-based character offset of the word in the input text |

## Notes

- Only **English** is supported.
- Positions are **0-based** character offsets in the original input string.
- Capitalisation is **preserved**: if the original word starts with an uppercase
  letter, the suggestion will too.
- When no mistakes are found, `corrections` is an empty array (`[]`) and
  `corrected` equals the input.
- Only **ASCII letter sequences** (`[a-zA-Z]`) are spell-checked. Non-ASCII
  characters (accented letters, CJK, emoji, etc.) are passed through unchanged
  and do not affect position counting.

## Error Codes

| Code                | Status | When                                |
| ------------------- | ------ | ----------------------------------- |
| `validation_failed` | 422    | `text` field is missing or empty    |
| `bad_request`       | 400    | Request body is missing or not JSON |
| `internal_error`    | 500    | Unexpected failure                  |

## Performance

Measured against production (`https://api.requiems.xyz`) with 1 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 1029 ms |
| p95     | 1029 ms |
| p99     | 1029 ms |
| Average | 1029 ms |

_Last updated: 2026-04-16_
