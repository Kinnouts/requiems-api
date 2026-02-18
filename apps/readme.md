# Apps

- **api** - Go backend. Handles business logic and database queries. Receives
  requests from the auth gateway only, no auth of its own.

- **dashboard** - Rails web app. User registration, subscription management, API
  key management, and admin panel.

- **workers** - Cloudflare Workers. Contains `auth-gateway` (public edge,
  validates keys and proxies to Go), `api-management` (internal, manages API
  keys and usage data for Rails), and `shared` (common types and utilities).
