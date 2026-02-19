# Unit Conversion API

## Status

✅ **Live**

## Overview

Convert between different units of measurement. Supports length, weight,
volume, temperature, area, speed, data, and time.

## Endpoint

### Convert Units

**Endpoint:** `GET /v1/misc/convert`

**Credit Cost:** 1 credit

Convert between any two compatible units of measurement.

### Query Parameters

| Parameter | Type   | Required | Description                     |
|-----------|--------|----------|---------------------------------|
| `from`    | string | Yes      | Source unit (e.g. `miles`)      |
| `to`      | string | Yes      | Target unit (e.g. `km`)         |
| `value`   | number | Yes      | Numeric value to convert        |

### Example Request

```
GET /v1/misc/convert?from=miles&to=km&value=10
```

### Example Response

```json
{
  "from": "miles",
  "to": "km",
  "input": 10,
  "result": 16.09344,
  "formula": "miles × 1.609344"
}
```

### Error Responses

| Status | Description                                    |
|--------|------------------------------------------------|
| 400    | Missing or invalid query parameters            |
| 400    | Unknown unit name                              |
| 400    | Incompatible unit types (e.g. length vs weight)|

## Supported Units

### Length
`mm`, `cm`, `m`, `km`, `inch`/`inches`, `ft`/`foot`/`feet`, `yard`/`yards`,
`mile`/`miles`, `nmi` (nautical mile)

### Weight
`mg`, `g`, `kg`, `tonne`, `lb`/`lbs`/`pound`/`pounds`, `oz`/`ounce`/`ounces`,
`stone`/`stones`

### Volume
`ml`, `l`/`liter`/`litre`, `fl_oz`, `cup`/`cups`, `pint`/`pints`,
`quart`/`quarts`, `gallon`/`gallons`

### Temperature
`celsius`/`c`/`°c`, `fahrenheit`/`f`/`°f`, `kelvin`/`k`

### Area
`mm2`, `cm2`, `m2`, `km2`, `sqft`, `acre`/`acres`, `hectare`/`hectares`/`ha`,
`sqmi`

### Speed
`m/s`, `km/h`/`kph`, `mph`, `knot`/`knots`, `ft/s`

### Data
`bit`/`bits`, `byte`/`bytes`, `kb`, `mb`, `gb`, `tb`

### Time
`ms`, `s`/`second`/`seconds`, `min`/`minute`/`minutes`,
`hr`/`hour`/`hours`, `day`/`days`, `week`/`weeks`

