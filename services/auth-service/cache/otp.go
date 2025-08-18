package cache

import (
	"context"
	"time"
)

func SetOTP(ctx context.Context, email, code string) error {
	return RedisClient.Set(ctx, "otp:"+email, code, 5*time.Minute).Err()
}

func GetOTP(ctx context.Context, email string) (string, error) {
	return RedisClient.Get(ctx, "otp:"+email).Result()
}
