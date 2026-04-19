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
