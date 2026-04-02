package config

import (
	"log"
	"os"
	"strconv"

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
	DBMaxIdleConns int
	DBMaxOpenConns int
	DBMaxLifetime  int // in minutes
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
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
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 3),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxLifetime:  getEnvAsInt("DB_MAX_LIFETIME", 30),
	}

	if config.JWTSecret == "" {
		log.Println("Warning: JWT_SECRET is not set")
	}

	return config, nil
}
