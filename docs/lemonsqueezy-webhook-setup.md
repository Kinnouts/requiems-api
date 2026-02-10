# LemonSqueezy Webhook Integration Setup

## Overview

The LemonSqueezy webhook integration is now fully implemented. This document explains how to complete the setup.

## What's Been Implemented

### 1. Database Migration ✅
- Added LemonSqueezy fields to `subscriptions` table:
  - `lemonsqueezy_subscription_id` (string, unique index)
  - `lemonsqueezy_customer_id` (string, indexed)
  - `plan_name` (string, indexed)
  - `cancel_at_period_end` (boolean, default: false)

### 2. Subscription Model ✅
- Validations for `plan_name` and `lemonsqueezy_subscription_id`
- Scopes: `active`, `cancelled`
- Auto-sync to Cloudflare KV when created or plan changes

### 3. Webhook Controller ✅
Location: `apps/dashboard/app/controllers/webhooks/lemonsqueezy_controller.rb`

Handles these webhook events:
- `subscription_created` - New subscription started
- `subscription_updated` - Plan changed or subscription modified
- `subscription_cancelled` / `subscription_expired` - Subscription ended
- `subscription_resumed` - Cancelled subscription reactivated
- `subscription_payment_success` - Successful payment processed

### 4. Cloudflare KV Sync Service ✅
Location: `apps/dashboard/app/services/cloudflare_kv_sync_service.rb`

Automatically syncs subscription changes to Cloudflare KV so the Worker can enforce:
- Plan limits
- Rate limits
- Credit quotas

## Setup Required

### Step 1: Configure Webhook in LemonSqueezy Dashboard

1. Go to your LemonSqueezy dashboard: https://app.lemonsqueezy.com/settings/webhooks
2. Click "Create Webhook"
3. Enter webhook URL: `https://requiems.xyz/webhooks/lemonsqueezy`
4. Select events to send:
   - ✅ `subscription_created`
   - ✅ `subscription_updated`
   - ✅ `subscription_cancelled`
   - ✅ `subscription_expired`
   - ✅ `subscription_resumed`
   - ✅ `subscription_payment_success`
5. Copy the **Signing Secret** (shown after creation)

### Step 2: Update Environment Variable

Update the signing secret in your environment:

```bash
# In infra/docker/.env
LEMONSQUEEZY_SIGNING_SECRET=your_actual_signing_secret_here
```

**Important:** Replace `your_webhook_signing_secret_from_lemonsqueezy` with the actual signing secret from Step 1.

### Step 3: Restart Services

```bash
cd infra/docker
docker compose down
docker compose up -d
```

## Testing the Integration

### Test with LemonSqueezy Test Mode

1. Create a test checkout using test mode variant IDs
2. Complete a test purchase
3. Check Rails logs for webhook events:

```bash
docker compose logs -f dashboard | grep LemonSqueezy
```

You should see logs like:
```
[LemonSqueezy Webhook] Received: subscription_created
[LemonSqueezy] Subscription created for user 123: developer
[CloudflareKV] Synced API key abc123 to plan: developer
```

### Verify Cloudflare KV Sync

After a successful webhook, check that the API key was synced to Cloudflare KV:

```bash
# Using Wrangler CLI
cd apps/edge-auth
npx wrangler kv:key get "key:YOUR_API_KEY" --namespace-id=7cc847da3f3143b2ba8f7c531f416b35
```

Should return:
```json
{
  "userId": "123",
  "plan": "developer",
  "billingCycleStart": "2026-02-10T00:00:00Z",
  "createdAt": "2026-01-15T10:30:00Z"
}
```

## Webhook Flow

```
User Checkout → LemonSqueezy → Webhook Event
                                     ↓
                          POST /webhooks/lemonsqueezy
                                     ↓
                          Verify HMAC-SHA256 Signature
                                     ↓
                          Update Subscription Record
                                     ↓
                          Sync to Cloudflare KV
                                     ↓
                          Worker Enforces New Plan Limits
```

## Subscription Status Mapping

| LemonSqueezy Status | Plan Name    | Cloudflare KV Plan |
|---------------------|--------------|---------------------|
| active              | developer    | developer           |
| trialing            | developer    | developer           |
| cancelled           | free         | free                |
| expired             | free         | free                |
| on_trial            | developer    | developer           |

## Variant ID to Plan Mapping

Configured in `.env`:

| Variant ID | Plan           | Billing |
|------------|----------------|---------|
| 822585     | Developer      | Monthly |
| 822595     | Developer      | Yearly  |
| 822596     | Business       | Monthly |
| 822601     | Business       | Yearly  |
| 822603     | Professional   | Monthly |
| 822604     | Professional   | Yearly  |

## Security

- ✅ Webhook signature verification using HMAC-SHA256
- ✅ CSRF token skipped for webhook endpoint
- ✅ Cloudflare API token secured in environment
- ✅ Backend secret validation between Worker and Rails

## Troubleshooting

### Webhook Returns 401 Unauthorized

**Cause:** Invalid signature or missing signing secret

**Fix:**
1. Verify `LEMONSQUEEZY_SIGNING_SECRET` is set correctly
2. Check webhook logs: `docker compose logs dashboard | grep "Invalid signature"`
3. Regenerate signing secret in LemonSqueezy dashboard

### Plan Not Syncing to Cloudflare KV

**Cause:** Cloudflare API credentials missing or incorrect

**Fix:**
1. Verify these environment variables are set:
   - `CLOUDFLARE_ACCOUNT_ID`
   - `CLOUDFLARE_KV_NAMESPACE_ID`
   - `CLOUDFLARE_API_TOKEN`
2. Test API token permissions: needs "Workers KV Storage:Edit"

### Subscription Created but User Not Found

**Cause:** Custom data `user_id` not passed in checkout URL

**Fix:** Verify checkout URL includes:
```
?checkout[custom][user_id]=123
```

This is automatically added by `billing_controller.rb:create_lemonsqueezy_checkout_url`

## Next Steps

1. ✅ Set webhook URL in LemonSqueezy dashboard
2. ✅ Copy signing secret to `.env`
3. ✅ Restart services
4. ✅ Test with test purchase
5. ✅ Monitor webhook logs
6. ✅ Verify KV sync works

## API Key Sync

When a subscription changes, **all active API keys** for that user are automatically updated in Cloudflare KV. This ensures:

- Existing API keys immediately get new rate limits
- Credit quotas are updated for the new billing cycle
- No need for users to regenerate API keys

## Support

If webhooks are failing:

1. Check Rails logs: `docker compose logs -f dashboard`
2. Check LemonSqueezy webhook attempts: Dashboard → Settings → Webhooks
3. Verify signature verification is working
4. Test Cloudflare KV connectivity manually
