package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthServiceUrl  string
	FileServiceUrl  string
	Port            string
	FrontendBaseUrl string
	JwtSecret       string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		AuthServiceUrl:  getEnv("AUTH_SERVICE_URL", "http://localhost:8001"),
		FileServiceUrl:  getEnv("FILE_SERVICE_URL", "http://localhost:8002"),
		Port:            getEnv("PORT", ":8000"),
		FrontendBaseUrl: getEnv("FRONTEND_BASE_URL", "http://localhost:5173"),
		JwtSecret:       getEnv("JWT_SECRET", "Aditya_Nagpal"),
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
