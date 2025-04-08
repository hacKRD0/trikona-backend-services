package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BASE_URL string
	PORT string
	DB_HOST string
	DB_PORT int
	DB_USER string
	DB_PASSWORD string
	DB_NAME string
	JWT_SECRET string
	JWT_RESET_SECRET string
	JWT_REGISTRATION_SECRET string
	LINKEDIN_CLIENT_ID string
	LINKEDIN_CLIENT_SECRET string
	LINKEDIN_REDIRECT_URI string
	MAILJET_API_KEY string
	MAILJET_API_SECRET string
	EMAIL_FROM string
	FRONTEND_BASE_URL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		BASE_URL: getEnv("BASE_URL", "http://localhost:5173"),
		PORT: getEnv("PORT", "8080"),
		DB_HOST: getEnv("DB_HOST", "localhost"),
		DB_PORT: getEnvAsInt("DB_PORT", 5432),	
		DB_USER: getEnv("DB_USER", "postgres"),	
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),	
		DB_NAME: getEnv("DB_NAME", "postgres"),
		JWT_SECRET: getEnv("JWT_SECRET", ""),	
		JWT_RESET_SECRET: getEnv("JWT_RESET_SECRET", ""),
		JWT_REGISTRATION_SECRET: getEnv("JWT_REGISTRATION_SECRET", ""),
		LINKEDIN_CLIENT_ID: getEnv("LINKEDIN_CLIENT_ID", ""),
		LINKEDIN_CLIENT_SECRET: getEnv("LINKEDIN_CLIENT_SECRET", ""),
		LINKEDIN_REDIRECT_URI: getEnv("LINKEDIN_REDIRECT_URI", ""),
		MAILJET_API_KEY: getEnv("MAILJET_API_KEY", ""),
		MAILJET_API_SECRET: getEnv("MAILJET_API_SECRET", ""),
		EMAIL_FROM: getEnv("EMAIL_FROM", ""),
		FRONTEND_BASE_URL: getEnv("FRONTEND_BASE_URL", "http://localhost:3000"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
    valueStr := getEnv(name, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
			return value
    }

    return defaultVal
}