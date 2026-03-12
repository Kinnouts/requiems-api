# Email Normalize API

Normalizes an email address to its canonical form using the
[go-email-normalizer](https://github.com/bobadilla-tech/go-email-normalizer)
package.

## Endpoint

`POST /v1/email/normalize`

## What Normalization Does

| Transformation | Example |
| --- | --- |
| Lowercase the full address | `User@Example.COM` → `user@example.com` |
| Trim leading/trailing whitespace | `  user@example.com  ` → `user@example.com` |
| Remove dots from Gmail local part | `te.st.user@gmail.com` → `testuser@gmail.com` |
| Strip plus tag (supported providers) | `user+spam@gmail.com` → `user@gmail.com` |
| Resolve alias domains | `user@googlemail.com` → `user@gmail.com` |

Normalizations are applied in order. The `changes` array in the response lists
every transformation that was applied.

## Request

```json
{
  "email": "Te.st.User+spam@Googlemail.com"
}
```

**Validation:** `email` is required and must be a syntactically valid address.
Invalid format returns `422 Unprocessable Entity`.

## Response Envelope

All responses are wrapped in the standard envelope:

```json
{
  "data": { ... },
  "metadata": { "timestamp": "2026-01-01T00:00:00Z" }
}
```

## Response Fields

| Field | Type | Description |
| --- | --- | --- |
| `original` | string | Exact input supplied by the caller |
| `normalized` | string | Canonical address after all transformations |
| `local` | string | Local part (before `@`) of the normalized address |
| `domain` | string | Domain part (after `@`) of the normalized address |
| `changes` | array | Ordered list of transformations applied (empty when none) |

### `changes` Values

| Value | When applied |
| --- | --- |
| `trimmed_whitespace` | Leading or trailing whitespace was present |
| `removed_trailing_dot` | One or more trailing dots were stripped from the raw input |
| `lowercase` | Any uppercase characters were found |
| `removed_dots` | Dots were removed from the local part (e.g. Gmail, Protonmail) |
| `removed_underscores` | Underscores were removed from the local part (Protonmail) |
| `removed_hyphens` | Hyphens were removed from the local part (Protonmail) |
| `replaced_hyphens_with_dots` | Hyphens in the local part were replaced with dots (Yandex) |
| `removed_plus_tag` | A plus-sign subaddress (`+tag`) was stripped (e.g. Gmail, Apple, Fastmail) |
| `removed_plus_signs` | All plus signs were removed regardless of position (e.g. Microsoft, Yahoo, Zoho) |
| `removed_subaddress` | A dash-delimited subaddress (`-tag`) was stripped (Yahoo) |
| `canonicalized_domain` | An alias domain was resolved (e.g. googlemail.com → gmail.com) |

## Example

```bash
curl -X POST https://api.requiems.xyz/v1/email/normalize \
  -H "requiems-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"email": "Te.st.User+spam@Googlemail.com"}'
```

Response:

```json
{
  "data": {
    "original": "Te.st.User+spam@Googlemail.com",
    "normalized": "testuser@gmail.com",
    "local": "testuser",
    "domain": "gmail.com",
    "changes": [
      "lowercase",
      "removed_dots",
      "removed_plus_tag",
      "canonicalized_domain"
    ]
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Error Codes

| Code | Status | When |
| --- | --- | --- |
| `validation_failed` | 422 | `email` field is missing or not a valid address |
| `bad_request` | 400 | Body is missing, invalid JSON, or contains unknown fields |
| `internal_error` | 500 | Unexpected server error |

## Use Cases

### De-duplicate accounts at registration

Store both the raw input and the normalized form. On new signups, query the
normalized address to catch users registering with subtle variations of an
address they already own.

```javascript
const { data } = await normalizeEmail(userInput);
const existing = await db.users.findByNormalizedEmail(data.normalized);
if (existing) {
  throw new ConflictError("An account with this email already exists");
}
await db.users.create({ email: userInput, email_normalized: data.normalized });
```

### Detect plus-tag abuse

If `changes` includes `removed_plus_tag`, the user is using a plus tag — often
a signal of throwaway signup attempts.

```javascript
const { data } = await normalizeEmail(userInput);
if (data.changes.includes("removed_plus_tag")) {
  flagForReview(userInput);
}
```
