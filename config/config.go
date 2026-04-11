package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	DBChannelBinding string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	ServerPort       string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "register_system"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		DBChannelBinding: getEnv("DB_CHANNEL_BINDING", "disable"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTAccessExpiry:  mustParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
		JWTRefreshExpiry: mustParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("invalid duration: " + s)
	}
	return d
}
