## 1. Make a user admin locally

Target email:

- `eliaz.bobadilladeva@gmail.com`

Run from host machine (inside `apps/dashboard`):

```bash
cd apps/dashboard
bin/rails runner "user = User.find_by!(email: 'eliaz.bobadilladeva@gmail.com'); user.update!(admin: true); puts \"#{user.email} admin=#{user.admin?}\""
```

Run in Docker container:

```bash
docker exec requiem-dev-dashboard-1 bin/rails runner "user = User.find_by!(email: 'eliaz.bobadilladeva@gmail.com'); user.update!(admin: true); puts \"#{user.email} admin=#{user.admin?}\""
```

Verify in Rails console (optional):

```bash
cd apps/dashboard
bin/rails console
```

Then:

```ruby
User.find_by!(email: "eliaz.bobadilladeva@gmail.com").admin?
```

Set back to non-admin if needed:

```bash
cd apps/dashboard
bin/rails runner "user = User.find_by!(email: 'eliaz.bobadilladeva@gmail.com'); user.update!(admin: false); puts \"#{user.email} admin=#{user.admin?}\""
```

### Persistence Notes

Local users should remain after normal restarts because the PostgreSQL data
lives in the named Docker volume `db_data`.

- `docker compose stop` keeps users
- `docker compose restart` keeps users
- `docker compose down` keeps users
- `docker compose down -v` deletes users because it removes the volume

## 2. Docker Data Persistence

### Stopping Containers (Preserves Data)

To stop containers without losing database data:

```bash
cd infra/docker
docker compose stop          # Stops all containers (data persists)
docker compose start         # Restarts containers with existing data
```

### Restarting Containers (Preserves Data)

```bash
cd infra/docker
docker compose restart       # Restarts all services (data persists)
```

### Full Cleanup (Deletes Everything)

⚠️ **WARNING**: This deletes all database data, volumes, and test users:

```bash
cd infra/docker
docker compose down -v       # Removes containers AND volumes (data lost!)
docker compose up -d         # Spins up fresh containers with empty database
docker compose logs -f       # Watch logs as services initialize and seed
```

### When to Use Each

| Command                  | Data Persists? | When to Use                      |
| ------------------------ | -------------- | -------------------------------- |
| `docker compose stop`    | ✅ Yes         | Quick pause, testing later       |
| `docker compose restart` | ✅ Yes         | Quick restart, debugging         |
| `docker compose down`    | ✅ Yes         | Full shutdown, keep data         |
| `docker compose down -v` | ❌ **No**      | Fresh start, clear all test data |

### Automatic Test User Seeding

The dashboard service automatically creates test users on startup via `db:seed`:

- **Email**: `eliaz.bobadilladeva@gmail.com` (admin: true)
- **Email**: `test@example.com` (regular user)

If containers are cleanly shut down with `docker compose stop` (not `down -v`),
these users persist. After `docker compose down -v`, they're recreated
automatically on the next `up`.

### Check Database Status

```bash
# View PostgreSQL logs
docker compose logs db

# Connect to database directly
docker exec -it requiem-dev-db-1 psql -U requiem -d requiem

# List users
docker exec requiem-dev-db-1 psql -U requiem -d requiem -c "SELECT email, admin FROM users;"

# Seed data manually if needed
docker exec requiem-dev-dashboard-1 bin/rails db:seed
```
