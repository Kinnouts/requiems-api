# Developer Tools APIs

## Overview

The boring-but-essential utilities every app eventually needs, already built.

## Endpoints

### [QR Code Generator](./qr-code.md) - ✅ MVP

Generate QR codes from any text or URL, returned as PNG or base64.

- **Status:** mvp
- **Endpoints:** `GET /v1/tech/qr`, `GET /v1/tech/qr/base64`
- **Credit Cost:** 1

### [Barcode Generator](./barcode.md) - ✅ MVP

Generate barcodes in multiple formats (Code 128, 93, 39, EAN-8, EAN-13).

- **Status:** mvp
- **Endpoints:** `GET /v1/tech/barcode`, `GET /v1/tech/barcode/base64`
- **Credit Cost:** 1

### [Password Generator](./password-generator.md) - ✅ MVP

Generate cryptographically secure random passwords.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/password`
- **Credit Cost:** 1

### [User Agent Parser](./user-agent.md) - ✅ MVP

Parse user agent strings: browser, OS, device type, and bot detection.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/useragent`
- **Credit Cost:** 1

### [Counter](./counter.md) - ✅ MVP

Atomic, namespace-isolated hit counter.

- **Status:** mvp
- **Endpoints:** `POST /v1/misc/counter/{namespace}`,
  `GET /v1/misc/counter/{namespace}`
- **Credit Cost:** 1

### [Random User](./random-user.md) - ✅ MVP

Generate random fake user profiles for testing and prototyping.

- **Status:** mvp
- **Endpoint:** `GET /v1/misc/random-user`
- **Credit Cost:** 1

### [Base64 Encode / Decode](./base64.md) - ✅ MVP

Encode strings to Base64 and decode Base64 back to plain text.

- **Status:** mvp
- **Endpoints:** `POST /v1/convert/base64/encode`,
  `POST /v1/convert/base64/decode`
- **Credit Cost:** 1

### [Number Base Conversion](./number-base-conversion.md) - ✅ MVP

Convert integers between binary, octal, decimal, and hexadecimal.

- **Status:** mvp
- **Endpoint:** `GET /v1/convert/base`
- **Credit Cost:** 1

### [Color Format Conversion](./color-conversion.md) - ✅ MVP

Convert color values between HEX, RGB, HSL, and CMYK.

- **Status:** mvp
- **Endpoint:** `GET /v1/convert/color`
- **Credit Cost:** 1

### [Unit Conversion](./unit-conversion.md) - ✅ MVP

Convert between units — length, weight, volume, temperature, area, speed.

- **Status:** mvp
- **Endpoint:** `GET /v1/misc/convert`
- **Credit Cost:** 1

### [Data Format Conversion](./data-format-conversion.md) - ✅ MVP

Convert structured data between JSON, YAML, CSV, XML, and TOML.

- **Status:** mvp
- **Endpoint:** `POST /v1/convert/data`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 11
- Live: 11
