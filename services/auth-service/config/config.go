package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JwtSecret   string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		Port:        getEnv("PORT", ":8001"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:bhaibhai10@localhost:5432/FastFiles"),
		JwtSecret:   getEnv("JWT_SECRET", "Aditya_Nagpal"),
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
