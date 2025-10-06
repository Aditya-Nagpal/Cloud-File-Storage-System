package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"strings"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/cache"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/services/sqs"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// request DTO
type ForgotStartRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Response: always generic
var genericResp = gin.H{"message": "If an account exists for this email, an OTP has been sent if allowed."}

func StartPasswordReset(c *gin.Context) {
	ipPwdResetRateLimit, _ := strconv.Atoi(config.AppConfig.IpPwdResetRateLimit)
	emailPwdResetRateLimit, _ := strconv.Atoi(config.AppConfig.EmailPwdResetRateLimit)
	otpRateLimitTtlInMinutes, _ := strconv.Atoi(config.AppConfig.OtpRateLimitTtlInMinutes)
	otpValidityTtlInMinutes, _ := strconv.Atoi(config.AppConfig.OtpValidityTtlInMinutes)
	cooldownTtlInSeconds, _ := strconv.Atoi(config.AppConfig.CooldownTtlInSeconds)
	pwdResetPepper := config.AppConfig.PwdResetPepper

	ctx := c.Request.Context()

	var req ForgotStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Step 1: Basic rate limiting: IP and email
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	// Step 2: small windows for basic protection
	ipCount, err := cache.IncrementRateIP(ctx, ip, 1*time.Minute)
	if err == nil && ipCount > int64(ipPwdResetRateLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests from this IP: " + ip})
		return
	}
	emailCount, err := cache.IncrementRateEmail(ctx, email, time.Duration(otpRateLimitTtlInMinutes)*time.Minute)
	if err == nil && emailCount > int64(emailPwdResetRateLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests for this email: " + email})
		return
	}

	// Step 3: If an active flow exists, expire it
	activeFlow, err := cache.GetActiveFlow(ctx, email)
	if err != nil {
		log.Printf("redis get active error: %v", err)
	}
	if activeFlow != "" {
		_ = cache.ExpireFlow(ctx, activeFlow, time.Duration(otpValidityTtlInMinutes)*time.Minute)
	}

	// Step 4: Create new flow
	flowId := uuid.New().String()
	log.Println("flowId: ", flowId)
	otp, err := utils.GenerateOtp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating OTP", "error": err.Error()})
		return
	}

	// Step 5: create salt and hash
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating salt", "error": err.Error()})
		return
	}
	otpHash := utils.HashOTP(otp, salt, pwdResetPepper)
	now := time.Now().UTC()
	ttl := time.Duration(otpValidityTtlInMinutes) * (time.Minute)
	cooldownTtl := time.Duration(cooldownTtlInSeconds) * (time.Second)

	flow := &cache.OTPFlow{
		FlowID:        flowId,
		Email:         email,
		Otp:           otp,
		OtpHash:       otpHash,
		OtpSalt:       salt,
		OtpExpiresAt:  now.Add(ttl),
		ResendCount:   0,
		Attempts:      0,
		Status:        "ACTIVE",
		CreatedAt:     now,
		CooldownUntil: now.Add(cooldownTtl),
	}

	// Step 6: Save flow and active pointer
	if err := cache.SaveFlow(ctx, flow, ttl); err != nil {
		log.Printf("redis save flow error: %v", err.Error())
	}
	if err := cache.SetActiveFlow(ctx, email, flowId, ttl); err != nil {
		log.Printf("redis save active flow error: %v", err.Error())
	}

	// Step 7: Insert audit STARTED (non-sensitive)
	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "PENDING", ip, userAgent, "", flow.Attempts)

	// Step 8: Async send OTP via notification-service (fire and forget)
	if err := sqs.PublishOTP(ctx, email, otp, flowId); err != nil {
		log.Println("Failed to push notification message:", err)
	}

	// Step 9: Always return generic response (do not leak existence)
	c.JSON(http.StatusOK, genericResp)
}

