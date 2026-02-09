package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        envOrDefault("PORT", "8080"),
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"),
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
