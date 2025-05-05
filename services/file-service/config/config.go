package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	BucketName         string
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	ASWRegion          string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		Port:               getEnv("PORT", ":8002"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:bhaibhai10@localhost:5432/FastFiles"),
		BucketName:         getEnv("BUCKET_NAME", "fastfiles-bucket"),
		AWSAccessKeyId:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY_ID", ""),
		ASWRegion:          getEnv("AWS_REGION", "ap-south-1"),
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
