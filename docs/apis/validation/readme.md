# Validation APIs

## Overview

Stop bad data at the door. These APIs validate emails, phone numbers, and text
content at the point of entry.

## Endpoints

### [Email Validator](./email-validate.md) - ✅ MVP

Full email validation: syntax check, MX record lookup, disposable domain
detection, normalization, and typo suggestions.

- **Status:** mvp
- **Endpoint:** `POST /v1/validation/email`
- **Credit Cost:** 1

### [Phone Validation](./phone-validation.md) - ✅ MVP

Validate phone numbers globally. Detect carrier, country, type, and VOIP risk.

- **Status:** mvp
- **Endpoints:** `GET /v1/validation/phone`, `POST /v1/validation/phone`
- **Credit Cost:** 1

### [Profanity Filter](./profanity.md) - ✅ MVP

Detect and censor profanity in text for content moderation.

- **Status:** mvp
- **Endpoint:** `POST /v1/text/profanity`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 3
- Live: 3
