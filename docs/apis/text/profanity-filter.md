# Profanity Filter API

Detect and censor profanity in text for content moderation.

## Endpoint

`POST /v1/text/profanity`

## Request

```json
{
  "text": "Some text to check"
}
```

| Field | Type   | Required | Description             |
|-------|--------|----------|-------------------------|
| text  | string | Yes      | The text to check       |

## Response

```json
{
  "data": {
    "hasProfanity": false,
    "censored": "Some text to check",
    "flaggedWords": []
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field        | Type             | Description                                               |
|--------------|------------------|-----------------------------------------------------------|
| hasProfanity | boolean          | Whether any profanity was detected                        |
| censored     | string           | Input text with profane words replaced by `*` characters |
| flaggedWords | array of strings | Deduplicated list of detected profane words (lowercase)   |

## Behaviour

- Detection is **case-insensitive** — `BULLSHIT`, `Bullshit`, and `bullshit` all match.
- Censoring replaces each character of a flagged word with `*`, preserving word length.
- Surrounding punctuation and whitespace are left unchanged.
- `flaggedWords` contains each word only once, even if it appears multiple times in the input.

## Error Codes

| Code               | Status | When                                    |
|--------------------|--------|-----------------------------------------|
| `validation_failed` | 422   | The `text` field is missing or empty    |
| `bad_request`       | 400   | The request body is missing or malformed |
| `internal_error`    | 500   | Unexpected server error                 |

