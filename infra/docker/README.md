# Docker Setup

## Development Mode (Hot Reloading)

Start all services with **hot reloading enabled**:

```bash
cd infra/docker
docker compose -f docker-compose.dev.yml up
```

### What You Get:

✅ **Go API** with Air hot reloading on port **6969**

- Changes to `.go` files automatically rebuild and restart
- Access: http://localhost:6969/healthz

✅ **Rails Dashboard** with native hot reloading on port **3000**

- Changes to Ruby files automatically reload
- Changes to views refresh on browser reload
- Access: http://localhost:3000

✅ **PostgreSQL** on port **5432**

- Database persists between restarts
- Shared between Go and Rails

✅ **Redis** on port **6379**

- For Sidekiq background jobs

✅ **Sidekiq** background worker

- Processes Rails jobs automatically

### First Time Setup:

The containers will automatically:

1. Install dependencies (Go modules, Ruby gems)
2. Run database migrations
3. Start all services

### Development Workflow:

```bash
# Start everything
docker compose -f docker-compose.dev.yml up

# Rebuild if you change dependencies (Gemfile or go.mod)
docker compose -f docker-compose.dev.yml up --build

# View logs for specific service
docker compose -f docker-compose.dev.yml logs -f api
docker compose -f docker-compose.dev.yml logs -f dashboard

# Stop everything
docker compose -f docker-compose.dev.yml down

# Stop and remove volumes (reset database)
docker compose -f docker-compose.dev.yml down -v
```

### Accessing Services:

| Service         | URL                   | Notes                                |
| --------------- | --------------------- | ------------------------------------ |
| Rails Dashboard | http://localhost:3000 | Sign up, sign in, dashboard          |
| Go API          | http://localhost:6969 | Internal API endpoints               |
| PostgreSQL      | localhost:5432        | User: `requiem`, Password: `requiem` |
| Redis           | localhost:6379        | For Sidekiq                          |

### Hot Reloading:

**Go (Air):**

- Edit any `.go` file
- Air detects changes automatically
- Rebuilds and restarts the server (~2-3 seconds)

**Rails:**

- Edit controllers, models, views
- Rails reloads code automatically (no restart needed)
- Refresh browser to see view changes

### Running Commands Inside Containers:

```bash
# Rails console
docker compose -f docker-compose.dev.yml exec dashboard rails console

# Go build/test
docker compose -f docker-compose.dev.yml exec api go test ./...

# Database migrations
docker compose -f docker-compose.dev.yml exec dashboard rails db:migrate

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

### Troubleshooting:

**Port already in use:**

```bash
# Find and kill process using port 3000
lsof -ti:3000 | xargs kill -9

# Or change port in docker-compose.dev.yml
ports:
  - "3001:3000"  # Access on localhost:3001 instead
```

**Database connection errors:**

```bash
# Wait for PostgreSQL to fully start (check logs)
docker compose -f docker-compose.dev.yml logs db

# If still failing, restart database
docker compose -f docker-compose.dev.yml restart db
```

**Bundle install errors (Rails):**

```bash
# Clear bundle cache and reinstall
docker compose -f docker-compose.dev.yml down
docker volume rm requiem-dev_bundle_cache
docker compose -f docker-compose.dev.yml up --build
```

**Go module errors:**

```bash
# Run go mod tidy
docker compose -f docker-compose.dev.yml exec api go mod tidy

# Or rebuild
docker compose -f docker-compose.dev.yml up --build api
```

---

## Production Mode

For production deployment:

```bash
cd infra/docker
docker compose up --build
```

This uses optimized Dockerfiles with:

- Multi-stage builds
- Compiled binaries
- No development dependencies
- Caddy reverse proxy (optional, prod profile)

### Production Differences:

- Go: Compiled to static binary (no Air)
- Rails: Uses Thruster + Puma in production mode
- All code baked into images (no volume mounts)
- Requires `SECRET_KEY_BASE` and other env vars

---

## Tips:

1. **Keep dev running:** Leave `docker-compose.dev.yml` running while you code.
   All changes are picked up automatically.

2. **Use separate terminals:** Run Docker in one terminal, keep another open for
   git commands.

3. **Database GUI:** Connect with TablePlus, DBeaver, or pgAdmin using:
   - Host: `localhost`
   - Port: `5432`
   - Database: `requiem`
   - User: `requiem`
   - Password: `requiem`

4. **Reset database:**

   ```bash
   docker compose -f docker-compose.dev.yml down -v
   docker compose -f docker-compose.dev.yml up
   ```

5. **Bundle cache:** The `bundle_cache` volume speeds up subsequent starts by
   caching gems.
