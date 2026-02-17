# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

We take the security of Requiems API seriously. If you discover a security
vulnerability, please report it responsibly.

### Where to Report

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report security vulnerabilities to:

**Email:** eliaz@bobadilla.tech

### What to Include

Please include the following information in your report:

- **Description** of the vulnerability
- **Steps to reproduce** the issue
- **Potential impact** of the vulnerability
- **Suggested fix** (if you have one)
- **Your contact information** for follow-up

### What to Expect

- **Acknowledgment:** We'll acknowledge receipt of your report within 48 hours
- **Updates:** We'll keep you informed about our progress as we investigate
- **Timeline:** We aim to validate and patch critical vulnerabilities within 7
  days
- **Credit:** We'll credit you in our security advisories (unless you prefer to
  remain anonymous)

## Security Best Practices

### For Developers and Contributors

1. **Secure Your API Keys**
   - Store API keys in environment variables, never in code
   - Use different keys for different environments (test vs live)
   - Rotate keys regularly
   - Revoke compromised keys immediately

2. **Code Security**
   - Never commit `.env` files or secrets to version control
   - Review code for common vulnerabilities (SQL injection, XSS, etc.)
   - Keep dependencies updated
   - Run security scans before submitting PRs

3. **Access Control**
   - Use strong passwords and enable 2FA on your GitHub account
   - Follow the principle of least privilege
   - Review and approve third-party access carefully

## Scope

This security policy covers:

- ✅ The Requiems API source code in this repository
- ✅ Official deployment configurations
- ✅ Dependencies we directly control
- ❌ Third-party services or integrations
- ❌ Infrastructure outside of this repository

## Questions?

For general security questions or concerns, please contact:
eliaz@bobadilla.tech

Thank you for helping keep Requiems API and our users secure! 🛡️
