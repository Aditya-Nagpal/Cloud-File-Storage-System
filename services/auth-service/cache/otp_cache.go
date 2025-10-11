package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/models"
	"github.com/redis/go-redis/v9"
)

// Redis key helpers
func flowKey(flowId string) string     { return "pwdreset:flow:" + flowId }
func activeKey(email string) string    { return "pwdreset:active:" + email }
func rateIpKey(ip string) string       { return "pwdreset:rate:ip:" + ip }
func rateEmailKey(email string) string { return "pwdreset:rate:email:" + email }

// SaveFlow saves the flow struct as JSON in Redis with TTL
func SaveFlow(ctx context.Context, flow *models.OTPFlow, ttl time.Duration) error {
	b, err := json.Marshal(flow)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, flowKey(flow.FlowID), b, ttl).Err()
}

// GetFlow loads a flow by id (returns nil, nil if not found)
func GetFlow(ctx context.Context, flowId string) (*models.OTPFlow, error) {
	val, err := RedisClient.Get(ctx, flowKey(flowId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var flow models.OTPFlow
	if err := json.Unmarshal([]byte(val), &flow); err != nil {
		return nil, err
	}
	return &flow, nil
}

// ExpireFlowMark sets status to EXPIRED for existing flowID (non-destructive)
func ExpireFlow(ctx context.Context, flowId string, ttl time.Duration) error {
	f, err := GetFlow(ctx, flowId)
	if err != nil || f == nil {
		return err
	}
	f.Status = "EXPIRED"
	b, _ := json.Marshal(f)
	return RedisClient.Set(ctx, flowKey(flowId), b, ttl).Err()
}

func DeleteFlow(ctx context.Context, flowId string) error {
	return RedisClient.Del(ctx, flowKey(flowId)).Err()
}

// SetActiveFlow sets pwdreset:active:{email} = flowID with TTL
func SetActiveFlow(ctx context.Context, email, flowId string, ttl time.Duration) error {
	return RedisClient.Set(ctx, activeKey(email), flowId, ttl).Err()
}

// GetActiveFlow returns flowID active for email, "" if none
func GetActiveFlow(ctx context.Context, email string) (string, error) {
	v, err := RedisClient.Get(ctx, activeKey(email)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

// DeleteActiveFlow removes the active pointer for an email
func DeleteActiveFlow(ctx context.Context, email string) error {
	return RedisClient.Del(ctx, activeKey(email)).Err()
}

// SetActiveTtl sets a new ttl to the active key
func SetActiveTtl(ctx context.Context, email string, ttl time.Duration) error {
	return RedisClient.Expire(ctx, activeKey(email), ttl).Err()
}

// Rate limiting helpers: increment counters with TTL and return current count
func IncrementRateIP(ctx context.Context, ip string, ttl time.Duration) (int64, error) {
	c, err := RedisClient.Incr(ctx, rateIpKey(ip)).Result()
	if err != nil {
		return 0, err
	}
	if c == 1 {
		_ = RedisClient.Expire(ctx, rateIpKey(ip), ttl).Err()
	}
	return c, nil
}

func IncrementRateEmail(ctx context.Context, email string, ttl time.Duration) (int64, error) {
	c, err := RedisClient.Incr(ctx, rateEmailKey(email)).Result()
	if err != nil {
		return 0, err
	}
	if c == 1 {
		_ = RedisClient.Expire(ctx, rateEmailKey(email), ttl).Err()
	}
	return c, nil
}
