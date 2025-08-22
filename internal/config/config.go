package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	JWT      JWTConfig
	OTP      OTPConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret string
}

type OTPConfig struct {
	ExpiryTime  time.Duration
	MaxAttempts int
	RateWindow  time.Duration
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "go_auth"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		},
		OTP: OTPConfig{
			ExpiryTime:  2 * time.Minute,
			MaxAttempts: 3,
			RateWindow:  10 * time.Minute,
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
