# Miscellaneous APIs

> **Note:** Live APIs from this directory are now grouped under **Developer
> Tools** in the public catalog. See `apps/dashboard/config/api_catalog.yml`.

## Live Endpoints

### [Counter](./counter.md) - ✅ MVP

Atomic, namespace-isolated hit counter.

- **Status:** mvp
- **Endpoints:** `POST /v1/misc/counter/{namespace}`,
  `GET /v1/misc/counter/{namespace}`
- **Credit Cost:** 1

### [Unit Conversion](./unit-conversion.md) - ✅ MVP

Convert between units of measurement — length, weight, volume, temperature,
area, speed.

- **Status:** mvp
- **Endpoint:** `GET /v1/misc/convert`
- **Credit Cost:** 1

### [Random User](./random-user.md) - ✅ MVP

Generate random fake user profiles for testing and prototyping.

- **Status:** mvp
- **Endpoint:** `GET /v1/misc/random-user`
- **Credit Cost:** 1

### [Color Format Conversion](./color-conversion.md) - ✅ MVP

Convert color values between HEX, RGB, HSL, and CMYK.

- **Status:** mvp
- **Endpoint:** `GET /v1/convert/color`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 4
- Live: 4
