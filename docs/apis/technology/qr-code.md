# QR Code API

## Status

✅ **MVP** - Live

## Overview

Generate QR codes from text or URLs. The endpoint supports PNG image output or
base64-encoded JSON response.

## Endpoint

### Generate QR Code

`GET /v1/tech/qr`

### Query Parameters

| Parameter | Required | Default | Description                                       |
| --------- | -------- | ------- | ------------------------------------------------- |
| `data`    | Yes      | —       | The text or URL to encode in the QR code          |
| `size`    | No       | `256`   | Image size in pixels (min: 50, max: 1000)         |
| `format`  | No       | `png`   | Response format: `png` (image) or `base64` (JSON) |

### Response

#### PNG format (default)

Returns a raw PNG image with `Content-Type: image/png`.

```
GET /v1/tech/qr?data=https://example.com&size=200
```

#### Base64 format

Returns a JSON envelope containing a base64-encoded PNG image.

```
GET /v1/tech/qr?data=https://example.com&size=200&format=base64
```

```json
{
  "data": {
    "image": "<base64-encoded PNG>",
    "width": 200,
    "height": 200
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```
