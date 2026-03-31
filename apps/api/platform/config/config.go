package config

import "os"

type Config struct {
	Port               string
	DatabaseURL        string
	BackendSecret      string
	RedisURL           string
	VPNDatabasePath    string
	VPNASNDatabasePath string
	IPCityDatabasePath string
}

func Load() Config {
	return Config{
		Port:               envOrDefault("PORT", "8080"),
		DatabaseURL:        envOrDefault("DATABASE_URL", "postgres://requiem:requiem@localhost:5432/requiem?sslmode=disable"),
		BackendSecret:      envOrDefault("BACKEND_SECRET", ""),
		RedisURL:           envOrDefault("REDIS_URL", "redis://localhost:6379/0"),
		VPNDatabasePath:    envOrDefault("VPN_DATABASE_PATH", "dbs/IP2PROXY-LITE-PX2.BIN"),
		VPNASNDatabasePath: envOrDefault("VPN_ASN_DATABASE_PATH", "dbs/GeoLite2-ASN.mmdb"),
		IPCityDatabasePath: envOrDefault("IP_CITY_DATABASE_PATH", "dbs/GeoLite2-City.mmdb"),
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
