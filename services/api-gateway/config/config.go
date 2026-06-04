package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	AuthServiceUrl  string
	FileServiceUrl  string
	UserServiceUrl  string
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
		Port:            getEnv("PORT"),
		AuthServiceUrl:  getEnv("AUTH_SERVICE_URL"),
		FileServiceUrl:  getEnv("FILE_SERVICE_URL"),
		UserServiceUrl:  getEnv("USER_SERVICE_URL"),
		FrontendBaseUrl: getEnv("FRONTEND_BASE_URL"),
		JwtSecret:       getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
