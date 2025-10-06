package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                     string
	DatabaseURL              string
	JwtSecret                string
	RedisURL                 string
	IpPwdResetRateLimit      string
	EmailPwdResetRateLimit   string
	PwdFlowTtlInMinutes      string
	OtpRateLimitTtlInMinutes string
	OtpValidityTtlInMinutes  string
	OtpCancelledTtlInMinutes string
	CooldownTtlInSeconds     string
	PwdResetPepper           string
	MaxOtpResends            string
	MaxOtpAttempts           string
	AWSAccessKeyId           string
	AWSSecretAccessKey       string
	AWSRegion                string
	SqsQueueUrl              string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = &Config{
		Port:                     getEnv("PORT"),
		DatabaseURL:              getEnv("DATABASE_URL"),
		JwtSecret:                getEnv("JWT_SECRET"),
		RedisURL:                 getEnv("REDIS_URL"),
		IpPwdResetRateLimit:      getEnv("IP_PWD_RESET_RATE_LIMIT"),
		EmailPwdResetRateLimit:   getEnv("EMAIL_PWD_RESET_RATE_LIMIT"),
		PwdFlowTtlInMinutes:      getEnv("PWD_FLOW_TTL_IN_MINUTES"),
		OtpRateLimitTtlInMinutes: getEnv("OTP_RATE_LIMIT_TTL_IN_MINUTES"),
		OtpValidityTtlInMinutes:  getEnv("OTP_VALIDITY_TTL_IN_MINUTES"),
		OtpCancelledTtlInMinutes: getEnv("OTP_FLOW_CANCELLED_TTL_IN_MINUTES"),
		CooldownTtlInSeconds:     getEnv("COOLDOWN_TTL_IN_SECONDS"),
		PwdResetPepper:           getEnv("PWD_RESET_PEPPER"),
		MaxOtpResends:            getEnv("MAX_OTP_RESENDS"),
		MaxOtpAttempts:           getEnv("MAX_OTP_ATTEMPTS"),
		AWSAccessKeyId:           getEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:       getEnv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:                getEnv("AWS_REGION"),
		SqsQueueUrl:              getEnv("SQS_QUEUE_URL"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
