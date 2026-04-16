# Random Word API

Get random words with definitions and parts of speech for vocabulary builders,
educational apps, and word games.

## Status

✅ **Live** - Available now at `GET /v1/text/words/random`

## Endpoint

`GET /v1/text/words/random`

## Query Parameters

None required.

## Response

```json
{
  "data": {
    "id": 123,
    "word": "ephemeral",
    "definition": "lasting for a very short time",
    "part_of_speech": "adjective"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field          | Type    | Description                                                      |
| -------------- | ------- | ---------------------------------------------------------------- |
| id             | integer | Unique identifier for the word                                   |
| word           | string  | The random word                                                  |
| definition     | string  | Dictionary definition of the word                                |
| part_of_speech | string  | Grammatical classification (e.g., noun, verb, adjective, adverb) |

## Error Codes

| Code                  | Status | When                           |
| --------------------- | ------ | ------------------------------ |
| `service_unavailable` | 503    | No words available in database |

## Code Examples

### cURL

```bash
curl https://api.requiems.xyz/v1/text/words/random \
  -H "requiems-api-key: YOUR_API_KEY"
```

### Python

```python
import requests

url = "https://api.requiems.xyz/v1/text/words/random"
headers = {"requiems-api-key": "YOUR_API_KEY"}

response = requests.get(url, headers=headers)
word_data = response.json()['data']
print(f"{word_data['word']} ({word_data['part_of_speech']})")
print(f"Definition: {word_data['definition']}")
```

### JavaScript

```javascript
const response = await fetch("https://api.requiems.xyz/v1/text/words/random", {
  headers: {
    "requiems-api-key": "YOUR_API_KEY",
  },
});

const { data } = await response.json();
console.log(`${data.word} (${data.part_of_speech})`);
console.log(`Definition: ${data.definition}`);
```

### Ruby

```ruby
require 'net/http'
require 'json'

uri = URI('https://api.requiems.xyz/v1/text/words/random')
request = Net::HTTP::Get.new(uri)
request['requiems-api-key'] = 'YOUR_API_KEY'

response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
  http.request(request)
end

word = JSON.parse(response.body)['data']
puts "#{word['word']} (#{word['part_of_speech']})"
puts "Definition: #{word['definition']}"
```

## Use Cases

- **Vocabulary Learning Apps** - Help users expand their vocabulary
- **Word-of-the-Day Features** - Display a new word each day
- **Educational Tools** - Create flashcards and learning materials
- **Writing Prompts** - Inspire creative writing with random words
- **Word Games** - Generate words for games and puzzles

## FAQ

**What parts of speech are included?** The API includes all major parts of
speech including nouns, verbs, adjectives, adverbs, prepositions, conjunctions,
and more.

**Can I filter by part of speech or difficulty level?** Currently, the API
returns completely random words. Filtering by part of speech or difficulty level
is planned for a future update.

**How many words are in the database?** Our collection contains thousands of
carefully curated English words with accurate definitions, covering a wide range
of vocabulary levels.

**Are the definitions from a specific dictionary?** The definitions are curated
from reputable dictionary sources to provide clear, concise explanations
suitable for general use.

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 864 ms  |
| p95     | 1081 ms |
| p99     | 1303 ms |
| Average | 915 ms  |

_Last updated: 2026-04-16_ Measured against production
(`https://api.requiems.xyz`) with 43 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 807 ms  |
| p95     | 1009 ms |
| p99     | 1092 ms |
| Average | 834 ms  |

_Last updated: 2026-04-16_
