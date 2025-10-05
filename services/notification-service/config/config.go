package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	AWSAccessKeyId       string
	AWSSecretAccessKey   string
	AWSRegion            string
	MailjetApiKeyPublic  string
	MailjetApiKeyPrivate string
	MailjetSenderEmail   string
	MailjetSenderName    string
	SqsQueueUrl          string
	SqsWaitTimeSeconds   int
	SqsMaxMessages       int
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		Port:                 getEnvAsString("PORT", ":8010"),
		AWSAccessKeyId:       mustGetEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:   mustGetEnv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:            getEnvAsString("AWS_REGION", "ap-south-1"),
		MailjetApiKeyPublic:  mustGetEnv("MAILJEY_API_KEY_PUBLIC"),
		MailjetApiKeyPrivate: mustGetEnv("MAILJEY_API_KEY_PRIVATE"),
		MailjetSenderEmail:   mustGetEnv("MAILJET_SENDER_EMAIL"),
		MailjetSenderName:    mustGetEnv("MAILJET_SENDER_NAME"),
		SqsQueueUrl:          mustGetEnv("SQS_QUEUE_URL"),
		SqsWaitTimeSeconds:   getEnvAsInt("SQS_WAIT_TIME_SECONDS", 20),
		SqsMaxMessages:       getEnvAsInt("SQS_MAX_MESSAGES", 5),
	}
}

func getEnvAsString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("environment variable %s is required", key)
	return ""
}
