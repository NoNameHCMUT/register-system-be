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
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	ServerPort       string
	UploadDir        string
	MaxUploadSize    int64
	SMTPHost         string
	SMTPPort         string
	SMTPUser         string
	SMTPPassword     string
	SMTPFrom         string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "register_system"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTAccessExpiry:  mustParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
		JWTRefreshExpiry: mustParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		UploadDir:        getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadSize:    mustParseInt64(getEnv("MAX_UPLOAD_SIZE", "5242880")),
		SMTPHost:         getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:         getEnv("SMTP_PORT", "587"),
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPassword:     getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:         getEnv("SMTP_FROM", ""),
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

func mustParseInt64(s string) int64 {
	var v int64
	for _, c := range s {
		if c < '0' || c > '9' {
			panic("invalid integer: " + s)
		}
		v = v*10 + int64(c-'0')
	}
	return v
}
