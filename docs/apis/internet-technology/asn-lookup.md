# ASN Lookup API

Look up Autonomous System Number (ASN), organization, ISP, and network route information for any IP address.

## Endpoint

`GET /v1/tech/ip/asn/{ip}`

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `ip` | string | No | IP address to look up (IPv4 or IPv6). If omitted, uses the requester's IP. |

## Response Envelope

All responses are wrapped in the standard envelope:

```json
{
  "data": {
    "ip": "8.8.8.8",
    "asn": "AS15169",
    "org": "Google LLC",
    "isp": "Google Public DNS",
    "domain": "google.com",
    "route": "8.8.8.0/24",
    "type": "hosting"
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `ip` | string | The IP address that was looked up |
| `asn` | string | Autonomous System Number (e.g., "AS15169") |
| `org` | string | Organization name owning the IP range |
| `isp` | string | Internet Service Provider |
| `domain` | string | Domain name associated with the IP |
| `route` | string | CIDR notation of the network route |
| `type` | string | Type of network (hosting, isp, business, cdn) |

## Error Codes

| Code | Status | When |
|------|--------|------|
| `bad_request` | 400 | Invalid IP address |
| `internal_error` | 500 | Unexpected failure |

## Examples

### cURL

```bash
curl "https://api.requiems.xyz/v1/tech/ip/asn/8.8.8.8" \
  -H "requiems-api-key: YOUR_API_KEY"
```

### Python

```python
import requests

url = "https://api.requiems.xyz/v1/tech/ip/asn/8.8.8.8"
headers = {"requiems-api-key": "YOUR_API_KEY"}

response = requests.get(url, headers=headers)
print(response.json())
```

## FAQ

**What is an ASN?**

An Autonomous System Number (ASN) is a unique identifier assigned to a group of IP networks and routers that operate under a common administration.

**Does this support IPv6?**

Yes, both IPv4 and IPv6 addresses are fully supported.

**What happens with private IP addresses?**

Private IP addresses (192.168.x.x, 10.x.x.x, 172.16-31.x.x) do not have ASN information. The API returns the IP with empty ASN fields.
