# Entertainment APIs

## Overview

Entertainment-focused endpoints providing access to advice, quotes, jokes,
trivia, and other entertaining content.

## Endpoints

### [Advice](./advice.md) - ✅ MVP

Get random advice and wisdom.

- **Status:** mvp
- **Endpoint:** `GET /v1/text/advice`
- **Credit Cost:** 1

### [Quotes](./quotes.md) - ✅ MVP

Get random inspirational and famous quotes.

- **Status:** mvp
- **Endpoint:** `GET /v1/text/quotes/random`
- **Credit Cost:** 1

### [Horoscope](./horoscope.md) - ✅ MVP

Get daily horoscope readings for all 12 zodiac signs.

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/horoscope/{sign}`
- **Credit Cost:** 1

### [Trivia](./trivia.md) - ✅ MVP

Get random trivia questions with multiple-choice answers. Filter by category and
difficulty.

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/trivia`
- **Credit Cost:** 1

### [Random Facts](./facts.md) - ✅ MVP

Get random interesting facts. Filter by category (science, history, technology,
nature, space, food).

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/facts`
- **Credit Cost:** 1

### [Emoji](./emoji.md) - ✅ MVP

Look up emoji by name, search by keyword, or get a random emoji with full
Unicode metadata.

- **Status:** mvp
- **Endpoints:** `GET /v1/entertainment/emoji/random`,
  `GET /v1/entertainment/emoji/search`, `GET /v1/entertainment/emoji/{name}`
- **Credit Cost:** 1

### [Sudoku](./sudoku.md) - ✅ MVP

Generate Sudoku puzzles with solutions across multiple difficulty levels.

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/sudoku`
- **Credit Cost:** 1

### [Dad Jokes](./dad-jokes.md) - ✅ MVP

Get a random dad joke.

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/jokes/dad`
- **Credit Cost:** 1

### [Chuck Norris Facts](./chuck-norris.md) - ✅ MVP

Get a random Chuck Norris fact from a curated database.

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/chuck-norris`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 9
- Live: 9
