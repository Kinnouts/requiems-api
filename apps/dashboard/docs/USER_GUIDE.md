# Requiems API User Guide

Welcome to Requiems API! This guide will help you get started with our platform and make the most of our APIs.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Dashboard Overview](#dashboard-overview)
3. [Managing API Keys](#managing-api-keys)
4. [Making API Requests](#making-api-requests)
5. [Monitoring Usage](#monitoring-usage)
6. [Billing & Subscriptions](#billing--subscriptions)
7. [Account Settings](#account-settings)
8. [Troubleshooting](#troubleshooting)

---

## Getting Started

###1. Create an Account

1. Visit [requiems.xyz](https://requiems.xyz)
2. Click "Get Started" or "Sign Up"
3. Enter your email and create a password
4. Verify your email address
5. You're ready to go! Start on the Free plan (500 requests/month)

### 2. Browse Available APIs

- Go to **APIs** in the navigation menu
- Explore categories: Email APIs, Text APIs, and more
- Click on any API to see documentation and examples

### 3. Create Your First API Key

1. Go to **Dashboard** → **API Keys**
2. Click "Create New Key"
3. Give it a name (e.g., "My First Key")
4. Choose environment: **Test** or **Live**
5. Click "Generate Key"
6. **IMPORTANT**: Copy and save your key immediately - it's only shown once!

---

## Dashboard Overview

Your dashboard shows:

- **Usage Progress**: Current month's API calls vs. your plan limit
- **Credits Remaining**: How many credits you have left
- **Recent Activity**: Your last 10 API requests
- **API Keys Widget**: Quick view of your active keys

### Quick Stats

- **Total Requests**: All-time API calls
- **This Month**: Requests made in the current billing period
- **Average Response Time**: Performance metric
- **Credits Remaining**: Based on your current plan

---

## Managing API Keys

### Creating API Keys

**All API keys share your monthly quota** - it doesn't matter which key you use, they all count toward the same limit.

#### Test vs. Live Keys

- **Test Keys** (`rq_test_...`): Use for development and testing
- **Live Keys** (`rq_live_...`): Use in production applications

Both types count toward your monthly quota.

### Key Actions

#### Regenerate a Key
1. Go to **API Keys**
2. Click "Regenerate" on the key you want to change
3. The old key is immediately revoked
4. Copy the new key (shown only once)
5. Update your applications with the new key

#### Revoke a Key
1. Go to **API Keys**
2. Click "Revoke" on the key
3. Confirm the action
4. The key is immediately disabled

**⚠️ Warning**: Revoking a key will break any applications using it.

### Best Practices

- ✅ Use descriptive names ("Production API", "Mobile App", "Testing")
- ✅ Create separate keys for different applications
- ✅ Rotate keys periodically for security
- ✅ Never share your API keys publicly
- ✅ Store keys in environment variables, not in code
- ❌ Don't commit keys to version control

---

## Making API Requests

### Authentication

Include your API key in the `X-API-Key` header:

```bash
curl -X POST https://api.requiems.xyz/v1/endpoint \
  -H "X-API-Key: rq_live_your_api_key_here" \
  -H "Content-Type: application/json" \
  -d '{"param": "value"}'
```

### Example: Check Disposable Email

```bash
curl -X POST https://api.requiems.xyz/v1/email/disposable/check \
  -H "X-API-Key: rq_live_xxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'
```

**Response**:
```json
{
  "email": "test@example.com",
  "is_disposable": false,
  "domain": "example.com"
}
```

### Rate Limits

Rate limits depend on your plan:

| Plan | Requests/Minute | Monthly Requests |
|------|----------------|------------------|
| Free | 30/min | 500 |
| Developer | 5,000/min | 100,000 |
| Business | 10,000/min | 1,000,000 |
| Professional | 50,000/min | 10,000,000 |

### Response Headers

Every response includes usage information:

```
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4995
X-RateLimit-Reset: 1640000000
X-Credits-Used: 1
X-Credits-Remaining: 99500
```

### Error Handling

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded. Please try again in 30 seconds.",
    "details": {
      "retry_after": 30
    }
  }
}
```

Common status codes:

- `200` - Success
- `400` - Bad Request (invalid parameters)
- `401` - Unauthorized (invalid or missing API key)
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error

---

## Monitoring Usage

### Usage & Analytics Dashboard

Access detailed usage statistics at **Dashboard** → **Usage & Analytics**.

#### Date Range Filters

- Last 7 Days
- Last 30 Days
- Last 90 Days
- Custom Range (with date picker)

#### Charts

1. **Requests Over Time**: Line chart showing daily request volume
2. **Requests by Endpoint**: Bar chart of most-used APIs
3. **Response Code Distribution**: Pie chart of success/error rates
4. **Response Times**: Average latency by endpoint

#### Recent Requests Table

View detailed logs of your last requests:

- Timestamp
- Endpoint called
- HTTP method
- Status code
- Response time
- Credits used
- API key used

#### Export Data

Click "Export CSV" to download your usage data for analysis.

### Understanding Credits

Each API call costs credits based on complexity:

- **Simple APIs** (Advice, Lorem Ipsum): 1 credit
- **Text APIs** (Quotes, Words): 1-2 credits
- **Email APIs** (Disposable Check): 1 credit
- **Batch Operations**: Credits per item

Your plan includes a monthly credit allowance. Track usage to avoid hitting your limit.

---

## Billing & Subscriptions

### Plans & Pricing

#### Free Plan
- 500 requests/month
- 30 requests/minute
- Community support
- Perfect for testing

#### Developer Plan - $30/month ($25/year)
- 100,000 requests/month
- 5,000 requests/minute
- Email support
- Ideal for small projects

#### Business Plan - $75/month ($62.50/year) ⭐ Most Popular
- 1,000,000 requests/month
- 10,000 requests/minute
- Priority support
- EU & US data centers
- For growing businesses

#### Professional Plan - $150/month ($125/year)
- 10,000,000 requests/month
- 50,000 requests/minute
- Dedicated support
- EU & US data centers
- SLA guarantee

#### Enterprise Plan - Custom Pricing
- Unlimited requests
- Custom rate limits
- Private servers
- White-label options
- Dedicated account manager

### Upgrading Your Plan

1. Go to **Dashboard** → **Billing**
2. Click "Upgrade Plan"
3. Select your desired plan
4. Choose Monthly or Yearly billing
5. Complete checkout via LemonSqueezy
6. Your new limits apply immediately!

### Managing Subscriptions

- **Change Plan**: Upgrade or downgrade anytime
- **Update Payment Method**: Click "Manage Subscription" → LemonSqueezy portal
- **View Invoices**: Access all past invoices in the Billing page
- **Cancel Subscription**: Cancel anytime - access continues until period end

### Billing Cycle

- Subscriptions renew automatically each month/year
- Usage resets at the start of each billing cycle
- Prorated charges apply when upgrading mid-cycle
- Cancel anytime without penalties

---

## Account Settings

### Profile Information

Update your account details at **Dashboard** → **Settings**:

- Name
- Email address
- Company name
- Password

### Email Preferences

Control which emails you receive:

- ✉️ **Usage Alerts**: Notified at 80% and 100% of quota
- ✉️ **Weekly Reports**: Summary of your usage and stats
- ✉️ **Product Updates**: New features and improvements

### Security

- Change your password regularly
- Enable two-factor authentication (coming soon)
- Review active API keys periodically
- Revoke unused keys

### Danger Zone

#### Delete Account

**⚠️ This action is permanent and cannot be undone.**

To delete your account:

1. Go to **Settings** → **Danger Zone**
2. Click "Delete Account"
3. Type your email to confirm
4. Click "Permanently Delete Account"

All data will be deleted:
- API keys (immediately revoked)
- Usage history
- Billing information
- Account settings

Your subscription will be canceled automatically.

---

## Troubleshooting

### Common Issues

#### "Invalid API Key" Error

**Solutions**:
- Verify you're using the correct key
- Check that the key hasn't been revoked
- Ensure the key is in the correct header: `X-API-Key`
- Make sure there are no extra spaces or characters

#### "Rate Limit Exceeded" Error

**Solutions**:
- Wait for the rate limit window to reset (check `X-RateLimit-Reset` header)
- Upgrade to a higher plan for increased limits
- Implement exponential backoff in your application
- Batch requests where possible

#### "Insufficient Credits" Error

**Solutions**:
- Check your usage at **Dashboard** → **Usage**
- Upgrade your plan for more monthly credits
- Wait for your billing cycle to reset
- Contact support for manual credit adjustment

#### API Keys Not Showing in Dashboard

**Solutions**:
- Refresh the page
- Clear browser cache
- Try a different browser
- Check if you're logged into the correct account

#### Billing Issues

**Solutions**:
- Verify payment method in LemonSqueezy portal
- Check for declined payments
- Contact LemonSqueezy support for payment issues
- Reach out to our support team

### Getting Help

#### Documentation
- **API Reference**: [docs.requiems.xyz](https://docs.requiems.xyz)
- **Examples**: Check out code examples for common use cases
- **Blog**: Tips, tutorials, and best practices

#### Support

- **Free Plan**: Community support (GitHub Discussions)
- **Paid Plans**: Email support (support@requiems.xyz)
- **Business/Professional**: Priority support (24-48h response)
- **Enterprise**: Dedicated support channel

#### Status Page

Check system status at [status.requiems.xyz](https://status.requiems.xyz):
- API uptime
- Scheduled maintenance
- Incident reports
- Performance metrics

---

## Best Practices

### Security

1. **Never expose API keys**: Store in environment variables
2. **Use HTTPS**: Always make requests over HTTPS
3. **Rotate keys**: Change keys every 90 days
4. **Separate keys**: Use different keys for dev/staging/production
5. **Monitor usage**: Watch for unusual activity

### Performance

1. **Cache responses**: Cache when data doesn't change frequently
2. **Batch requests**: Use batch endpoints when available
3. **Handle rate limits**: Implement exponential backoff
4. **Use webhooks**: For async operations (coming soon)
5. **Monitor latency**: Track response times in your analytics

### Cost Optimization

1. **Start with Free**: Test thoroughly before upgrading
2. **Annual billing**: Save ~33% with yearly plans
3. **Monitor usage**: Set up usage alerts
4. **Optimize calls**: Reduce unnecessary API requests
5. **Right-size plan**: Choose based on actual usage, not estimates

---

## FAQ

**Q: Do API keys share the same monthly quota?**
A: Yes! All your API keys (test and live) count toward the same monthly request limit.

**Q: What happens when I hit my limit?**
A: You'll receive a `429 Too Many Requests` error. Requests will work again when your billing cycle resets or you upgrade.

**Q: Can I test APIs without an account?**
A: Yes! Use the interactive playground on each API documentation page for testing.

**Q: Is there a free tier?**
A: Yes! The Free plan includes 500 requests/month with no credit card required.

**Q: Can I cancel anytime?**
A: Absolutely. Cancel anytime and retain access until your billing period ends.

**Q: Do you offer refunds?**
A: We offer a 30-day money-back guarantee for new subscriptions.

**Q: Where is data processed?**
A: We have servers in both US and EU regions. Business+ plans can choose their preferred region.

**Q: Is Requiems API open source?**
A: Yes! All our code is open source at [github.com/bobadilla-tech/requiems-api](https://github.com/bobadilla-tech/requiems-api). You can self-host if desired.

**Q: What's your SLA?**
A: Professional and Enterprise plans include a 99.9% uptime SLA.

---

## Next Steps

- 📚 [Browse API Documentation](https://docs.requiems.xyz)
- 💡 [View Code Examples](https://requiems.xyz/examples)
- 🚀 [Upgrade Your Plan](https://requiems.xyz/pricing)
- 🐛 [Report Issues](https://github.com/bobadilla-tech/requiems-api/issues)
- 💬 [Join Community](https://github.com/bobadilla-tech/requiems-api/discussions)

---

**Need help?** Contact us at support@requiems.xyz or visit our [Help Center](https://help.requiems.xyz).

Built with ❤️ by Bobadilla Tech | [GitHub](https://github.com/bobadilla-tech/requiems-api) | [Twitter](https://twitter.com/requiemsapi) | [Blog](https://blog.requiems.xyz)
