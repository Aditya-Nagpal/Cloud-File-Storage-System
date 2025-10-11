package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JwtSecret   string
	RedisURL    string

	IpPwdResetRateLimit      int
	EmailPwdResetRateLimit   int
	OtpRateLimitTtlInMinutes int

	OtpValidityTtlInMinutes     int
	OtpExpiryBufferTtlInMinutes int

	PwdResetTtlInMinutes              int
	ResetFlowCancelBufferTtlInMinutes int
	ResetFlowBlockBufferTtlInMinutes  int
	CooldownTtlInSeconds              int

	PwdResetPepper string

	MaxOtpResends  int
	MaxOtpAttempts int

	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string

	SqsQueueUrl string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		Port:        mustGetEnv("PORT"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		JwtSecret:   mustGetEnv("JWT_SECRET"),
		RedisURL:    mustGetEnv("REDIS_URL"),

		IpPwdResetRateLimit:      getEnvAsInt("IP_PWD_RESET_RATE_LIMIT"),
		EmailPwdResetRateLimit:   getEnvAsInt("EMAIL_PWD_RESET_RATE_LIMIT"),
		OtpRateLimitTtlInMinutes: getEnvAsInt("OTP_RATE_LIMIT_TTL_IN_MINUTES"),

		OtpValidityTtlInMinutes:     getEnvAsInt("OTP_VALIDITY_TTL_IN_MINUTES"),
		OtpExpiryBufferTtlInMinutes: getEnvAsInt("OTP_EXPIRY_BUFFER_TTL_IN_MINUTES"),

		PwdResetTtlInMinutes:              getEnvAsInt("PWD_RESET_TTL_IN_MINUTES"),
		ResetFlowCancelBufferTtlInMinutes: getEnvAsInt("RESET_FLOW_CANCEL_BUFFER_TTL_IN_MINUTES"),
		ResetFlowBlockBufferTtlInMinutes:  getEnvAsInt("RESET_FLOW_BLOCK_BUFFER_TTL_IN_MINUTES"),
		CooldownTtlInSeconds:              getEnvAsInt("COOLDOWN_TTL_IN_SECONDS"),

		PwdResetPepper: getEnvAsString("PWD_RESET_PEPPER"),

		MaxOtpResends:  getEnvAsInt("MAX_OTP_RESENDS"),
		MaxOtpAttempts: getEnvAsInt("MAX_OTP_ATTEMPTS"),

		AWSAccessKeyId:     mustGetEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: mustGetEnv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          mustGetEnv("AWS_REGION"),

		SqsQueueUrl: mustGetEnv("SQS_QUEUE_URL"),
	}
}

func getEnvAsString(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return ""
}

func getEnvAsInt(key string) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 0
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("environment variable %s is required", key)
	return ""
}
