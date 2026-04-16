# Sudoku API

## Status

✅ **MVP** - Live

## Overview

Generate Sudoku puzzles of varying difficulty levels, complete with their
solutions. Each request returns a freshly generated, unique puzzle.

## Endpoints

### Get Sudoku Puzzle

**Endpoint:** `GET /v1/entertainment/sudoku`

Generate a random Sudoku puzzle.

#### Query Parameters

| Parameter    | Type   | Required | Description                                                           |
| ------------ | ------ | -------- | --------------------------------------------------------------------- |
| `difficulty` | string | No       | Puzzle difficulty: `easy`, `medium`, or `hard`. Defaults to `medium`. |

#### Response Fields

| Field        | Type         | Description                                         |
| ------------ | ------------ | --------------------------------------------------- |
| `difficulty` | string       | The difficulty level of the returned puzzle.        |
| `puzzle`     | array[array] | 9×9 grid — `0` represents an empty cell to fill in. |
| `solution`   | array[array] | 9×9 grid containing the complete solution.          |

#### Example Response

```json
{
  "data": {
    "difficulty": "hard",
    "puzzle": [
      [5, 3, 0, 0, 7, 0, 0, 0, 0],
      [6, 0, 0, 1, 9, 5, 0, 0, 0],
      [0, 9, 8, 0, 0, 0, 0, 6, 0],
      [8, 0, 0, 0, 6, 0, 0, 0, 3],
      [4, 0, 0, 8, 0, 3, 0, 0, 1],
      [7, 0, 0, 0, 2, 0, 0, 0, 6],
      [0, 6, 0, 0, 0, 0, 2, 8, 0],
      [0, 0, 0, 4, 1, 9, 0, 0, 5],
      [0, 0, 0, 0, 8, 0, 0, 7, 9]
    ],
    "solution": [
      [5, 3, 4, 6, 7, 8, 9, 1, 2],
      [6, 7, 2, 1, 9, 5, 3, 4, 8],
      [1, 9, 8, 3, 4, 2, 5, 6, 7],
      [8, 5, 9, 7, 6, 1, 4, 2, 3],
      [4, 2, 6, 8, 5, 3, 7, 9, 1],
      [7, 1, 3, 9, 2, 4, 8, 5, 6],
      [9, 6, 1, 5, 3, 7, 2, 8, 4],
      [2, 8, 7, 4, 1, 9, 6, 3, 5],
      [3, 4, 5, 2, 8, 6, 1, 7, 9]
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Performance

Measured against production (`https://api.requiems.xyz`) with 1 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 1095 ms |
| p95     | 1095 ms |
| p99     | 1095 ms |
| Average | 1095 ms |

_Last updated: 2026-04-16_ Measured against production
(`https://api.requiems.xyz`) with 1 samples.

| Metric  | Value  |
| ------- | ------ |
| p50     | 933 ms |
| p95     | 933 ms |
| p99     | 933 ms |
| Average | 933 ms |

_Last updated: 2026-04-16_
