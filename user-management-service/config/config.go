package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	DB_HOST string
	DB_PORT int
	DB_USER string
	DB_PASSWORD string
	DB_NAME string
	JWT_SECRET string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		PORT: getEnv("PORT", "8080"),
		DB_HOST: getEnv("DB_HOST", "localhost"),
		DB_PORT: getEnvAsInt("DB_PORT", 5432),	
		DB_USER: getEnv("DB_USER", "postgres"),	
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),	
		DB_NAME: getEnv("DB_NAME", "postgres"),
		JWT_SECRET: getEnv("JWT_SECRET", ""),	
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