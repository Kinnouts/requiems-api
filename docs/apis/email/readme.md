# Email APIs

> **Note:** These APIs are now grouped under **Validation** and **Networking &
> Internet** in the public catalog. See `apps/dashboard/config/api_catalog.yml`.

## Endpoints

### [Email Validator](./validate.md) - ✅ MVP

Full email validation: syntax, MX records, disposable detection, normalization,
typo suggestions.

- **Status:** mvp
- **Endpoint:** `POST /v1/email/validate`
- **Credit Cost:** 1

### [Disposable Domain Checker](./disposable.md) - ✅ MVP

Check whether an email domain belongs to a known disposable/temporary provider.

- **Status:** mvp
- **Endpoints:**
  - `POST /v1/email/disposable/check`
  - `POST /v1/email/disposable/check-batch`
  - `GET /v1/email/disposable/domain/{domain}`
  - `GET /v1/email/disposable/domains`
  - `GET /v1/email/disposable/stats`
- **Credit Cost:** 1

### [Email Normalizer](./normalize.md) - ✅ MVP

Normalize email addresses to canonical form with provider-specific rules.

- **Status:** mvp
- **Endpoint:** `POST /v1/email/normalize`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 3
- Live: 3
