package models

import "time"

// OTPFlow represents the Redis-stored flow structure (JSON stored as value)
type OTPFlow struct {
	FlowID        string    `json:"flow_id"`
	Email         string    `json:"email"`
	Otp           string    `json:"otp"`
	OtpHash       string    `json:"otp_hash"`
	OtpSalt       string    `json:"otp_salt"`
	OtpExpiresAt  time.Time `json:"otp_expires_at"`
	ResendCount   int       `json:"resend_count"`
	Attempts      int       `json:"attempts"`
	Status        string    `json:"status"` // ACTIVE | EXPIRED | COMPLETED | CANCELLED
	CreatedAt     time.Time `json:"created_at"`
	CooldownUntil time.Time `json:"cooldown_until"`
}
