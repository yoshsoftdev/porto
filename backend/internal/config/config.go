package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port          string
	DatabaseURL   string
	SQLitePath    string
	AdminPassword string
	JWTSecret     string
	EncryptionKey string
	CORSOrigins   []string
	MaxUploadMB   int64
}

func Load() Config {
	c := Config{
		Port:          env("PORT", "8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		SQLitePath:    env("SQLITE_PATH", "porto.db"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
		MaxUploadMB:   int64(envInt("MAX_UPLOAD_MB", 5)),
	}
	for _, o := range strings.Split(env("CORS_ORIGINS", "http://localhost:5173"), ",") {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			c.CORSOrigins = append(c.CORSOrigins, o)
		}
	}
	return c
}

func env(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
		return v
	}
	return def
}
