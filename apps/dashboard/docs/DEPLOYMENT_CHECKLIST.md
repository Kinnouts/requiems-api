# Deployment Checklist

Use this checklist to ensure a smooth deployment of Requiems API Dashboard to
production.

## Pre-Deployment

### 1. Dependencies

```bash
cd apps/dashboard

# Install Ruby dependencies
bundle install

# Install JavaScript dependencies
npm install

# Verify versions
ruby -v  # Should be 3.3+
rails -v # Should be 8.1+
node -v  # Should be 18+
```

### 2. Run Tests

```bash
# Run all tests
bundle exec rails test

# Expected: 71+ tests passing
# 0 failures, 0 errors

# Run specific test suites
bundle exec rails test:models
bundle exec rails test:controllers
```

### 3. Database Setup

```bash
# Create databases
rails db:create

# Run migrations
rails db:migrate

# Verify schema
rails db:schema:load
```

### 4. Environment Variables

Create `.env` file with required variables:

```env
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/requiem_production

# Rails
RAILS_ENV=production
SECRET_KEY_BASE=generate_with_rails_secret
RAILS_LOG_LEVEL=info

# Redis
REDIS_URL=redis://localhost:6379/0

# Email (SMTP)
SMTP_ADDRESS=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your@email.com
SMTP_PASSWORD=your_app_password
SMTP_DOMAIN=yourdomain.com

# LemonSqueezy
LEMONSQUEEZY_API_KEY=your_api_key
LEMONSQUEEZY_STORE_SLUG=your_store
LEMONSQUEEZY_DEVELOPER_MONTHLY_VARIANT_ID=variant_id
LEMONSQUEEZY_DEVELOPER_YEARLY_VARIANT_ID=variant_id
LEMONSQUEEZY_BUSINESS_MONTHLY_VARIANT_ID=variant_id
LEMONSQUEEZY_BUSINESS_YEARLY_VARIANT_ID=variant_id
LEMONSQUEEZY_PROFESSIONAL_MONTHLY_VARIANT_ID=variant_id
LEMONSQUEEZY_PROFESSIONAL_YEARLY_VARIANT_ID=variant_id

# Go API Backend
BACKEND_API_URL=http://api-backend:8080
BACKEND_SECRET=generate_secure_secret

# Feature Flags (optional)
ENABLE_REGISTRATION=true
ENABLE_API_PLAYGROUND=true
RATE_LIMIT_ENABLED=true

# Monitoring (optional)
SENTRY_DSN=your_sentry_dsn
```

Generate secrets:

```bash
# Generate SECRET_KEY_BASE
bundle exec rails secret

# Generate BACKEND_SECRET
openssl rand -hex 32
```

### 5. Assets

```bash
# Precompile assets for production
RAILS_ENV=production bundle exec rails assets:precompile

# Verify assets compiled
ls -la public/assets
```

### 6. Security Check

- [ ] All secrets are in environment variables (not committed)
- [ ] `.env` is in `.gitignore`
- [ ] Database credentials are secure
- [ ] API keys are rotated from development
- [ ] SSL certificates are configured
- [ ] CORS is properly configured
- [ ] Rate limiting is enabled

---

## Deployment Methods

### Option 1: Docker (Recommended)

```bash
# Build production image
docker compose -f docker-compose.prod.yml build

# Start services
docker compose -f docker-compose.prod.yml up -d

# Run migrations
docker compose -f docker-compose.prod.yml exec dashboard rails db:migrate

# Create admin user
docker compose -f docker-compose.prod.yml exec dashboard rails console
# User.create!(email: 'admin@example.com', password: 'secure_password', admin: true)

# Check logs
docker compose -f docker-compose.prod.yml logs -f dashboard
```

### Option 2: Traditional Server

```bash
# On server
git clone https://github.com/bobadilla-tech/requiems-api.git
cd requiems-api/apps/dashboard

# Install dependencies
bundle install --deployment --without development test

# Setup database
RAILS_ENV=production bundle exec rails db:migrate

# Precompile assets
RAILS_ENV=production bundle exec rails assets:precompile

# Start server
bundle exec rails server -e production -p 3000
```

