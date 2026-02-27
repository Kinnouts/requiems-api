# The Ruby on Rails App

## Overview

The Rails application serves as the public-facing dashboard and admin panel for
Requiem API, handling user authentication, API key management, usage statistics,
and billing.

**Location:** [apps/dashboard/](../apps/dashboard/) **Technology:** Ruby on
Rails 8.1, Tailwind CSS, Hotwire (Turbo + Stimulus) **Database:** PostgreSQL
(shared with Go backend) **Background Jobs:** Sidekiq with Redis

## Responsibilities

### Public Website

- Landing page ([requiems.xyz](https://requiems.xyz))
- Marketing content and feature showcase
- Documentation portal
- Live API playground

### User Dashboard

- User authentication and registration (Devise)
- API key generation and management
- Usage statistics and analytics
- Billing and subscription management
- Account settings

### Admin Panel

- User management
- API key oversight
- Usage monitoring
- Abuse detection and reporting
- System health monitoring

## Architecture

```
┌─────────────────────────────────────┐
│         Rails App (Port 3000)       │
│                                     │
│  ┌──────────────────────────────┐   │
│  │     Controllers              │   │
│  │  - HomeController            │   │
│  │  - DashboardController       │   │
│  │  - ApiKeysController         │   │
│  └──────────────────────────────┘   │
│                                     │
│  ┌──────────────────────────────┐   │
│  │     Models                   │   │
│  │  - User (Devise)             │   │
│  │  - ApiKey                    │   │
│  │  - Subscription              │   │
│  │  - UsageLog                  │   │
│  └──────────────────────────────┘   │
│                                     │
│  ┌──────────────────────────────┐   │
│  │     Services                 │   │
│  │  - CloudflareKvSyncService   │   │
│  │  - UsageSyncService          │   │
│  └──────────────────────────────┘   │
│                                     │
│  ┌──────────────────────────────┐   │
│  │     Background Jobs          │   │
│  │  - SyncUsageJob              │   │
│  │  - CleanupLogsJob            │   │
│  └──────────────────────────────┘   │
└─────────────────────────────────────┘
         ↓                    ↓
    PostgreSQL             Redis
```

## Database Tables

### User Management

- `users` - User accounts (email, password, plan, etc.)
- `api_keys` - API keys associated with users
- `subscriptions` - Billing and subscription data

### Usage Tracking

- `usage_logs` - Detailed API usage records
- `daily_usage_summaries` - Aggregated daily usage
- `credit_adjustments` - Manual credit additions/deductions

### Administration

- `audit_logs` - System activity audit trail
- `abuse_reports` - Flagged suspicious activity

## Key Features

### API Key Management

Users can:

- Generate new API keys
- View key statistics (requests, credits used)
- Regenerate keys
- Delete keys
- Set key permissions (future feature)

When a key is created, updated, or deleted, it's managed through the API
Management worker via
[CloudflareApiManagementService](../apps/dashboard/app/services/cloudflare/api_management_service.rb).

### Usage Analytics

- Real-time usage tracking
- Historical usage charts
- Credit consumption breakdown
- Top endpoints by usage
- Monthly/daily aggregations

### Background Jobs (Solid Queue)

**[SyncD1UsageJob](../apps/dashboard/app/jobs/sync_d1_usage_job.rb)**

- Pulls usage data from Cloudflare D1 into PostgreSQL
- Runs every 5 minutes via Solid Queue recurring tasks

**[AggregateDailyUsageJob](../apps/dashboard/app/jobs/aggregate_daily_usage_job.rb)**

- Builds daily usage summaries for analytics/reporting
- Runs once per day at 00:05 UTC via Solid Queue recurring tasks

## Development

### Local Setup with Docker

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

The Rails app runs on http://localhost:3000 with:

- Hot reloading (changes reflected immediately)
- Tailwind CSS compilation
- Sidekiq worker running
- PostgreSQL and Redis available

### Running Commands

```bash
# Rails console
docker compose -f docker-compose.dev.yml exec dashboard rails console

# Database migrations
docker compose -f docker-compose.dev.yml exec dashboard rails db:migrate

# Run tests
docker compose -f docker-compose.dev.yml exec dashboard rails test

# Create admin user
docker compose -f docker-compose.dev.yml exec dashboard rails runner "
User.create!(
  name: 'Admin',
  email: 'admin@requiems.xyz',
  password: 'password123',
  password_confirmation: 'password123',
  admin: true,
  confirmed_at: Time.current
)
"
```

### Environment Variables

Required in production:

- `RAILS_ENV` - `production`
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `SECRET_KEY_BASE` - Rails secret (generate with `rails secret`)
- `RAILS_MASTER_KEY` - Encrypted credentials key
- `CLOUDFLARE_ACCOUNT_ID` - Cloudflare account ID
- `CLOUDFLARE_KV_NAMESPACE_ID` - KV namespace for API keys
- `CLOUDFLARE_API_TOKEN` - Cloudflare API token with KV write permissions

## Database Migrations

Rails uses separate migration tracking from the Go backend to avoid conflicts:

**Go migrations table:** `schema_migrations` **Rails migrations table:**
`schema_migrations` (separate schema)

Both apps share the same PostgreSQL database but maintain independent migration
histories.

## Styling

- **Tailwind CSS** - Utility-first CSS framework
- **ViewComponent** - Reusable view components
- **Hotwire** - Modern SPA-like experience without JavaScript frameworks
  - **Turbo Drive** - Fast navigation
  - **Turbo Frames** - Partial page updates
  - **Stimulus** - JavaScript sprinkles

## Testing

```bash
# Run all tests
docker compose -f docker-compose.dev.yml exec dashboard rails test

# Run specific test
docker compose -f docker-compose.dev.yml exec dashboard rails test test/models/api_key_test.rb

# Run system tests (browser tests)
docker compose -f docker-compose.dev.yml exec dashboard rails test:system
```

## Deployment

### With Docker Compose (VPS)

```bash
cd infra/docker
docker compose up -d --build
```

## Cloudflare Integration

### API Management Service

API key operations (create, update, revoke) are sent to the API Management
worker via
[CloudflareApiManagementService](../apps/dashboard/app/services/cloudflare/api_management_service.rb),
which keeps KV in sync. Rails never writes to KV directly.

### Usage Data Flow

```
User makes API request
    ↓
Worker records usage in D1
    ↓
Sidekiq job pulls from D1 (hourly)
    ↓
Usage stored in PostgreSQL
    ↓
Displayed in Rails dashboard
```

## Performance Considerations

- **Database Indexes:** All foreign keys and frequently queried columns indexed
- **Caching:** Fragment caching for dashboard views
- **Background Jobs:** Long-running operations moved to Sidekiq
- **Asset Pipeline:** Tailwind CSS compiled and minified
- **CDN:** Static assets served via CDN in production

## Related Documentation

- [Architecture Overview](./architecture.md)
- [Auth Gateway Documentation](./auth-gateway.md)
- [Infrastructure Guide](./infrastructure.md)
