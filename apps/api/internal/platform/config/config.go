package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	BackendSecret string
	RedisURL      string
}

func Load() Config {
	return Config{
		Port:          envOrDefault("PORT", "8080"),
		DatabaseURL:   envOrDefault("DATABASE_URL", "postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"),
		BackendSecret: envOrDefault("BACKEND_SECRET", ""),
		RedisURL:      envOrDefault("REDIS_URL", "redis://localhost:6379/0"),
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
