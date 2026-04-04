# Entertainment APIs

## Overview

Entertainment-focused endpoints providing access to advice, quotes, jokes,
trivia, and other entertaining content.

## Endpoints

### [Advice](./advice.md) - ✅ MVP

Get random advice for various life situations

- **Status:** mvp
- **Endpoint:** `GET /v1/text/advice`

### [Quotes](./quotes.md) - ✅ MVP

Get random quotes and inspirational messages

- **Status:** mvp
- **Endpoint:** `GET /v1/text/quotes/random`

### [Bucket List](./bucket-list.md) - ⏳ Planned

Get bucket list ideas and adventure suggestions

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/bucket-list`

### [Celebrity](./celebrity.md) - ⏳ Planned

Get information about celebrities and famous people

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/celebrity`

### [Chuck Norris](./chuck-norris.md) - ⏳ Planned

Get Chuck Norris jokes and facts

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/chuck-norris`

### [Dad Jokes](./dad-jokes.md) - ✅ Live

Get classic dad jokes

- **Status:** live
- **Endpoint:** `GET /v1/entertainment/jokes/dad`

### [Day in History](./day-in-history.md) - ⏳ Planned

Get historical events that happened on a specific day

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/day-in-history`

### [Emoji](./emoji.md) - ✅ Live

Get emoji information and random emojis

- **Status:** live
- **Endpoint:** `GET /v1/entertainment/emoji/random`
- **Endpoint:** `GET /v1/entertainment/emoji/search?q=happy`
- **Endpoint:** `GET /v1/entertainment/emoji/:name`

### [Facts](./facts.md) - ⏳ Planned

Get random interesting facts

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/facts`

### [Hobbies](./hobbies.md) - ⏳ Planned

Get hobby suggestions and recommendations

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/hobbies`

### [Horoscope](./horoscope.md) - ✅ MVP

Get horoscope readings for zodiac signs

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/horoscope/{sign}`

### [Jokes](./jokes.md) - ⏳ Planned

Get random jokes of various types

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/jokes`

### [Riddles](./riddles.md) - ⏳ Planned

Get riddles and brain teasers

- **Status:** planned
- **Planned Endpoint:** `GET /v1/entertainment/riddles`

### [Sudoku](./sudoku.md) - ✅ MVP

Get Sudoku puzzles and solutions

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/sudoku`

### [Trivia](./trivia.md) - ✅ MVP

Get trivia questions and answers

- **Status:** mvp
- **Endpoint:** `GET /v1/entertainment/trivia`
- **Endpoint:**
  `GET /v1/entertainment/trivia?category=science&difficulty=medium`

## Category Statistics

- Total Endpoints: 18
- Live: 5
- Planned: 11
