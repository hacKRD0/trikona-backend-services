package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	return godotenv.Load(".env.prod")
}

// GetEnv returns the value of the environment variable or the default value if not set
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// LinkedInConfig holds LinkedIn OAuth configuration
type LinkedInConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	TokenURL     string
	ProfileURL   string
}

// LoadLinkedInConfig loads LinkedIn OAuth configuration from environment variables
func LoadLinkedInConfig() *LinkedInConfig {
	return &LinkedInConfig{
		ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
		ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("LINKEDIN_REDIRECT_URI"),
		TokenURL:     "https://www.linkedin.com/oauth/v2/accessToken",
		ProfileURL:   "https://api.linkedin.com/v2/me?projection=(id,firstName,lastName,profilePicture(displayImage~:playableStreams))",
	}
} 