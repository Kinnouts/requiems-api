# Text Language API

## Status

✅ **Live** - Implemented as `POST /v1/ai/detect-language`

## Overview

Detect language of text. This endpoint identifies the language of provided
text, returning the full language name, ISO 639-1 code, and a confidence score.

## Endpoint

`POST /v1/ai/detect-language`

## Request

```json
{
  "text": "Bonjour, comment ça va?"
}
```

## Response

```json
{
  "data": {
    "language": "French",
    "code": "fr",
    "confidence": 0.98
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Supported Languages

75 languages are supported including English, French, Spanish, German, Italian,
Portuguese, Dutch, Russian, Chinese, Japanese, Arabic, and many more.

## Response Fields

| Field        | Type   | Description                                                              |
| ------------ | ------ | ------------------------------------------------------------------------ |
| `language`   | string | Full name of the detected language                                       |
| `code`       | string | ISO 639-1 two-letter language code. Empty string when unreliable.        |
| `confidence` | float  | Confidence score between 0.0 and 1.0. 0.0 when detection is unreliable. |

## Error Codes

| Code                | Status | When                          |
| ------------------- | ------ | ----------------------------- |
| `validation_failed` | 422    | Invalid or missing text field |
| `bad_request`       | 400    | Malformed or missing body     |
| `internal_error`    | 500    | Unexpected failure            |
