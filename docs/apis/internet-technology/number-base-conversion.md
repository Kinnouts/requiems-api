# Number Base Conversion API

Convert integers between binary, octal, decimal, and hexadecimal.

## Endpoint

`GET /v1/convert/base`

## Query Parameters

| Parameter | Type    | Required | Description                                                                                                               |
| --------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------- |
| `from`    | integer | ✅       | Source base — must be one of `2`, `8`, `10`, `16`                                                                         |
| `to`      | integer | ✅       | Target base — must be one of `2`, `8`, `10`, `16`                                                                         |
| `value`   | string  | ✅       | Number as a string (e.g. `255`, `ff`, `11111111`). Accepts optional prefixes `0x`, `0b`, `0o` for their respective bases. |

## Supported Bases

| Value | Name        | Example input |
| ----- | ----------- | ------------- |
| `2`   | Binary      | `11111111`    |
| `8`   | Octal       | `377`         |
| `10`  | Decimal     | `255`         |
| `16`  | Hexadecimal | `ff`          |

## Response Envelope

All responses are wrapped in the standard envelope:

```json
{
  "data": {
    "input": "255",
    "from": 10,
    "to": 16,
    "result": "ff"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Prefix Handling

The `value` parameter accepts optional base prefixes that are stripped before
parsing:

| Prefix | Base | Example              |
| ------ | ---- | -------------------- |
| `0x`   | 16   | `0xff` → `255`       |
| `0b`   | 2    | `0b11111111` → `255` |
| `0o`   | 8    | `0o377` → `255`      |

Negative values are supported: `-255` in decimal converts to `-ff` in hex.

## Error Codes

| Code          | Status | When                                                                          |
| ------------- | ------ | ----------------------------------------------------------------------------- |
| `bad_request` | 400    | Missing parameter, unsupported base, or value not parseable in the given base |
