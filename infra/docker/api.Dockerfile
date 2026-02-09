FROM golang:1.23-alpine AS build

WORKDIR /app

# Copy only deps first (cache-friendly)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/api .

FROM alpine:3.20

# Security hardening
RUN adduser -D appuser

WORKDIR /app
ENV PORT=8080

COPY --from=build /out/api /app/api
COPY infra/migrations /app/infra/migrations

# Run as non-root
USER appuser

EXPOSE 8080

CMD ["/app/api"]


