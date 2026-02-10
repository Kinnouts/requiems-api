# Self-Hosting Requiems API

This guide will help you deploy and run your own instance of Requiems API.

## Table of Contents

1. [Why Self-Host?](#why-self-host)
2. [Architecture Overview](#architecture-overview)
3. [Prerequisites](#prerequisites)
4. [Quick Start with Docker](#quick-start-with-docker)
5. [Manual Installation](#manual-installation)
6. [Configuration](#configuration)
7. [Production Deployment](#production-deployment)
8. [Monitoring & Maintenance](#monitoring--maintenance)
9. [Troubleshooting](#troubleshooting)

---

## Why Self-Host?

### Benefits

- **Full Control**: Own your data and infrastructure
- **Privacy**: Keep sensitive data in-house
- **Customization**: Modify the platform to fit your needs
- **Cost**: No subscription fees for high-volume usage
- **Compliance**: Meet specific regulatory requirements

### When to Self-Host

✅ Self-hosting is great for:

- Enterprise compliance requirements
- Very high API volume (>10M requests/month)
- Custom feature requirements
- On-premise infrastructure mandates
- Development and testing

❌ Use our managed service instead if:

- You want zero maintenance
- You need guaranteed uptime (99.9% SLA)
- You want automatic updates
- You prefer pay-as-you-go pricing
- You don't have DevOps resources

---

## Architecture Overview

Requiems API consists of three main components:

```
┌─────────────┐      ┌──────────────────┐      ┌──────────────┐
│  Dashboard  │◄─────┤  Edge Gateway    │◄─────┤   Go API     │
│ (Rails 8)   │      │  (Cloudflare)    │      │  (Backend)   │
└─────────────┘      └──────────────────┘      └──────────────┘
       │                      │                          │
       ▼                      ▼                          ▼
  PostgreSQL            KV + D1 SQLite            PostgreSQL
```

### Components

1. **Dashboard (Rails 8)**
   - User management and authentication
   - API key generation
   - Usage analytics and billing
   - Admin panel

2. **Edge Gateway (Cloudflare Worker)**
   - API authentication
   - Rate limiting
   - Credit tracking
   - Request proxying

3. **Go API Backend**
   - Business logic
   - API endpoints
   - Data processing

4. **Databases**
   - **PostgreSQL**: User data, subscriptions, usage logs
   - **Cloudflare KV**: API key storage (fast lookups)
   - **Cloudflare D1**: Usage tracking (SQLite at the edge)

---

## Prerequisites

### Required

- **Docker** 20.10+ & Docker Compose 2.0+
- **PostgreSQL** 14+ (or use Docker)
- **Redis** 7+ (or use Docker)
- **Node.js** 18+ (for Cloudflare Worker)
- **Go** 1.23+ (for API backend)
- **Ruby** 3.3+ (for Dashboard)

### Optional

- **Cloudflare Account** (for edge deployment)
- **Domain name** (for custom URLs)
- **SSL certificate** (Let's Encrypt recommended)

---

## Quick Start with Docker

### 1. Clone the Repository

```bash
git clone https://github.com/bobadilla-tech/requiems-api.git
cd requiems-api
```

### 2. Start All Services

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

This starts:

- Go API (http://localhost:8080)
- Rails Dashboard (http://localhost:3000)
- PostgreSQL (localhost:5432)
- Redis (localhost:6379)

### 3. Set Up the Database

```bash
# In another terminal
docker compose -f docker-compose.dev.yml exec dashboard rails db:create db:migrate db:seed
```

### 4. Create Your First Admin User

```bash
docker compose -f docker-compose.dev.yml exec dashboard rails console

# In the console:
User.create!(
  email: 'admin@example.com',
  password: 'password123',
  password_confirmation: 'password123',
  admin: true
)
```

### 5. Access the Dashboard

Visit http://localhost:3000 and log in with your admin credentials.

---

## Manual Installation

### 1. Install Dependencies

#### macOS

```bash
brew install go ruby node postgresql redis
```

#### Ubuntu/Debian

```bash
sudo apt update
sudo apt install golang-go ruby-full nodejs npm postgresql redis-server
```

### 2. Set Up PostgreSQL

```bash
# Create database user
createuser -s requiem

# Create database
createdb -O requiem requiem

# Set password
psql requiem -c "ALTER USER requiem WITH PASSWORD 'your_password';"
```

### 3. Set Up Environment Variables

Create `.env` files:

#### Go API (`apps/api/.env`)

```env
DATABASE_URL=postgresql://requiem:your_password@localhost:5432/requiem
PORT=8080
ENV=development
```

#### Rails Dashboard (`apps/dashboard/.env`)

```env
DATABASE_URL=postgresql://requiem:your_password@localhost:5432/requiem
REDIS_URL=redis://localhost:6379/0
SECRET_KEY_BASE=generate_with_rails_secret
LEMONSQUEEZY_API_KEY=your_key_here
LEMONSQUEEZY_STORE_SLUG=your_store
BACKEND_API_URL=http://localhost:8080
```

Generate Rails secret:

```bash
cd apps/dashboard
bundle exec rails secret
```

### 4. Install & Run Go API

```bash
cd apps/api
go mod download
go run main.go
```

### 5. Install & Run Rails Dashboard

```bash
cd apps/dashboard
bundle install
rails db:create db:migrate db:seed
rails server
```

### 6. Deploy Cloudflare Worker (Optional)

```bash
cd apps/edge-auth
npm install
npx wrangler login
npx wrangler deploy
```

---

## Configuration

### Environment Variables

#### Core Settings

| Variable          | Description                  | Default | Required |
| ----------------- | ---------------------------- | ------- | -------- |
| `DATABASE_URL`    | PostgreSQL connection string | -       | Yes      |
| `REDIS_URL`       | Redis connection string      | -       | Yes      |
| `SECRET_KEY_BASE` | Rails secret key             | -       | Yes      |
| `PORT`            | Server port                  | 3000    | No       |

#### External Services

| Variable                  | Description         | Required       |
| ------------------------- | ------------------- | -------------- |
| `LEMONSQUEEZY_API_KEY`    | Payment processing  | For billing    |
| `LEMONSQUEEZY_STORE_SLUG` | Store identifier    | For billing    |
| `CLOUDFLARE_API_TOKEN`    | For edge deployment | For Cloudflare |

#### Email (SMTP)

| Variable        | Description    | Default        |
| --------------- | -------------- | -------------- |
| `SMTP_ADDRESS`  | SMTP server    | smtp.gmail.com |
| `SMTP_PORT`     | SMTP port      | 587            |
| `SMTP_USERNAME` | Email username | -              |
| `SMTP_PASSWORD` | Email password | -              |
| `SMTP_DOMAIN`   | Email domain   | -              |

#### Feature Flags

| Variable                | Description         | Default |
| ----------------------- | ------------------- | ------- |
| `ENABLE_REGISTRATION`   | Allow new signups   | true    |
| `ENABLE_API_PLAYGROUND` | Interactive testing | true    |
| `RATE_LIMIT_ENABLED`    | Rate limiting       | true    |

### Database Configuration

#### PostgreSQL Performance Tuning

Edit `postgresql.conf`:

```conf
# Memory
shared_buffers = 256MB
effective_cache_size = 1GB

# Connections
max_connections = 100

# Query Planning
random_page_cost = 1.1  # For SSD

# Write Ahead Log
wal_buffers = 16MB
checkpoint_completion_target = 0.9
```

#### Database Backups

Set up automated backups:

```bash
# Daily backup script
#!/bin/bash
pg_dump requiem | gzip > /backups/requiem_$(date +%Y%m%d).sql.gz

# Keep last 30 days
find /backups -name "requiem_*.sql.gz" -mtime +30 -delete
```

Add to crontab:

```bash
0 2 * * * /path/to/backup_script.sh
```

---

## Production Deployment

### Option 1: Docker Production

```bash
# Build production images
docker compose -f docker-compose.prod.yml build

# Start services
docker compose -f docker-compose.prod.yml up -d

# Run migrations
docker compose -f docker-compose.prod.yml exec dashboard rails db:migrate

# Check status
docker compose -f docker-compose.prod.yml ps
```

### Option 2: Kubernetes

Coming soon! See `infra/k8s/` for manifests.

### Option 3: Traditional VPS

#### System Requirements

**Minimum**:

- 2 CPU cores
- 4GB RAM
- 40GB SSD
- Ubuntu 22.04 LTS

**Recommended**:

- 4 CPU cores
- 8GB RAM
- 100GB SSD
- Load balancer

#### Install System Dependencies

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y \
  postgresql-14 \
  redis-server \
  nginx \
  certbot \
  python3-certbot-nginx \
  git \
  build-essential

# Install Ruby via rbenv
git clone https://github.com/rbenv/rbenv.git ~/.rbenv
git clone https://github.com/rbenv/ruby-build.git ~/.rbenv/plugins/ruby-build
echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc
echo 'eval "$(rbenv init -)"' >> ~/.bashrc
source ~/.bashrc
rbenv install 3.3.0
rbenv global 3.3.0

# Install Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Configure Nginx

Create `/etc/nginx/sites-available/requiems`:

```nginx
upstream dashboard {
  server 127.0.0.1:3000;
}

upstream api {
  server 127.0.0.1:8080;
}

server {
  listen 80;
  server_name yourdomain.com;
  return 301 https://$server_name$request_uri;
}

server {
  listen 443 ssl http2;
  server_name yourdomain.com;

  ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

  # Dashboard
  location / {
    proxy_pass http://dashboard;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }

  # API Backend
  location /api/ {
    proxy_pass http://api/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

Enable site:

```bash
sudo ln -s /etc/nginx/sites-available/requiems /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

#### SSL with Let's Encrypt

```bash
sudo certbot --nginx -d yourdomain.com
```

#### Process Management with Systemd

Create `/etc/systemd/system/requiems-dashboard.service`:

```ini
[Unit]
Description=Requiems Dashboard
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=deploy
WorkingDirectory=/var/www/requiems/apps/dashboard
Environment="RAILS_ENV=production"
Environment="PORT=3000"
ExecStart=/home/deploy/.rbenv/shims/bundle exec rails server
Restart=always

[Install]
WantedBy=multi-user.target
```

Create `/etc/systemd/system/requiems-api.service`:

```ini
[Unit]
Description=Requiems API
After=network.target postgresql.service

[Service]
Type=simple
User=deploy
WorkingDirectory=/var/www/requiems/apps/api
Environment="ENV=production"
ExecStart=/var/www/requiems/apps/api/bin/api
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start services:

```bash
sudo systemctl daemon-reload
sudo systemctl enable requiems-dashboard requiems-api
sudo systemctl start requiems-dashboard requiems-api
```

---

## Monitoring & Maintenance

### Health Checks

Check service health:

```bash
# Dashboard
curl http://localhost:3000/healthz

# API
curl http://localhost:8080/healthz
```

### Logs

#### Docker

```bash
docker compose logs -f dashboard
docker compose logs -f api
```

#### Systemd

```bash
sudo journalctl -u requiems-dashboard -f
sudo journalctl -u requiems-api -f
```

### Monitoring Tools

#### Prometheus + Grafana

Coming soon! See `infra/monitoring/` for configs.

#### Application Performance

Use tools like:

- **New Relic** - APM
- **DataDog** - Infrastructure monitoring
- **Sentry** - Error tracking

### Database Maintenance

#### Vacuum & Analyze

```bash
# Run weekly
docker compose exec postgres psql -U requiem -c "VACUUM ANALYZE;"
```

#### Index Maintenance

```sql
-- Find missing indexes
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE schemaname = 'public'
ORDER BY abs(correlation) DESC;

-- Reindex if needed
REINDEX DATABASE requiem;
```

### Updates

#### Pull Latest Code

```bash
git pull origin main
```

#### Update Dependencies

```bash
# Go API
cd apps/api
go mod tidy
go build

# Rails Dashboard
cd apps/dashboard
bundle update
rails db:migrate RAILS_ENV=production
```

#### Restart Services

```bash
docker compose restart
# OR
sudo systemctl restart requiems-dashboard requiems-api
```

---

## Troubleshooting

### Common Issues

#### Database Connection Failed

```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U requiem -h localhost -d requiem

# Check credentials in .env file
```

#### Port Already in Use

```bash
# Find process using port
sudo lsof -i :3000

# Kill process
kill -9 <PID>
```

#### Migrations Failing

```bash
# Reset database (CAUTION: Deletes data!)
rails db:drop db:create db:migrate

# Or fix specific migration
rails db:migrate:status
rails db:migrate:up VERSION=20240101000000
```

#### High Memory Usage

```bash
# Check memory
free -h

# Optimize Rails
# Set RAILS_MAX_THREADS=2 in .env

# Optimize PostgreSQL
# Reduce shared_buffers in postgresql.conf
```

#### Slow API Responses

```bash
# Check database performance
SELECT * FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;

# Add missing indexes
rails db:migrate

# Enable query caching
# Set CACHE_STORE=redis in .env
```

### Getting Help

- **Documentation**: Check `docs/` folder
- **GitHub Issues**:
  [github.com/bobadilla-tech/requiems-api/issues](https://github.com/bobadilla-tech/requiems-api/issues)
- **Discussions**:
  [github.com/bobadilla-tech/requiems-api/discussions](https://github.com/bobadilla-tech/requiems-api/discussions)
- **Community**: Join our Discord (coming soon)

---

## Security Considerations

### API Keys

- Store in environment variables, not code
- Use different keys for dev/staging/prod
- Rotate regularly
- Never commit to version control

### Database

- Use strong passwords (20+ characters)
- Restrict network access (firewall rules)
- Enable SSL connections
- Regular backups (test restoration!)
- Encrypt backups

### Server

- Keep system updated
- Use firewall (ufw/iptables)
- Disable SSH password auth (use keys)
- Enable fail2ban
- Regular security audits

### Application

- Set `RAILS_ENV=production`
- Use HTTPS only
- Configure CORS properly
- Enable CSRF protection
- Set secure session cookies

---

## License

Requiems API is open source under the MIT License. See [LICENSE](../../LICENSE)
file.

---

## Support

- 🌟 Star us on [GitHub](https://github.com/bobadilla-tech/requiems-api)
- 🐛 Report bugs in
  [Issues](https://github.com/bobadilla-tech/requiems-api/issues)
- 💬 Ask questions in
  [Discussions](https://github.com/bobadilla-tech/requiems-api/discussions)
- 📧 Enterprise support: enterprise@requiems.xyz

Built with ❤️ by Bobadilla Tech
