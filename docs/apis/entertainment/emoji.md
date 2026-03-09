# Emoji API

## Status

✅ **Live**

## Overview

Get emoji information and search emojis. This endpoint provides emoji data
including the rendered glyph, CLDR name, Unicode category, and code-point.

## Endpoints

### Get Random Emoji

**Endpoint:** `GET /v1/entertainment/emoji/random`
**Credit Cost:** 1 credit

Returns a randomly selected emoji.

### Search Emoji

**Endpoint:** `GET /v1/entertainment/emoji/search?q=happy`
**Credit Cost:** 1 credit

Searches for emojis whose name or category contains the given query string
(case-insensitive). Returns a list of all matches.

**Query Parameters:**

| Parameter | Required | Description                        |
| --------- | -------- | ---------------------------------- |
| `q`       | Yes      | Term to match against names/categories |

### Get Emoji by Name

**Endpoint:** `GET /v1/entertainment/emoji/:name`
**Credit Cost:** 1 credit

Returns a specific emoji by its CLDR snake\_case name.

## Response

```json
{
  "data": {
    "emoji": "😀",
    "name": "grinning_face",
    "category": "Smileys & Emotion",
    "unicode": "U+1F600"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

For the search endpoint the response wraps a list:

```json
{
  "data": {
    "items": [ { "emoji": "😀", "name": "grinning_face", ... } ],
    "total": 1
  },
  "metadata": { "timestamp": "2026-01-01T00:00:00Z" }
}
```

## Categories

| Category             | Example emojis                    |
| -------------------- | --------------------------------- |
| Smileys & Emotion    | 😀 😂 ❤️ 💀                     |
| People & Body        | 👋 👍 💪 🙏                     |
| Animals & Nature     | 🐶 🦁 🌹 🍄                     |
| Food & Drink         | 🍕 🍔 ☕ 🍺                     |
| Travel & Places      | 🚗 ✈️ 🏠 🌍                     |
| Activities           | ⚽ 🎮 🏆 🎨                     |
| Objects              | 💻 📷 🔑 💡                     |
| Symbols              | ✅ ❌ ♻️ ❓                     |

## Error Codes

| Code          | Status | When                                  |
| ------------- | ------ | ------------------------------------- |
| `bad_request` | 400    | Missing or empty `q` query parameter  |
| `not_found`   | 404    | No emoji found with the requested name |