### Option 3: Heroku

```bash
# Install Heroku CLI
heroku login

# Create app
heroku create your-app-name

# Add PostgreSQL
heroku addons:create heroku-postgresql:essential-0

# Add Redis
heroku addons:create heroku-redis:mini

# Set environment variables
heroku config:set SECRET_KEY_BASE=$(rails secret)
heroku config:set LEMONSQUEEZY_API_KEY=your_key
# ... set all other env vars

# Deploy
git push heroku main

# Run migrations
heroku run rails db:migrate

# Create admin user
heroku run rails console
```

---

## Post-Deployment

### 1. Database Verification

```bash
# Check database connection
rails runner "puts User.count"

# Verify migrations
rails db:migrate:status

# Check indexes
rails runner "puts ActiveRecord::Base.connection.indexes(:usage_logs).map(&:name)"
```

### 2. Create Admin User

```bash
rails console

# In console:
User.create!(
  email: 'admin@yourdomain.com',
  password: ENV['ADMIN_PASSWORD'] || 'change_this_password',
  password_confirmation: ENV['ADMIN_PASSWORD'] || 'change_this_password',
  admin: true,
  name: 'Admin'
)
```

### 3. Smoke Tests

Visit these URLs and verify they load:

- [ ] `https://yourdomain.com` - Homepage
- [ ] `https://yourdomain.com/apis` - API Directory
- [ ] `https://yourdomain.com/pricing` - Pricing page
- [ ] `https://yourdomain.com/users/sign_in` - Login page
- [ ] `https://yourdomain.com/users/sign_up` - Registration
- [ ] `https://yourdomain.com/dashboard` - Dashboard (after login)
- [ ] `https://yourdomain.com/admin` - Admin panel (admin user)

### 4. Functional Tests

- [ ] User can register for new account
- [ ] User can log in
- [ ] User can create API key
- [ ] User can view usage analytics
- [ ] User can upgrade plan (test mode)
- [ ] Admin can view dashboard
- [ ] Admin can search users
- [ ] Admin can view analytics
- [ ] Email delivery works
- [ ] Password reset works

### 5. Performance Check

```bash
# Check response times
time curl https://yourdomain.com
time curl https://yourdomain.com/apis

# Should be < 500ms for cached pages
# Should be < 2s for dynamic pages

# Check database query performance
rails runner "
  require 'benchmark'
  puts Benchmark.measure {
    User.includes(:api_keys, :subscription).limit(100).to_a
  }
"
```

### 6. Monitoring Setup

```bash
# Check logs
tail -f log/production.log

# Monitor errors
# Configure Sentry/Rollbar/Bugsnag

# Set up uptime monitoring
# UptimeRobot, Pingdom, or StatusCake

# Configure APM
# New Relic, DataDog, or Scout APM
```

### 7. Background Jobs

```bash
# Verify Solid Queue is running
rails runner "puts SolidQueue::Job.count"

# Check recurring jobs
rails runner "puts SolidQueue::RecurringTask.all.map(&:key)"

# Expected recurring jobs:
# - sync_d1_usage (every 5 minutes)
# - aggregate_daily_usage (daily at 00:05 UTC)

# Start background workers (if not using Docker)
bundle exec rails solid_queue:start
```

### 8. Backups

```bash
# Set up automated database backups
# Cron job example:

# Daily backup at 2 AM
0 2 * * * pg_dump requiem_production | gzip > /backups/db_$(date +\%Y\%m\%d).sql.gz

# Weekly cleanup (keep last 30 days)
0 3 * * 0 find /backups -name "db_*.sql.gz" -mtime +30 -delete
```

### 9. SSL/TLS

```bash
# Verify SSL is working
curl -I https://yourdomain.com | grep -i "HTTP/2"

# Check SSL certificate
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com

# Renew Let's Encrypt (if using)
certbot renew --dry-run
```

