package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	RedisAddr  string
	RedisPass  string
	JWTSecret  string
	GinMode    string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadConfig(path string) (*Config, error) {
	// Try loading from .env file, but don't fail if it's missing
	_ = godotenv.Load(path + "/.env")

	config := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", ""),
		DBName:     getEnv("DB_NAME", ""),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:  getEnv("REDIS_PASS", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}

	if config.JWTSecret == "" {
		log.Println("Warning: JWT_SECRET is not set")
	}

	return config, nil
}
