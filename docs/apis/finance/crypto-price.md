# Crypto Prices API

Get live cryptocurrency prices, 24-hour change, market cap, and trading volume.
Prices are sourced from CoinGecko and cached for 5 minutes.

## Endpoint

`GET /v1/finance/crypto/{symbol}`

## Supported Symbols

| Symbol  | Name           |
| ------- | -------------- |
| `BTC`   | Bitcoin        |
| `ETH`   | Ethereum       |
| `BNB`   | BNB            |
| `XRP`   | XRP            |
| `ADA`   | Cardano        |
| `SOL`   | Solana         |
| `DOGE`  | Dogecoin       |
| `DOT`   | Polkadot       |
| `MATIC` | Polygon        |
| `AVAX`  | Avalanche      |
| `LINK`  | Chainlink      |
| `LTC`   | Litecoin       |
| `UNI`   | Uniswap        |
| `ATOM`  | Cosmos         |
| `TRX`   | TRON           |
| `XLM`   | Stellar        |
| `ALGO`  | Algorand       |
| `NEAR`  | NEAR Protocol  |
| `FTM`   | Fantom         |
| `SHIB`  | Shiba Inu      |

Symbols are case-insensitive — `btc` and `BTC` both work.

## Response Envelope

```json
{
  "data": {
    "symbol": "BTC",
    "name": "Bitcoin",
    "price_usd": 42000.50,
    "change_24h": 2.5,
    "market_cap": 820000000000,
    "volume_24h": 25000000000
  },
  "metadata": {
    "timestamp": "2026-01-01T00:00:00Z"
  }
}
```

## Error Codes

| Code             | Status | When                                      |
| ---------------- | ------ | ----------------------------------------- |
| `unknown_symbol` | 422    | Symbol not in the supported coin list     |
| `upstream_error` | 503    | CoinGecko unavailable or unexpected error |

## Caching

Responses are cached in Redis with a 5-minute TTL. The `metadata.timestamp`
reflects when the response was served, not when prices were last fetched from
upstream.