### 10. DNS & CDN

- [ ] DNS A records point to server IP
- [ ] CNAME for www subdomain (if applicable)
- [ ] CDN configured (CloudFlare recommended)
- [ ] DNS propagation complete (`dig yourdomain.com`)

---

## Go-Live Checklist

### Before Going Live

- [ ] All tests passing
- [ ] Database migrations successful
- [ ] Environment variables configured
- [ ] SSL certificate installed
- [ ] Admin user created
- [ ] Email delivery tested
- [ ] Payment integration tested (LemonSqueezy)
- [ ] Backups configured
- [ ] Monitoring tools set up
- [ ] Error tracking configured
- [ ] Performance benchmarks met
- [ ] Security audit completed
- [ ] Load testing performed (if expecting high traffic)

### Launch Day

- [ ] Announce on social media
- [ ] Update GitHub README with live URL
- [ ] Monitor error logs closely
- [ ] Watch server resources (CPU, memory, disk)
- [ ] Monitor database connections
- [ ] Check email delivery
- [ ] Verify payment webhooks
- [ ] Have rollback plan ready

### First Week

- [ ] Monitor error rates
- [ ] Check user signups
- [ ] Review API key creation
- [ ] Analyze usage patterns
- [ ] Collect user feedback
- [ ] Fix any reported bugs
- [ ] Optimize slow queries
- [ ] Scale resources if needed

---

## Rollback Plan

If deployment fails:

```bash
# Docker
docker compose -f docker-compose.prod.yml down
git checkout previous-working-tag
docker compose -f docker-compose.prod.yml up -d

# Traditional
# Restore database from backup
pg_restore -d requiem_production backup.sql.gz

# Revert code
git revert HEAD
git push

# Restart services
systemctl restart requiems-dashboard
systemctl restart requiems-api
```

---

## Performance Tuning

### Database

```sql
-- Check slow queries
SELECT query, calls, mean_exec_time, stddev_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Check index usage
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;

-- Vacuum analyze
VACUUM ANALYZE;
```

### Rails

```ruby
# config/environments/production.rb

# Enable caching
config.action_controller.perform_caching = true
config.cache_store = :redis_cache_store, { url: ENV['REDIS_URL'] }

# Optimize assets
config.assets.compile = false
config.assets.digest = true

# Set timeouts
Rack::Timeout.timeout = 20  # 20 seconds
```

### Web Server

```nginx
# Nginx optimization
worker_processes auto;
worker_connections 1024;

# Enable gzip
gzip on;
gzip_types text/css application/javascript application/json;

# Cache static assets
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

---

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
systemctl status postgresql

# Test connection
psql -U requiem -h localhost -d requiem_production

# Check max connections
psql -c "SHOW max_connections;"
psql -c "SELECT count(*) FROM pg_stat_activity;"
```

### High Memory Usage

```bash
# Check memory
free -h

# Reduce Rails threads
# .env: RAILS_MAX_THREADS=2

# Optimize PostgreSQL
# postgresql.conf: shared_buffers = 256MB
```

### Slow Responses

```bash
# Profile slow requests
# Add to Gemfile: gem 'rack-mini-profiler'

# Check database query performance
# Rails console:
ActiveRecord::Base.logger = Logger.new(STDOUT)

# Check for N+1 queries
# Add to Gemfile (development): gem 'bullet'
```

---

## Support Resources

- **Documentation**: `/docs` folder
- **GitHub Issues**:
  [github.com/bobadilla-tech/requiems-api/issues](https://github.com/bobadilla-tech/requiems-api/issues)
- **Community**:
  [github.com/bobadilla-tech/requiems-api/discussions](https://github.com/bobadilla-tech/requiems-api/discussions)
- **Self-Hosting Guide**: `docs/SELF_HOSTING.md`
- **User Guide**: `docs/USER_GUIDE.md`

---

**Good luck with your deployment! 🚀**
