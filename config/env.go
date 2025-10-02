package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	AppEnv     string
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	JWTExpire  string
	LogLevel   string
	LogPath    string
}

var AppConfig *Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		AppName:    getEnv("APP_NAME", ""),
		AppEnv:     getEnv("APP_ENV", ""),
		AppPort:    getEnv("APP_PORT", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		JWTExpire:  getEnv("JWT_EXPIRE", ""),
		LogLevel:   getEnv("LOG_LEVEL", ""),
		LogPath:    getEnv("LOG_FILE_PATH", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}