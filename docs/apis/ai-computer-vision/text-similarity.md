# Text Similarity API

## Status

✅ **Live** - `/v1/ai/similarity`

## Overview

Compare two text documents for similarity using cosine similarity on
word-frequency vectors. Returns a score between 0 (no overlap) and 1 (identical
word distribution).

## Endpoints

### Compare Text Similarity

**Endpoint:** `POST /v1/ai/similarity`

#### Request

```json
{
  "text1": "The cat sat on the mat",
  "text2": "A cat was sitting on a mat"
}
```

| Field   | Type   | Required | Description             |
| ------- | ------ | -------- | ----------------------- |
| `text1` | string | ✅       | First text to compare.  |
| `text2` | string | ✅       | Second text to compare. |

#### Response

```json
{
  "data": {
    "similarity": 0.4364,
    "method": "cosine"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field        | Type   | Description                                    |
| ------------ | ------ | ---------------------------------------------- |
| `similarity` | number | Cosine similarity score in the range \[0, 1\]. |
| `method`     | string | Algorithm used. Currently always `"cosine"`.   |

## Algorithm

Cosine similarity is computed on term-frequency (bag-of-words) vectors:

1. Both texts are lowercased and tokenised into alphanumeric words.
2. A word-frequency map is built for each text.
3. The dot product of the two frequency vectors is divided by the product of
   their magnitudes: `similarity = (A · B) / (|A| × |B|)`.

A score of **1** means the texts have an identical word distribution.\
A score of **0** means the texts share no words in common.
