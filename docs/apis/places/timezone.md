# Timezone API

## Status

✅ **MVP** - Basic timezone lookup by coordinates and city name

## Overview

Get timezone information for any location. This endpoint provides timezone data
for geographic coordinates (latitude/longitude) or city names worldwide.

## Live Endpoints

### Get Timezone Information

**Endpoint:** `GET /v1/places/timezone`

Get timezone information for a location by coordinates or city name.

#### Query Parameters

| Parameter | Type   | Required | Description                                               |
| --------- | ------ | -------- | --------------------------------------------------------- |
| `lat`     | float  | *        | Latitude (-90 to 90). Required when using coordinates.    |
| `lon`     | float  | *        | Longitude (-180 to 180). Required when using coordinates. |
| `city`    | string | *        | City name. Required when not using coordinates.           |

\* Either `city` **or** both `lat` + `lon` must be provided.

#### Example Requests

```
GET /v1/places/timezone?lat=51.5&lon=-0.1
GET /v1/places/timezone?city=Tokyo
```

#### Example Response

```json
{
  "data": {
    "timezone": "Europe/London",
    "offset": "+00:00",
    "current_time": "2024-12-15T14:30:00Z",
    "is_dst": false
  },
  "metadata": {
    "timestamp": "2024-12-15T14:30:00Z"
  }
}
```

#### Response Fields

| Field          | Type    | Description                                       |
| -------------- | ------- | ------------------------------------------------- |
| `timezone`     | string  | IANA timezone identifier (e.g. `"Europe/London"`) |
| `offset`       | string  | UTC offset in `+HH:MM` / `-HH:MM` format          |
| `current_time` | string  | Current UTC time in RFC 3339 format               |
| `is_dst`       | boolean | Whether the location is currently observing DST   |
