# Changelog

All notable changes to Requiems API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-02-17

### Added

- **Continuous Integration Pipeline**: Established comprehensive CI workflow
  with automated testing across all three applications
  - Go API: Tests with race detection, golangci-lint (21 linters)
  - Rails Dashboard: Full test suite (60 tests), security scans (Brakeman,
    bundler-audit), RuboCop linting
  - Cloudflare Worker: TypeScript checks, Vitest suite (71 tests), Biome linting
  - Path-based execution for efficient builds
  - Security scans as blocking requirements
- Rails Dashboard test suite fully passing with 207 assertions
- Go API test coverage with race detection
- Cloudflare Worker test suite with 29% coverage

### Fixed

- Rails Dashboard: ApiKey model callback execution order (before_validation)
- Rails Dashboard: Devise mapping initialization in test environment
- Rails Dashboard: User model status field handling for suspended/banned states
- Rails Dashboard: Pagy pagination integration
- Rails Dashboard: Admin controller status field updates
- Go API: HTTP server timeout configurations for security
- Go API: Code formatting and linting issues (gofmt, goimports, gosec)
- Cloudflare Worker: Biome formatting across all source files

## [0.0.2] - 2026-02-08

### Added

- **Backend API Deployed**: Go-based internal API with PostgreSQL
  - Email validation endpoints (disposable email detection)
  - Text utility endpoints (advice, lorem ipsum, quotes, words)
  - Domain-driven design architecture
  - Health check endpoints
  - Database migrations system
- **Dashboard Application**: Rails 8 web application deployed
  - User authentication and registration (Devise)
  - API key management system
  - Usage tracking and analytics
  - Admin panel for user management
  - Subscription and billing integration
  - Tailwind CSS + Turbo/Stimulus frontend
- PostgreSQL database with dual migration systems
- Docker Compose development environment

### Infrastructure

- Production deployment configuration
- Database backup and monitoring systems
- Application health monitoring

## [0.0.1] - 2025-12-15

### Added

- **Project Initialization**: Multi-language monorepo architecture established
- **Design and Planning**: Initial architecture design for API platform
  - Three-tier architecture: Edge Gateway → Backend API → Database
  - Technology stack selection (Go, Rails, Cloudflare Workers)
  - Database schema design
  - API endpoint planning
- **Cloudflare Worker Gateway**: Authentication and rate limiting layer
  - API key validation using Cloudflare KV
  - Rate limiting with KV counters
  - Credit usage tracking with D1 SQLite
  - Request proxying to backend
  - Usage analytics and logging
- **Development Environment**: Local development setup
  - Docker Compose configuration
  - Hot reload for all applications
  - Database seeding scripts
  - Environment variable management

### Infrastructure

- GitHub repository structure
- Monorepo organization (apps/, infra/, docs/)
- CI/CD groundwork
- Documentation framework

---

## Version History

- **0.1.0** (2026-02-17): CI pipeline established, comprehensive test coverage
- **0.0.2** (2026-02-08): Backend API and Dashboard deployed to production
- **0.0.1** (2025-12-15): Project kickoff, architecture design, auth gateway
  development

[Unreleased]: https://github.com/yourusername/requiems-api/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yourusername/requiems-api/compare/v0.0.2...v0.1.0
[0.0.2]: https://github.com/yourusername/requiems-api/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/yourusername/requiems-api/releases/tag/v0.0.1
