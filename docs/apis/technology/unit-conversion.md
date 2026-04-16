# Unit Conversion API

## Status

✅ **Live** - Available now

## Overview

Convert between different units of measurement. Supports length, weight, volume,
temperature, area, and speed conversions.

## Endpoint

### Convert Units

**Endpoint:** `GET /v1/technology/convert`

Convert a value from one unit to another.

#### Query Parameters

| Parameter | Type   | Required | Description                |
| --------- | ------ | -------- | -------------------------- |
| `from`    | string | Yes      | Source unit (e.g. `miles`) |
| `to`      | string | Yes      | Target unit (e.g. `km`)    |
| `value`   | number | Yes      | Numeric value to convert   |

#### Example Request

```
GET /v1/technology/convert?from=miles&to=km&value=10
```

#### Example Response

```json
{
  "from": "miles",
  "to": "km",
  "input": 10,
  "result": 16.09344,
  "formula": "miles × 1.609344"
}
```

#### Response Fields

| Field     | Type   | Description                           |
| --------- | ------ | ------------------------------------- |
| `from`    | string | Source unit                           |
| `to`      | string | Target unit                           |
| `input`   | number | Original value                        |
| `result`  | number | Converted value (rounded to 6 places) |
| `formula` | string | Human-readable conversion formula     |

## Supported Units

### Length

| Unit          | Key     |
| ------------- | ------- |
| Millimeter    | `mm`    |
| Centimeter    | `cm`    |
| Meter         | `m`     |
| Kilometer     | `km`    |
| Inch          | `in`    |
| Foot          | `ft`    |
| Yard          | `yd`    |
| Mile          | `miles` |
| Nautical Mile | `nmi`   |

### Weight

| Unit       | Key     |
| ---------- | ------- |
| Milligram  | `mg`    |
| Gram       | `g`     |
| Kilogram   | `kg`    |
| Metric Ton | `t`     |
| Ounce      | `oz`    |
| Pound      | `lb`    |
| Stone      | `stone` |

### Volume

| Unit        | Key     |
| ----------- | ------- |
| Milliliter  | `ml`    |
| Liter       | `l`     |
| Teaspoon    | `tsp`   |
| Tablespoon  | `tbsp`  |
| Fluid Ounce | `fl_oz` |
| Cup         | `cup`   |
| Pint        | `pt`    |
| Quart       | `qt`    |
| Gallon      | `gal`   |

### Temperature

| Unit       | Key |
| ---------- | --- |
| Celsius    | `c` |
| Fahrenheit | `f` |
| Kelvin     | `k` |

### Area

| Unit              | Key    |
| ----------------- | ------ |
| Square Millimeter | `mm2`  |
| Square Centimeter | `cm2`  |
| Square Meter      | `m2`   |
| Square Kilometer  | `km2`  |
| Square Inch       | `in2`  |
| Square Foot       | `ft2`  |
| Square Yard       | `yd2`  |
| Acre              | `acre` |
| Hectare           | `ha`   |

### Speed

| Unit                | Key     |
| ------------------- | ------- |
| Meters per Second   | `m_s`   |
| Kilometers per Hour | `km_h`  |
| Miles per Hour      | `mph`   |
| Knots               | `knots` |

## Error Responses

| HTTP Status | Reason                                              |
| ----------- | --------------------------------------------------- |
| 400         | Missing required parameters                         |
| 400         | `value` is not a valid number                       |
| 400         | Unknown unit key                                    |
| 400         | `from` and `to` belong to different unit categories |
