# Color Format Conversion API

## Status

✅ **Live** - Available now at `GET /v1/technology/color`

## Overview

Convert color values between the four most common formats: HEX, RGB, HSL, and
CMYK. Useful for design tools, CSS utilities, palette generators, and any
application that handles color inputs from different sources.

## Endpoint

### Convert Color

**Endpoint:** `GET /v1/technology/color`

#### Query Parameters

| Parameter | Type   | Required | Description                                           |
| --------- | ------ | -------- | ----------------------------------------------------- |
| `from`    | string | Yes      | Source format: `hex`, `rgb`, `hsl`, or `cmyk`         |
| `to`      | string | Yes      | Target format: `hex`, `rgb`, `hsl`, or `cmyk`         |
| `value`   | string | Yes      | Color value in the source format (see examples below) |

#### Accepted Value Formats

| Format | Example input                          |
| ------ | -------------------------------------- |
| `hex`  | `#ff5733` or shorthand `#f53`          |
| `rgb`  | `rgb(255, 87, 51)` or `rgb(255,87,51)` |
| `hsl`  | `hsl(11, 100%, 60%)`                   |
| `cmyk` | `cmyk(0%, 66%, 80%, 0%)`               |

#### Example Request

```
GET /v1/technology/color?from=hex&to=hsl&value=%23ff5733
```

#### Example Response

```json
{
  "data": {
    "input": "#ff5733",
    "result": "hsl(11, 100%, 60%)",
    "formats": {
      "hex": "#ff5733",
      "rgb": "rgb(255, 87, 51)",
      "hsl": "hsl(11, 100%, 60%)",
      "cmyk": "cmyk(0%, 66%, 80%, 0%)"
    }
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

#### Response Fields

| Field          | Type   | Description                                        |
| -------------- | ------ | -------------------------------------------------- |
| `input`        | string | The original value passed in the `value` parameter |
| `result`       | string | The color expressed in the `to` format             |
| `formats.hex`  | string | HEX representation (`#rrggbb`)                     |
| `formats.rgb`  | string | RGB representation (`rgb(r, g, b)`)                |
| `formats.hsl`  | string | HSL representation (`hsl(h, s%, l%)`)              |
| `formats.cmyk` | string | CMYK representation (`cmyk(c%, m%, y%, k%)`)       |

## Error Responses

| HTTP Status | Code                | Reason                                            |
| ----------- | ------------------- | ------------------------------------------------- |
| 400         | `bad_request`       | Missing or invalid `from`, `to`, or `value`       |
| 422         | `invalid_color`     | The `value` cannot be parsed in the `from` format |
| 422         | `validation_failed` | `from` or `to` is not one of the allowed values   |

## Performance

Measured against production (`https://api.requiems.xyz`) with 50 samples.

| Metric  | Value   |
| ------- | ------- |
| p50     | 798 ms  |
| p95     | 935 ms  |
| p99     | 1263 ms |
| Average | 830 ms  |

_Last updated: 2026-04-16_
