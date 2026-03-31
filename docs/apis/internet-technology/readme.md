# Internet / Technology APIs

## Overview

Technology and internet-focused endpoints for domain validation, network
lookups, code generation, and web utilities.

## Endpoints

### [Barcode](./barcode.md) - ⏳ Planned

Generate and validate barcodes

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/barcode`

### [Base64 Encode / Decode](./base64.md) - ✅ MVP

Encode and decode Base64 strings

- **Status:** mvp
- **Endpoints:** `POST /v1/convert/base64/encode`, `POST /v1/convert/base64/decode`
- **Credit Cost:** 1

### [Disposable Email Checker](./disposable-email-checker.md) - ⏳ Planned

Check if an email is from a disposable email service

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/disposable-email-checker`

### [DNS Lookup](./dns-lookup.md) - ⏳ Planned

Perform DNS lookups for domains

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/dns-lookup`

### [Domain](./domain.md) - ⏳ Planned

Get domain information and availability

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/domain`

### [ASN Lookup](./asn-lookup.md) - ✅ MVP

Look up Autonomous System Number (ASN), organization, ISP, and network route information for any IP address.

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/ip/asn/{ip}`

### [IP Geolocation](./ip-lookup.md) - ✅ MVP

Get IP geolocation data including country, city, ISP, and VPN detection.

- **Status:** mvp
- **Endpoints:** `GET /v1/tech/ip` (caller IP), `GET /v1/tech/ip/{ip}` (specific IP)

### [MX Lookup](./mx-lookup.md) - ⏳ Planned

Perform MX record lookups

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/mx-lookup`

### [Password Generator](./password-generator.md) - ✅ MVP

Generate secure passwords

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/password`

### [QR Code](./qr-code.md) - ✅ MVP

Generate QR codes

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/qr`

### [URL Lookup](./url-lookup.md) - ⏳ Planned

Get information about URLs

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/url-lookup`

### [User Agent](./user-agent.md) - ✅ MVP

Parse and analyze user agent strings

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/useragent`

### [Validate Phone](./validate-phone.md) - ✅ MVP

Validate phone numbers

- **Status:** mvp
- **Endpoint:** `GET /v1/tech/validate/phone`

### [Validate Email](./validate-email.md) - ⏳ Planned

Validate email addresses

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/validate-email`

### [Webpage](./webpage.md) - ⏳ Planned

Get webpage information and metadata

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/webpage`

### [Web Scraper](./web-scraper.md) - ⏳ Planned

Scrape web pages for content

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/web-scraper`

### [Whois](./whois.md) - ⏳ Planned

Get WHOIS information for domains

- **Status:** planned
- **Planned Endpoint:** `GET /v1/internet-technology/whois`

## Category Statistics

- Total Endpoints: 18
- Live: 7
- Planned: 11
