package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthServiceUrl  string
	Port            string
	FrontendBaseUrl string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		AuthServiceUrl:  getEnv("AUTH_SERVICE_URL", "http://localhost:8001"),
		Port:            getEnv("PORT", ":8000"),
		FrontendBaseUrl: getEnv("FRONTEND_BASE_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		os.Setenv(key, fallback)
		return fallback
	}
	return value
}
