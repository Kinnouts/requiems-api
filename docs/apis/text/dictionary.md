# Dictionary API

## Status

✅ **MVP** - Available

## Overview

Get comprehensive word definitions, phonetics, usage examples, and synonyms for
any word in the dataset. Useful for vocabulary tools, educational apps, and
writing assistants.

## Endpoint

`GET /v1/text/dictionary/{word}`

## Path Parameters

| Parameter | Type   | Required | Description         |
| --------- | ------ | -------- | ------------------- |
| `word`    | string | Yes      | The word to look up |

## Response

```json
{
  "data": {
    "word": "ephemeral",
    "phonetic": "/ɪˈfɛm(ə)rəl/",
    "definitions": [
      {
        "partOfSpeech": "adjective",
        "definition": "lasting for a very short time",
        "example": "ephemeral pleasures"
      }
    ],
    "synonyms": ["transient", "fleeting", "momentary", "brief", "short-lived"]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Notes

- Lookup is **case-insensitive**: `Ephemeral`, `EPHEMERAL`, and `ephemeral` all
  return the same result.
- The `word` field in the response is always normalized to lowercase.
- Some words have multiple definitions (e.g. a word that is both a noun and an
  adjective).
- The `example` field may be omitted if no example sentence is available.
- `synonyms` is always present in the response (may be an empty array).

## Error Codes

| Code          | Status | When                                 |
| ------------- | ------ | ------------------------------------ |
| `not_found`   | 404    | Word not found in dictionary dataset |
| `bad_request` | 400    | Missing word path parameter          |
