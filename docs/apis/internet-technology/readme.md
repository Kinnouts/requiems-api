# Internet / Technology APIs

## Overview

Network intelligence, domain lookups, developer utilities, and code generation
tools. This directory covers both the **Networking & Internet** and **Developer
Tools** categories from the public API catalog.

## Networking & Internet

### [IP Geolocation](./ip-lookup.md) - ✅ MVP

Get geolocation data for any IP address including country, city, ISP, and VPN
detection.

- **Status:** mvp
- **Endpoints:** `GET /v1/tech/ip` (caller IP), `GET /v1/tech/ip/{ip}`
- **Credit Cost:** 1

### [ASN Lookup](./asn-lookup.md) - ✅ MVP

Look up ASN, organization, ISP, and network route information for any IP
address.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/ip/asn/{ip}`
- **Credit Cost:** 1

### [VPN & Proxy Detection](./vpn-detection.md) - ✅ MVP

Detect if an IP belongs to a VPN, proxy, Tor exit node, or hosting provider.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/ip/vpn/{ip}`
- **Credit Cost:** 1

### [WHOIS Lookup](./whois.md) - ✅ MVP

Get domain registration details including registrar, name servers, and expiry
dates.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/whois/{domain}`
- **Credit Cost:** 1

### [Domain Info](./domain.md) - ✅ MVP

Look up DNS records and check domain availability.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/domain/{domain}`
- **Credit Cost:** 1

### [MX Lookup](./mx-lookup.md) - ✅ MVP

Look up MX records for any domain, sorted by priority.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/mx/{domain}`
- **Credit Cost:** 1

## Developer Tools

### [QR Code Generator](./qr-code.md) - ✅ MVP

Generate QR codes from any text or URL, returned as PNG or base64.

- **Status:** mvp
- **Endpoints:** `GET /v1/tech/qr`, `GET /v1/tech/qr/base64`
- **Credit Cost:** 1

### [Barcode Generator](./barcode.md) - ✅ MVP

Generate barcodes in multiple formats, returned as PNG or base64.

- **Status:** mvp
- **Endpoints:** `GET /v1/tech/barcode`, `GET /v1/tech/barcode/base64`
- **Credit Cost:** 1

### [Password Generator](./password-generator.md) - ✅ MVP

Generate cryptographically secure random passwords.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/password`
- **Credit Cost:** 1

### [User Agent Parser](./user-agent.md) - ✅ MVP

Parse user agent strings to extract browser, OS, device type, and bot detection.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/useragent`
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
- **Endpoint:** `GET /v1/tech/color`
- **Credit Cost:** 1

### [Data Format Conversion](./data-format-conversion.md) - ✅ MVP

Convert structured data between JSON, YAML, CSV, XML, and TOML.

- **Status:** mvp
- **Endpoint:** `POST /v1/convert/data`
- **Credit Cost:** 1

## Category Statistics

- Networking & Internet: 6 live
- Developer Tools: 8 live
- **Total: 14 live**
