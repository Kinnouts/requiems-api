# Places APIs

## Overview

Location and geography-focused endpoints for timezones, calendars, geocoding,
and place data.

## Endpoints

### [Timezone](./timezone.md) - ✅ MVP

Get timezone information for any location by coordinates or city name.

- **Status:** mvp
- **Endpoint:** `GET /v1/places/timezone`
- **Credit Cost:** 1

### [World Time](./world-time.md) - ✅ MVP

Get the current time for any IANA timezone by name.

- **Status:** mvp
- **Endpoint:** `GET /v1/places/time/{timezone}`
- **Credit Cost:** 1

### [Working Days Calculator](./working-days.md) - ✅ MVP

Calculate the number of working days between two dates with optional
country-specific holidays.

- **Status:** mvp
- **Endpoint:** `GET /v1/places/working-days`
- **Credit Cost:** 1

### [Holidays](./holidays.md) - ✅ MVP

Get a list of public holidays for a specific country and year.

- **Status:** mvp
- **Endpoint:** `GET /v1/places/holidays`
- **Credit Cost:** 1

### [Geocoding](./geocode.md) - ✅ MVP

Convert addresses to coordinates and coordinates back to addresses.

- **Status:** mvp
- **Endpoints:** `GET /v1/places/geocode`, `GET /v1/places/reverse-geocode`
- **Credit Cost:** 1

### [Postal Code](./postal-code.md) - ✅ MVP

Look up city, state, and coordinates for any postal code worldwide.

- **Status:** mvp
- **Endpoint:** `GET /v1/places/postal-code/{code}`
- **Credit Cost:** 1

### [Cities](./cities.md) - ✅ MVP

Look up city metadata including population, timezone, and coordinates.

- **Status:** mvp
- **Endpoint:** `GET /v1/places/cities/{city}`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 7
- Live: 7
