# Thesaurus API

## Status

✅ **MVP** - Available

## Overview

Find synonyms and antonyms for any word. Returns a curated list of related and
opposite words for the given input.

## Endpoint

`GET /v1/text/thesaurus/{word}`

## Path Parameters

| Parameter | Type   | Required | Description             |
| --------- | ------ | -------- | ----------------------- |
| `word`    | string | Yes      | The word to look up     |

## Response

```json
{
  "data": {
    "word": "happy",
    "synonyms": ["joyful", "cheerful", "content", "pleased", "delighted", "glad", "elated", "blissful"],
    "antonyms": ["sad", "unhappy", "miserable", "sorrowful", "dejected", "gloomy", "melancholy"]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Notes

- Lookup is **case-insensitive**: `Happy`, `HAPPY`, and `happy` all return the same result.
- The `word` field in the response is always normalized to lowercase.
- Both `synonyms` and `antonyms` are always present in the response (may be empty arrays).

## Error Codes

| Code          | Status | When                                |
| ------------- | ------ | ----------------------------------- |
| `not_found`   | 404    | Word not found in thesaurus dataset |
| `bad_request` | 400    | Missing word path parameter         |

## Credit Cost

1 credit per request.
