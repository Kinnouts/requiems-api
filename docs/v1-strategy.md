# V1

Currently there is NO clear "killer" value proposition.
Useful utilities spread across too many weak domains.
Developers don't pay for lorem ipsum or horoscopes.

---

## Two Teams, Two Tracks

### Team A — Validation

**Owner:** Engineer 1

Everything about knowing your data is clean and trustworthy — who someone is,
where they're coming from, whether their contact info is real.

**Sections:**

#### Email Intelligence

**Tagline:** _Know your emails before they bounce._

The highest willingness-to-pay in the API market. Email list hygiene has clear
ROI — fewer bounces, cleaner databases, lower spam rates. ZeroBounce and
NeverBounce charge $20–$500/month for this exact category.

What we already have:

- Disposable email detection (single, batch, domain, list, stats)

What we need to build:

- Full email validation — syntax check, MX record lookup, SMTP deliverability
- Email normalization — handle Gmail dots, plus-addressing, case folding
- Batch endpoint that combines disposable + deliverability in one call

Use cases: user registration cleanup, email marketing list hygiene, CRM data
quality, lead validation.

#### IP Intelligence

**Tagline:** _Turn any IP into actionable context._

Completely absent from our current roadmap but one of the most monetized API
categories in the market. ipinfo.io and ipstack.com built real businesses on
this. Every fraud detection system, geo-gating feature, and analytics dashboard
needs IP data.

What we already have: nothing.

What we need to build:

- IP geolocation — country, city, region, lat/lng, timezone
- VPN / proxy / Tor detection — fraud prevention signal
- ISP and ASN lookup — carrier, org, network info
- IPv4 and IPv6 support

Use cases: fraud detection, geo-gating content, personalizing UX by country,
blocking VPN abuse, analytics enrichment.

**Backend path:** `internal/tech/ip/` → `/v1/tech/ip/*`

---

### Team B — Finance

**Owner:** Engineer 2

Everything about money — prices, rates, and financial data that fintech apps,
trading tools, and e-commerce platforms need.

**Sections:**

#### Finance Data

**Tagline:** _Real financial data without the complexity._

Every fintech app, budgeting tool, and e-commerce platform needs market data or
financial lookups. These are tedious to source correctly — worth paying to skip.

What we already have: nothing.

What we need to build:

- Credit card BIN lookup — issuer, card type, country, brand
- Currency conversion — live exchange rates, multi-currency
- Crypto price — BTC, ETH, and major tokens (current price, 24h change)
- Commodity price — gold, silver, oil, natural gas
- Gold price — spot price in multiple currencies

Use cases: checkout flows, payment form autofill, crypto dashboards, trading
tools, personal finance apps, commodity tracking.

**Backend path:** `internal/finance/` (new) → `/v1/finance/*`

---

## Weekly Iteration Plan

Each week one engineer ships 2–3 endpoints from their track. No blocking
between tracks.

| Team A — Validation   | Team B — Finance            |
| --------------------- | --------------------------- |
| Email full validation | BIN lookup                  |
| Email normalization   | Currency conversion         |
| IP geolocation        | Crypto price (top 10 coins) |
| VPN / proxy detection | Commodity price             |
| ISP / ASN lookup      | Gold price                  |
| Full Bot Detector     | Mortgage calculator         |
