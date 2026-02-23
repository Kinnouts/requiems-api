# Available Conversion Units API

Retrieve a complete list of all available unit conversion types and their codes
for the Unit Conversion API.

## Status

✅ **Live** - Available now at `GET /v1/misc/convert/units`

## Endpoint

`GET /v1/misc/convert/units`

## Query Parameters

None required.

## Response

```json
{
  "data": {
    "length": ["cm", "ft", "in", "km", "m", "miles", "mm", "nmi", "yd"],
    "weight": ["g", "kg", "lb", "mg", "oz", "stone", "t"],
    "volume": ["cup", "fl_oz", "gal", "l", "ml", "pt", "qt", "tbsp", "tsp"],
    "temperature": ["c", "f", "k"],
    "area": ["acre", "cm2", "ft2", "ha", "in2", "km2", "m2", "mm2", "yd2"],
    "speed": ["km_h", "knots", "m_s", "mph"]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field       | Type             | Description                                                                                                                                                                                           |
| ----------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| length      | array of strings | Available length units: millimeter (mm), centimeter (cm), meter (m), kilometer (km), inch (in), foot (ft), yard (yd), mile (miles), nautical mile (nmi)                                               |
| weight      | array of strings | Available weight units: milligram (mg), gram (g), kilogram (kg), metric ton (t), ounce (oz), pound (lb), stone (stone)                                                                                |
| volume      | array of strings | Available volume units: milliliter (ml), liter (l), teaspoon (tsp), tablespoon (tbsp), fluid ounce (fl_oz), cup (cup), pint (pt), quart (qt), gallon (gal)                                            |
| temperature | array of strings | Available temperature units: celsius (c), fahrenheit (f), kelvin (k)                                                                                                                                  |
| area        | array of strings | Available area units: square millimeter (mm2), square centimeter (cm2), square meter (m2), square kilometer (km2), square inch (in2), square foot (ft2), square yard (yd2), acre (acre), hectare (ha) |
| speed       | array of strings | Available speed units: meters per second (m_s), kilometers per hour (km_h), miles per hour (mph), knots (knots)                                                                                       |

## Unit Categories

The API supports **39 units** across 6 measurement categories:

### Length (9 units)

- **Metric**: mm, cm, m, km
- **Imperial**: in, ft, yd, miles
- **Nautical**: nmi

### Weight (7 units)

- **Metric**: mg, g, kg, t
- **Imperial**: oz, lb, stone

### Volume (9 units)

- **Metric**: ml, l
- **Imperial**: tsp, tbsp, fl_oz, cup, pt, qt, gal

### Temperature (3 units)

- c (Celsius)
- f (Fahrenheit)
- k (Kelvin)

### Area (9 units)

- **Metric**: mm2, cm2, m2, km2, ha
- **Imperial**: in2, ft2, yd2, acre

### Speed (4 units)

- m_s (meters per second)
- km_h (kilometers per hour)
- mph (miles per hour)
- knots

## Code Examples

### cURL

```bash
curl https://api.requiems.xyz/v1/misc/convert/units \
  -H "requiems-api-key: YOUR_API_KEY"
```

### Python

```python
import requests

url = "https://api.requiems.xyz/v1/misc/convert/units"
headers = {"requiems-api-key": "YOUR_API_KEY"}

response = requests.get(url, headers=headers)
units = response.json()['data']

# Print all length units
print("Length units:", units['length'])

# Build a unit conversion dropdown
for category, unit_list in units.items():
    print(f"\n{category.upper()}:")
    for unit in unit_list:
        print(f"  - {unit}")
```

### JavaScript

```javascript
const response = await fetch("https://api.requiems.xyz/v1/misc/convert/units", {
  headers: {
    "requiems-api-key": "YOUR_API_KEY",
  },
});

const { data } = await response.json();

// Build a dropdown with length units
const lengthSelect = document.createElement("select");
data.length.forEach((unit) => {
  const option = document.createElement("option");
  option.value = unit;
  option.textContent = unit;
  lengthSelect.appendChild(option);
});

// Display all categories
console.log("Available unit categories:", Object.keys(data));
```

### Ruby

```ruby
require 'net/http'
require 'json'

uri = URI('https://api.requiems.xyz/v1/misc/convert/units')
request = Net::HTTP::Get.new(uri)
request['requiems-api-key'] = 'YOUR_API_KEY'

response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
  http.request(request)
end

units = JSON.parse(response.body)['data']

# Display temperature units
puts "Temperature units: #{units['temperature'].join(', ')}"

# Count total units
total = units.values.map(&:length).sum
puts "Total units available: #{total}"
```

## Use Cases

- **Build Unit Selection Dropdowns** - Populate UI dropdowns with available
  conversion options
- **Discover Available Conversions** - Find out which units can be converted
- **Validate User Input** - Check if a user-provided unit code is valid
- **Display Available Options** - Show users what conversions are possible
- **Cache Unit List** - Store the list locally to avoid repeated API calls

## FAQ

**Do I need to call this endpoint every time I want to convert units?** No. This
endpoint is meant to be called once to discover available units. You can cache
the response and use it to build your UI or validate user input. The list of
available units changes very infrequently.

**Can I convert between different categories (e.g., length to weight)?** No. You
can only convert between units within the same category. For example, you can
convert meters to feet (both length), but not meters to kilograms (length to
weight). The conversion API will return a 400 error if you try to convert
incompatible units.

**What do the unit codes mean?** Each unit has a short code (e.g., "km" for
kilometer, "lb" for pound). These codes are used as parameters in the conversion
API. The codes are standardized and case-sensitive (always lowercase).

**Will new units be added in the future?** Yes, we periodically add new units
based on user requests. If you need a specific unit that's not currently
available, please contact support@requiems.xyz.

**How many units are currently supported?** The API currently supports 39 units
across 6 categories: length (9), weight (7), volume (9), temperature (3), area
(9), and speed (4).
