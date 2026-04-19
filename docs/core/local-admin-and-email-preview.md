## Local Admin + Email Preview

Use this guide for two common local-dev tasks in the Rails dashboard:

- Promote a specific user to admin
- Open development emails in Letter Opener Web

### 1. Make a user admin locally

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

### 2. Letter Opener Web in development

Route is mounted only in development at:

- `http://localhost:3000/letter_opener`

How it works:

- `config/environments/development.rb` uses `:letter_opener_web`
- `config/routes.rb` mounts `LetterOpenerWeb::Engine` at `/letter_opener` in development

If page is empty, trigger an email flow (signup confirmation, password reset, etc.), then refresh `/letter_opener`.
