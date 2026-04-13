# Code Quality

Commands for maintaining code quality across the monorepo.

## Go API (apps/api)

All commands run inside the Docker container:

```bash
docker exec requiem-dev-api-1 sh -lc 'export PATH=/usr/local/go/bin:$PATH; cd /app; go fmt ./... && /app/bin/golangci-lint run --fix'
```

This runs:
- `go fmt` - Formats all Go files
- `golangci-lint run --fix` - Lints and auto-fixes issues

Run tests before pushing:

```bash
docker exec requiem-dev-api-1 go test ./...
```

Run tests with coverage:

```bash
docker exec requiem-dev-api-1 go test -race -coverprofile=coverage.out ./...
```

Run specific test:

```bash
docker exec requiem-dev-api-1 go test ./services/text/advice -v -run TestGetAdvice
```
