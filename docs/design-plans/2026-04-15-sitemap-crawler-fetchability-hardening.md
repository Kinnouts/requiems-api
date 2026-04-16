# Sitemap & Crawler Fetchability Hardening

Google Search Console reported "Couldn't fetch" / "Unknown" type for
`https://requiems.xyz/sitemap.xml` while direct `curl` requests returned valid
XML with `application/xml; charset=utf-8`. Chrome rendered the sitemap as plain
text rather than the XML tree viewer it shows for other sitemaps.

---

## Problem Statement

Three distinct bugs were found through investigation:

### Bug 1 — Rails appends `; charset=utf-8` to `application/xml`

When `render ..., content_type: "application/xml"` is called, Rails
automatically appends `; charset=utf-8`, producing
`Content-Type: application/xml; charset=utf-8`.

RFC 7303 §6.1 states the charset parameter **SHOULD NOT** be used with
`application/xml`; the encoding is determined from the XML declaration instead.
Chrome's XML tree viewer treats `application/xml; charset=utf-8` differently
from plain `application/xml`, rendering the latter correctly and the former as
plain text. Comparison against a working sitemap confirmed the difference.

### Bug 2 — Caddy streaming gzip yields `Content-Length: 0` on HEAD requests

Caddy's `encode gzip` compresses response bodies as a streaming pass-through.
For HEAD requests there is no body to compress, so Caddy cannot compute the
compressed content-length and instead emits `Content-Length: 0`.

Google Search Console (and some other crawlers) issue a HEAD request before
fetching to check file metadata. Seeing `Content-Length: 0` caused GSC to report
"Unknown" type and "Couldn't fetch".

### Bug 3 — `respond_to { |f| f.xml }` leaves content-type resolution to MIME

negotiation

The original controller used `before_action :set_response_content_type` + a bare
`respond_to` block. Under some proxy/CDN configurations the MIME negotiation
path can override the before-action value. This was replaced with explicit
`render` calls that own format, layout, and content-type directly.

---

## Goals

1. Serve `sitemap.xml` with exactly `Content-Type: application/xml` (no charset
   parameter).
2. Ensure HEAD requests to sitemap/llms paths return a real `Content-Length`.
3. Make all sitemap/llms responses explicitly typed and publicly cacheable from
   the Rails layer.
4. Validate behavior with controller-level tests.

---

## Non-Goals

1. Changing sitemap XML structure or content generation logic.
2. Introducing long-lived cache durations.

---

## Design

### 1. Explicit `render` with direct Content-Type header override (Rails)

Replaced `respond_to { |f| f.xml/f.text }` + `before_action` with explicit
`render` calls that specify `formats:`, `layout: false`, and `content_type:`.

For the sitemap action, `response.headers["Content-Type"]` is set to
`"application/xml"` **after** render to strip the `; charset=utf-8` that Rails
auto-appends:

```ruby
def sitemap
  expires_in 5.minutes, public: true
  @apis = live_apis
  @last_modified = Time.current.beginning_of_day
  render "sitemap/sitemap", formats: [ :xml ], layout: false, content_type: "application/xml"
  response.headers["Content-Type"] = "application/xml"
end
```

Implementation:
[apps/dashboard/app/controllers/sitemap_controller.rb](../../apps/dashboard/app/controllers/sitemap_controller.rb)

### 2. Caddy matcher to skip gzip for sitemap/llms paths

Added a named `@compressible` matcher in the Caddyfile that excludes
`/sitemap.xml`, `/llms.txt`, and `/llms-full.txt` from gzip encoding:

```caddy
@compressible not path /sitemap.xml /llms.txt /llms-full.txt
encode @compressible gzip
```

This means these paths are served uncompressed. HEAD requests then return the
actual uncompressed `Content-Length` instead of `0`.

Implementation: [infra/caddy/Caddyfile](../../infra/caddy/Caddyfile)

---

## Validation

### Must-pass Rails checks

```bash
docker exec requiem-dev-dashboard-1 bin/rails test
docker exec requiem-dev-dashboard-1 bundle exec bundler-audit
docker exec requiem-dev-dashboard-1 bin/importmap audit
docker exec requiem-dev-dashboard-1 bundle exec brakeman --no-pager
```

### Targeted test

```bash
docker exec requiem-dev-dashboard-1 bin/rails test test/controllers/sitemap_controller_test.rb
```

### Manual verification after deploy

```bash
# Content-Type must be exactly "application/xml" — no charset
curl -sI https://requiems.xyz/sitemap.xml | grep content-type
# → content-type: application/xml

# HEAD must return a non-zero Content-Length
curl -sI https://requiems.xyz/sitemap.xml | grep content-length
# → content-length: <actual size>

# Body must start with XML declaration
curl -s https://requiems.xyz/sitemap.xml | head -c 40
# → <?xml version="1.0" encoding="UTF-8"?>
```

---

## Tradeoffs

### Why remove `; charset=utf-8`?

RFC 7303 §6.1 says charset SHOULD NOT be used with `application/xml`. Chrome's
XML viewer behaves differently for `application/xml` vs
`application/xml; charset=utf-8`. Omitting it matches RFC guidance and
real-world working sitemap behavior.

### Why disable gzip for sitemap paths instead of fixing Caddy globally?

Caddy's streaming gzip is by design — it avoids buffering large responses. The
`Content-Length: 0` on HEAD is a known trade-off. Disabling gzip only for the
small set of crawler-consumed static-ish paths is a targeted fix. The sitemap is
~54 KB uncompressed; the overhead is acceptable and the crawl reliability gain
is worth it.

### Why short public cache (300s)?

Helps crawler fetch stability and intermediary caching while keeping updates
near-real-time for a dynamic API catalog.

---

## Rollout Notes

1. Deploy Caddy config change (`infra/caddy/Caddyfile`).
2. Deploy dashboard changes (`apps/dashboard`).
3. Hard-refresh browser (Cmd+Shift+R) to bypass local cached response.
4. Re-submit sitemap in Google Search Console.
5. Use "Test Live URL" on `https://requiems.xyz/sitemap.xml` in GSC URL
   Inspection to trigger an immediate recrawl.
6. Monitor GSC sitemap status over 24–72 hours.

---

## Files Changed

1. [apps/dashboard/app/controllers/sitemap_controller.rb](../../apps/dashboard/app/controllers/sitemap_controller.rb)
2. [apps/dashboard/test/controllers/sitemap_controller_test.rb](../../apps/dashboard/test/controllers/sitemap_controller_test.rb)
3. [infra/caddy/Caddyfile](../../infra/caddy/Caddyfile)
