# API Management Deployment Guide

## First-Time Setup

### 1. Set Up Cloudflare Secrets

First, create the required secrets in Cloudflare:

```bash
cd apps/api-management

# Required: API Management authentication key (64+ character random string)
wrangler secret put API_MANAGEMENT_API_KEY --env production

# Optional: Swagger UI basic auth (production only)
wrangler secret put SWAGGER_USERNAME --env production
wrangler secret put SWAGGER_PASSWORD --env production
```

### 2. Initial Manual Deployment

Deploy manually the first time to verify everything works:

```bash
cd apps/api-management

# Install dependencies
bun install

# Run tests locally first
bun run typecheck
bunx vitest run

# Deploy to production
bunx wrangler deploy --env production
```

This will deploy to: **https://api-management.requiems.xyz**

### 3. Set Up GitHub Secrets for Auto-Deploy

To enable automatic deployments when you push to `main`, add your Cloudflare API token to GitHub:

1. **Get your Cloudflare API Token:**
   - Go to https://dash.cloudflare.com/profile/api-tokens
   - Click "Create Token"
   - Use "Edit Cloudflare Workers" template
   - Or create custom token with permissions:
     - Account > Workers Scripts > Edit
     - Zone > Workers Routes > Edit

2. **Add to GitHub Secrets:**
   - Go to your GitHub repo → Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `CLOUDFLARE_API_TOKEN`
   - Value: Your Cloudflare API token
   - Click "Add secret"

## How Auto-Deploy Works

Once set up, the workflow automatically:

1. **On every push to `main`** that changes `apps/api-management/**`:
   - Runs TypeScript type checking
   - Runs tests
   - **If tests pass**: Deploys to Cloudflare Workers production

2. **On pull requests**:
   - Runs TypeScript type checking
   - Runs tests
   - Does NOT deploy (deploy only happens on main branch)

## Manual Deployment (When Needed)

You can always deploy manually from your local machine:

```bash
cd apps/api-management

# Deploy to production
bunx wrangler deploy --env production

# Or just run "deploy" (defaults to production)
bunx wrangler deploy
```

## Verify Deployment

After deployment (manual or automatic), verify it's working:

```bash
# Health check (no auth required)
curl https://api-management.requiems.xyz/healthz

# Test with API key (replace with your actual key)
curl https://api-management.requiems.xyz/healthz \
  -H "X-API-Management-Key: your-api-key-here"
```

Expected response:
```json
{
  "status": "ok",
  "service": "api-management"
}
```

## Checking Logs

View production logs:

```bash
cd apps/api-management
bunx wrangler tail requiem-api-management-production
```

## Troubleshooting

### Deployment fails with "Authentication error"
- Check that `CLOUDFLARE_API_TOKEN` is set in GitHub Secrets
- Verify the token has Workers Scripts Edit permissions

### "KV namespace not found" error
- Check `wrangler.toml` has correct KV namespace ID
- Verify KV namespace exists in Cloudflare dashboard

### "D1 database not found" error
- Check `wrangler.toml` has correct D1 database ID
- Verify D1 database exists in Cloudflare dashboard

### Secrets not working
- Secrets are environment-specific
- Use `--env production` when setting secrets
- Secrets don't show in dashboard, verify by checking logs

## Environment Variables

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `API_MANAGEMENT_API_KEY` | Secret | Yes | Authentication key for Rails dashboard (min 64 chars) |
| `SWAGGER_USERNAME` | Secret | No | Basic auth username for `/docs` in production |
| `SWAGGER_PASSWORD` | Secret | No | Basic auth password for `/docs` in production |
| `ENVIRONMENT` | Var | No | Set to "production" in prod environment |

**Secrets vs Variables:**
- **Secrets**: Sensitive data (API keys, passwords) - set with `wrangler secret put`
- **Variables**: Non-sensitive config - set in `wrangler.toml` under `[env.production.vars]`

## Rollback

If you need to rollback to a previous version:

```bash
# List recent deployments
bunx wrangler deployments list --env production

# Rollback to a specific deployment
bunx wrangler rollback [DEPLOYMENT_ID] --env production
```

## Monitoring

- **Cloudflare Dashboard**: https://dash.cloudflare.com
  - View analytics, errors, and performance metrics
  - Monitor request volume and error rates

- **GitHub Actions**: Check deployment status
  - Go to Actions tab in your repo
  - View deployment logs and test results