type ResendOTPRequest struct {
	FlowID string `json:"flowId" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

func ResendForgotPassword(c *gin.Context) {
	maxOtpResends, _ := strconv.Atoi(config.AppConfig.MaxOtpResends)
	otpValidityTtlInMinutes, _ := strconv.Atoi(config.AppConfig.OtpValidityTtlInMinutes)
	otpCancelledTtlInMinutes, _ := strconv.Atoi(config.AppConfig.OtpCancelledTtlInMinutes)
	cooldownTtlInSeconds, _ := strconv.Atoi(config.AppConfig.CooldownTtlInSeconds)
	pwdResetPepper := config.AppConfig.PwdResetPepper

	ctx := context.Background()

	var req ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	flowId := strings.TrimSpace(req.FlowID)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Step 1: Validate active flow
	activeFlowID, err := cache.GetActiveFlow(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if activeFlowID == "" || activeFlowID != flowId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid or expired flow"})
		return
	}

	// Step 2: Fetch flow data
	flow, err := cache.GetFlow(ctx, flowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if flow == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "flow not found or expired"})
		return
	}

	// Step 3: Check status
	if flow.Status != "ACTIVE" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Flow is not active"})
		return
	}

	now := time.Now().UTC()

	// Step 4: Cooldown check
	if !flow.CooldownUntil.IsZero() && now.Before(flow.CooldownUntil) {
		remaining := max(int(flow.CooldownUntil.Sub(now).Seconds()), 0)
		c.JSON(http.StatusTooManyRequests, gin.H{"message": fmt.Sprintf("Please wait %d seconds requesting another OTP", remaining)})
		return
	}

	// Step 5: Resend count check
	if flow.ResendCount >= maxOtpResends {
		cancelTtl := time.Duration(otpCancelledTtlInMinutes) * (time.Minute)
		// Mark flow cancelled and delete active pointer
		flow.Status = "CANCELLED"
		// Save flow with small TTL so record remains for audit (keep a few minutes)
		if err := cache.SaveFlow(ctx, flow, cancelTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		// remove active pointer
		if err := cache.DeleteActiveFlow(ctx, email); err != nil {
			log.Printf("redis delete active flow error: %v", err.Error())
		}

		// Insert audit CANCELLED
		failureReason := "reset attempts exceeded"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "CANCELLED", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusTooManyRequests, gin.H{"message": "maximum resend attempts exceeded"})
		return
	}

	// Step 6: Generate new OTP, salt, hash
	otp, err := utils.GenerateOtp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating OTP", "error": err})
		return
	}

	salt, err := utils.GenerateSalt(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating salt", "error": err.Error()})
		return
	}
	otpHash := utils.HashOTP(otp, salt, pwdResetPepper)
	otpTtl := time.Duration(otpValidityTtlInMinutes) * (time.Minute)
	cooldownTtl := time.Duration(cooldownTtlInSeconds) * (time.Second)

	flow.Otp = otp
	flow.OtpHash = otpHash
	flow.OtpSalt = salt
	flow.OtpExpiresAt = now.Add(otpTtl)
	flow.ResendCount = flow.ResendCount + 1
	flow.CooldownUntil = now.Add(cooldownTtl)
	flow.Status = "ACTIVE"

	// Step 7: Compute TTL for Redis save (keep flow alive until otp expiry)
	ttl := time.Until(flow.OtpExpiresAt)
	if ttl <= 0 {
		ttl = otpTtl
	}

	if err := cache.SaveFlow(ctx, flow, ttl); err != nil {
		log.Printf("redis error saving flow: %v", err.Error())
	}

	if err := cache.SetActiveFlow(ctx, email, flowId, ttl); err != nil {
		log.Printf("redis save active flow error: %v", err.Error())
	}

	failureReason := "otp missed or expired"
	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "RESENT", ip, userAgent, failureReason, flow.Attempts)

	// Step 8: Async send OTP via notification-service (fire and forget)
	if err := sqs.PublishOTP(ctx, email, otp, flowId); err != nil {
		log.Println("Failed to push notification message:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "OTP resent successfully.",
		"cooldownSeconds": cooldownTtlInSeconds,
	})
}

type VerifyOTPRequest struct {
	FlowID string `json:"flowId" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	OTP    string `json:"otp" binding:"required"`
}

func VerifyForgotPasswordOTP(c *gin.Context) {
	maxOtpAttempts, _ := strconv.Atoi(config.AppConfig.MaxOtpAttempts)
	pwdFlowTtlInMinutes, _ := strconv.Atoi(config.AppConfig.PwdFlowTtlInMinutes)
	pwdResetPepper := config.AppConfig.PwdResetPepper

	ctx := context.Background()

	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	flowId := strings.TrimSpace(req.FlowID)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	otp := strings.TrimSpace(req.OTP)

	if len(otp) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid OTP"})
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Step 1: Validate active flow
	activeFlowID, err := cache.GetActiveFlow(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if activeFlowID == "" || activeFlowID != flowId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid or expired flow"})
		return
	}

	// Step 2: Fetch flow data
	flow, err := cache.GetFlow(ctx, flowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if flow == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "flow not found or expired"})
		return
	}

	// Step 3: Check status
	if flow.Status != "ACTIVE" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Flow is not active"})
		return
	}

	now := time.Now().UTC()
	pwdFlowTtl := time.Duration(pwdFlowTtlInMinutes) * (time.Minute)

	// Step 4: Expiry check
	if now.After(flow.OtpExpiresAt) {
		flow.Status = "EXPIRED"

		if err := cache.SaveFlow(ctx, flow, pwdFlowTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		if err := cache.DeleteActiveFlow(ctx, email); err != nil {
			log.Printf("redis delete active flow error: %v", err.Error())
		}

		failureReason := "otp expired"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "EXPIRED", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusBadRequest, gin.H{"message": "OTP expired"})
		return
	}

	// Step 5: Attempt count check
	if flow.Attempts >= maxOtpAttempts {
		flow.Status = "BLOCKED"

		if err := cache.SaveFlow(ctx, flow, pwdFlowTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		if err := cache.DeleteActiveFlow(ctx, email); err != nil {
			log.Printf("redis delete active flow error: %v", err.Error())
		}

		failureReason := "verify attempts exceeded"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "BLOCKED", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusBadRequest, gin.H{"message": "Too many failed attempts"})
		return
	}

	calculatedHash := utils.HashOTP(otp, flow.OtpSalt, pwdResetPepper)

	if calculatedHash != flow.OtpHash {
		flow.Attempts++
		if flow.Attempts >= maxOtpAttempts {
			flow.Status = "BLOCKED"

			if err := cache.SaveFlow(ctx, flow, pwdFlowTtl); err != nil {
				log.Printf("redis save flow error: %v", err.Error())
			}

			if err := cache.DeleteActiveFlow(ctx, email); err != nil {
				log.Printf("redis delete active flow error: %v", err.Error())
			}

			failureReason := "verify attempts limit exceeded"
			_ = db.InsertPasswordResetAudit(ctx, flowId, email, "BLOCKED", ip, userAgent, failureReason, flow.Attempts)

			c.JSON(http.StatusTooManyRequests, gin.H{"message": "OTP blocked due to repeated failed attempts"})
			return
		}

		if err := cache.SaveFlow(ctx, flow, pwdFlowTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		failureReason := "incorrect otp"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "FAILED_ATTEMPT", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}

	// Step 7: OTP verified successfully
	flow.Status = "VERIFIED"

	if err := cache.SaveFlow(ctx, flow, pwdFlowTtl); err != nil {
		log.Printf("redis save flow error: %v", err.Error())
	}

	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "VERIFIED", ip, userAgent, "", flow.Attempts)

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
