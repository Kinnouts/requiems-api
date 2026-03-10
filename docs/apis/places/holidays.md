# Holidays API

## Status

✅ **MVP** - Basic holiday lookup by country and year

## Overview

Get public holiday information for any country and year. This endpoint provides
a list of national holidays with their dates.

## Live Endpoints

### Get Holidays

**Endpoint:** `GET /v1/places/holidays`

Get public holidays for a specific country and year.

#### Query Parameters

| Parameter | Type   | Required | Description                                    |
| --------- | ------ | -------- | ---------------------------------------------- |
| `country` | string | Yes      | ISO 3166-1 alpha-2 country code (e.g., `US`) |
| `year`    | int    | Yes      | Year (e.g., `2025`)                            |

#### Example Requests

```
GET /v1/places/holidays?country=US&year=2025
GET /v1/places/holidays?country=GB&year=2025
```

#### Example Response

```json
{
  "data": {
    "country": "US",
    "year": 2025,
    "holidays": [
      {
        "date": "2025-01-01",
        "name": "New Year's Day"
      },
      {
        "date": "2025-07-04",
        "name": "Independence Day"
      }
    ]
  },
  "metadata": {
    "timestamp": "2025-01-15T10:30:00Z"
  }
}
```

#### Response Fields

| Field     | Type           | Description                                  |
| --------- | -------------- | -------------------------------------------- |
| `country` | string         | ISO 3166-1 alpha-2 country code             |
| `year`    | int            | Year                                         |
| `holidays` | array of objects | List of holidays with date and name       |

#### Holiday Object Fields

| Field | Type   | Description              |
| ------| ------ | ------------------------ |
| `date`| string | Date in `YYYY-MM-DD` format |
| `name`| string | Name of the holiday      |
