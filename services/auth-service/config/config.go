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
	RedisURL    string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		Port:        getEnv("PORT"),
		DatabaseURL: getEnv("DATABASE_URL"),
		JwtSecret:   getEnv("JWT_SECRET"),
		RedisURL:    getEnv("REDIS_URL"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
