# Validate Phone API

Validate phone numbers and retrieve country, type, and international format.

## Endpoint

`GET /v1/tech/validate/phone`

## Query Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `number`  | string | ✓        | Phone number in E.164 format (e.g. `+12015551234`) |

## Response

```json
{
  "data": {
    "number": "+12015551234",
    "valid": true,
    "country": "US",
    "type": "landline_or_mobile",
    "formatted": "+1 201-555-1234"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

When the number is invalid, `valid` is `false` and the optional fields are omitted:

```json
{
  "data": {
    "number": "12345",
    "valid": false
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Number Types

| Value | Description |
|---|---|
| `mobile` | Mobile / cell phone |
| `landline` | Fixed-line telephone |
| `landline_or_mobile` | Cannot be distinguished (common for US numbers) |
| `toll_free` | Toll-free number |
| `premium_rate` | Premium-rate number |
| `shared_cost` | Shared-cost number |
| `voip` | Voice over IP |
| `personal_number` | Personal number |
| `pager` | Pager |
| `uan` | Universal Access Number |
| `voicemail` | Voicemail number |
| `unknown` | Type could not be determined |

## Response Fields

| Field     | Type    | Description |
|-----------|---------|-------------|
| number    | string  | The original number as supplied in the request |
| valid     | boolean | Whether the number is a valid, dialable phone number |
| country   | string  | ISO 3166-1 alpha-2 country code (omitted when valid is false) |
| type      | string  | Number type (see table above, omitted when valid is false) |
| formatted | string  | International format of the number (e.g., +1 201-555-1234, omitted when valid is false) |

## Error Codes

| Code | Status | When |
|---|---|---|
| `bad_request` | 400 | The `number` query parameter is missing |

## Code Examples

### cURL

```bash
curl "https://api.requiems.xyz/v1/tech/validate/phone?number=%2B12015551234" \
  -H "requiems-api-key: YOUR_API_KEY"
```

### Python

```python
import requests

url = "https://api.requiems.xyz/v1/tech/validate/phone"
headers = {"requiems-api-key": "YOUR_API_KEY"}
params = {"number": "+12015551234"}

response = requests.get(url, headers=headers, params=params)
result = response.json()['data']

if result['valid']:
    print(f"Valid {result['type']} number in {result['country']}")
    print(f"Formatted: {result['formatted']}")
else:
    print(f"Invalid number: {result['number']}")
```

### JavaScript

```javascript
const number = encodeURIComponent('+12015551234');
const response = await fetch(
  `https://api.requiems.xyz/v1/tech/validate/phone?number=${number}`,
  { headers: { 'requiems-api-key': 'YOUR_API_KEY' } }
);

const { data } = await response.json();

if (data.valid) {
  console.log(`Valid ${data.type} number in ${data.country}`);
  console.log(`Formatted: ${data.formatted}`);
} else {
  console.log(`Invalid number: ${data.number}`);
}
```

### Ruby

```ruby
require 'net/http'
require 'json'

uri = URI('https://api.requiems.xyz/v1/tech/validate/phone')
uri.query = URI.encode_www_form(number: '+12015551234')

request = Net::HTTP::Get.new(uri)
request['requiems-api-key'] = 'YOUR_API_KEY'

response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
  http.request(request)
end

data = JSON.parse(response.body)['data']

if data['valid']
  puts "Valid #{data['type']} number in #{data['country']}"
  puts "Formatted: #{data['formatted']}"
else
  puts "Invalid number: #{data['number']}"
end
```

## Use Cases

- **User Registration Forms** - Validate phone numbers during signup
- **Fraud Prevention** - Verify phone numbers for identity verification
- **Data Normalization** - Store phone numbers in consistent international format
- **UI Formatting** - Display phone numbers in user-friendly format

## FAQ

**What format should I use for the phone number?**
Always include the country calling code prefixed with a plus sign (E.164 format), for example +12015551234 or +447400123456. Numbers without a country code cannot be validated reliably.

**What does landline_or_mobile mean?**
Some numbering plans (notably the United States) do not distinguish between mobile and landline numbers at the format level. When the type cannot be determined more precisely, the API returns landline_or_mobile.

**What happens when the number is invalid?**
The API still returns HTTP 200 with valid set to false. The country, type, and formatted fields are omitted from the response.
