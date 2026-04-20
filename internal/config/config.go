package config

import (
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
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

// LoadDotEnv loads environment variables from a .env file
func LoadDotEnv(path string) {
	file, err := os.Open(path + "/.env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func LoadConfig(path string) (*Config, error) {
	// Try loading from .env file, but don't fail if it's missing
	LoadDotEnv(path)

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
