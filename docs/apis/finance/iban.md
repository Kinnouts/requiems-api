# IBAN Validator API

Validate IBAN numbers and extract the bank code and account number.
Supports all countries in the official SWIFT IBAN Registry (~80 countries).

## Endpoint

`GET /v1/finance/iban/{iban}`

## Parameters

| Name   | Type   | Location | Required | Description                                              |
| ------ | ------ | -------- | -------- | -------------------------------------------------------- |
| `iban` | string | path     | yes      | The IBAN to validate. Spaces stripped. Case-insensitive. |

## Response

Always HTTP `200`. Check the `valid` field to determine pass/fail.

```json
{
  "data": {
    "iban": "DE89370400440532013000",
    "valid": true,
    "country": "Germany",
    "bank_code": "37040044",
    "account": "0532013000"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

| Field       | Type    | Description                                                        |
| ----------- | ------- | ------------------------------------------------------------------ |
| `iban`      | string  | Normalised IBAN (spaces stripped, uppercased)                      |
| `valid`     | boolean | `true` if length and checksum passed                               |
| `country`   | string  | Full country name (empty if country code not in registry)          |
| `bank_code` | string  | Bank identifier from the BBAN (empty if positions not in registry) |
| `account`   | string  | Account number from the BBAN (empty if positions not in registry)  |

### Invalid IBAN example

```json
{
  "data": {
    "iban": "DE00370400440532013000",
    "valid": false,
    "country": "Germany",
    "bank_code": "",
    "account": ""
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Validation Algorithm

ISO 13616 mod-97 checksum:

1. Rearrange — move the first 4 characters to the end.
2. Replace each letter with its numeric equivalent (`A=10` … `Z=35`).
3. Compute the integer mod 97. A remainder of **1** means valid.

Length is also validated against the country's expected IBAN length from
the SWIFT IBAN Registry.

## Data Source

Country format data (expected IBAN length, bank code position, account
number position) is seeded from the
[php-iban project](https://github.com/globalcitizen/php-iban), a maintained
mirror of the official SWIFT IBAN Registry (ISO 13616).

Re-run `cmd/seed-iban` to refresh when new countries join or formats change.

## Error Codes

| Code             | Status | When                  |
| ---------------- | ------ | --------------------- |
| `internal_error` | 500    | Database unreachable  |
