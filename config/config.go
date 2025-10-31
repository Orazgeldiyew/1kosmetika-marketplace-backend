package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
	FromEmail  string
}

func Load() *Config {

	_ = godotenv.Load()

	return &Config{
		DBHost:     mustEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     mustEnv("DB_USER"),
		DBPassword: mustEnv("DB_PASSWORD"),
		DBName:     mustEnv("DB_NAME"),
		JWTSecret:  mustEnv("JWT_SECRET"),
		ServerPort: getEnv("PORT", "8080"),
		SMTPHost:   getEnv("SMTP_HOST", ""),
		SMTPPort:   getEnv("SMTP_PORT", ""),
		SMTPUser:   getEnv("SMTP_USER", ""),
		SMTPPass:   getEnv("SMTP_PASS", ""),
		FromEmail:  getEnv("FROM_EMAIL", ""),
	}
}

func mustEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("missing required env: " + key)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
