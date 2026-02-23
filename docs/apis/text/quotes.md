# Random Quotes API

Get random inspirational quotes from famous people, thinkers, and leaders.

## Status

✅ **Live** - Available now at `GET /v1/text/quotes/random`

## Endpoint

`GET /v1/text/quotes/random`

## Query Parameters

None required.

## Response

```json
{
  "data": {
    "id": 42,
    "text": "The only way to do great work is to love what you do.",
    "author": "Steve Jobs"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field  | Type    | Description                    |
|--------|---------|--------------------------------|
| id     | integer | Unique identifier for the quote |
| text   | string  | The quote text                 |
| author | string  | Name of the person who said or wrote the quote |

## Error Codes

| Code                  | Status | When                            |
|-----------------------|--------|---------------------------------|
| `service_unavailable` | 503    | No quotes available in database |

## Code Examples

### cURL

```bash
curl https://api.requiems.xyz/v1/text/quotes/random \
  -H "requiems-api-key: YOUR_API_KEY"
```

### Python

```python
import requests

url = "https://api.requiems.xyz/v1/text/quotes/random"
headers = {"requiems-api-key": "YOUR_API_KEY"}

response = requests.get(url, headers=headers)
quote = response.json()['data']
print(f'"{quote["text"]}" - {quote["author"]}')
```

### JavaScript

```javascript
const response = await fetch('https://api.requiems.xyz/v1/text/quotes/random', {
  headers: {
    'requiems-api-key': 'YOUR_API_KEY'
  }
});

const { data } = await response.json();
console.log(`"${data.text}" - ${data.author}`);
```

### Ruby

```ruby
require 'net/http'
require 'json'

uri = URI('https://api.requiems.xyz/v1/text/quotes/random')
request = Net::HTTP::Get.new(uri)
request['requiems-api-key'] = 'YOUR_API_KEY'

response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
  http.request(request)
end

quote = JSON.parse(response.body)['data']
puts "\"#{quote['text']}\" - #{quote['author']}"
```

## Use Cases

- **Daily Motivation Apps** - Display a new quote each day
- **Social Media Content** - Generate shareable quote graphics
- **Quote of the Day Features** - Add inspiration to websites and dashboards
- **Content Generation** - Enhance blog posts and articles with relevant quotes

## FAQ

**Can I request quotes from a specific author?**
Currently, the API returns random quotes from our entire collection. Author filtering is planned for a future update.

**How many quotes are in the database?**
Our collection contains hundreds of curated inspirational quotes from famous thinkers, leaders, and innovators, and we're constantly adding more.

**Will I get the same quote if I call the endpoint multiple times?**
No, quotes are selected randomly on each request, so consecutive calls will typically return different quotes.
