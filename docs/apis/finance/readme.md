# Finance APIs

## Overview

Financial data endpoints for payment card validation, currency exchange, banking
reference data, and macroeconomic indicators.

## Endpoints

### [Exchange Rate](./exchange-rate.md) - ✅ MVP

Get live currency exchange rates and convert amounts between currencies.

- **Status:** mvp
- **Endpoints:** `GET /v1/finance/exchange`, `GET /v1/finance/convert`
- **Credit Cost:** 1

### [Crypto Prices](./crypto-price.md) - ✅ MVP

Get live cryptocurrency prices, 24h change, market cap, and volume for 20+ major
coins.

- **Status:** mvp
- **Endpoint:** `GET /v1/finance/crypto`
- **Credit Cost:** 1

### [BIN Lookup](./bin.md) - ✅ MVP

Look up issuing bank, card network, type, and country from the first 6–8 digits
of a payment card.

- **Status:** mvp
- **Endpoint:** `GET /v1/finance/bin/{bin}`
- **Credit Cost:** 1

### [IBAN Validator](./iban.md) - ✅ MVP

Validate IBAN numbers and extract bank code and account number.

- **Status:** mvp
- **Endpoint:** `GET /v1/finance/iban/{iban}`
- **Credit Cost:** 1

### [SWIFT Code](./swift-code.md) - ✅ MVP

Validate and look up bank metadata by SWIFT/BIC code.

- **Status:** mvp
- **Endpoints:** `GET /v1/finance/swift/{code}`, `GET /v1/finance/swift`,
  `GET /v1/finance/swift/search`
- **Credit Cost:** 1

### [Mortgage Calculator](./mortgage.md) - ✅ MVP

Calculate monthly mortgage payments and full amortization schedules.

- **Status:** mvp
- **Endpoint:** `POST /v1/finance/mortgage`
- **Credit Cost:** 1

### [Commodity Prices](./commodities.md) - ✅ MVP

Historical and current annual average prices for 16 major commodities sourced
from FRED.

- **Status:** mvp
- **Endpoint:** `GET /v1/finance/commodities`
- **Credit Cost:** 1

### [Inflation](./inflation.md) - ✅ MVP

Historical and current CPI inflation rates for 241 countries, sourced from the
World Bank.

- **Status:** mvp
- **Endpoint:** `GET /v1/finance/inflation`
- **Credit Cost:** 1

## Category Statistics

- Total Endpoints: 8
- Live: 8
