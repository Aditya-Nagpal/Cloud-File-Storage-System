package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"

	"strings"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/cache"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/services/sqs"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/utils"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// request DTO
type ForgotStartRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func StartPasswordReset(c *gin.Context) {
	otpValidityTtlInMinutes := config.AppConfig.OtpValidityTtlInMinutes
	otpExpiryBufferTtlInMinutes := config.AppConfig.OtpExpiryBufferTtlInMinutes
	otpTotalFlowTtlInMinutes := otpValidityTtlInMinutes + otpExpiryBufferTtlInMinutes
	cooldownTtlInSeconds := config.AppConfig.CooldownTtlInSeconds

	ctx := c.Request.Context()

	var req ForgotStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Step 1: Basic rate limiting: IP and email
	ipPwdResetRateLimit := config.AppConfig.IpPwdResetRateLimit
	emailPwdResetRateLimit := config.AppConfig.EmailPwdResetRateLimit
	otpRateLimitTtlInMinutes := config.AppConfig.OtpRateLimitTtlInMinutes

	ipCount, err := cache.IncrementRateIP(ctx, ip, time.Duration(otpRateLimitTtlInMinutes)*time.Minute)
	if err == nil && ipCount > int64(ipPwdResetRateLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": fmt.Sprintf("Too many requests from this IP %s. Wait for %v minutes .", ip, otpRateLimitTtlInMinutes)})
		return
	}

	emailCount, err := cache.IncrementRateEmail(ctx, email, time.Duration(otpRateLimitTtlInMinutes)*time.Minute)
	if err == nil && emailCount > int64(emailPwdResetRateLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": fmt.Sprintf("Too many requests for this email %s. Wait for %v minutes .", email, otpRateLimitTtlInMinutes)})
		return
	}

	// Step 2: If an active flow exists, expire it
	activeFlow, err := cache.GetActiveFlow(ctx, email)
	if err != nil {
		log.Printf("redis get active error: %v", err)
	}
	if activeFlow != "" {
		_ = cache.ExpireFlow(ctx, activeFlow, time.Duration(otpValidityTtlInMinutes)*time.Minute)
	}

	// Step 3: Create new flow
	flowId := uuid.New().String()
	otp, err := utils.GenerateOtp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating OTP", "error": err.Error()})
		return
	}

	// Step 4: create salt and hash
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating salt", "error": err.Error()})
		return
	}

	otpHash := utils.HashOTP(otp, salt, config.AppConfig.PwdResetPepper)

	now := time.Now().UTC()
	otpValidityTtl := time.Duration(otpValidityTtlInMinutes) * (time.Minute)
	otpTotalFlowTtl := time.Duration(otpTotalFlowTtlInMinutes) * (time.Minute)
	cooldownTtl := time.Duration(cooldownTtlInSeconds) * (time.Second)

	flow := &models.OTPFlow{
		FlowID:        flowId,
		Email:         email,
		Otp:           otp,
		OtpHash:       otpHash,
		OtpSalt:       salt,
		OtpExpiresAt:  now.Add(otpValidityTtl),
		ResendCount:   0,
		Attempts:      0,
		Status:        "ACTIVE",
		CreatedAt:     now,
		CooldownUntil: now.Add(cooldownTtl),
	}

	// Step 5: Save flow and active pointer
	if err := cache.SaveFlow(ctx, flow, otpTotalFlowTtl); err != nil {
		log.Printf("redis save flow error: %v", err.Error())
	}
	if err := cache.SetActiveFlow(ctx, email, flowId, otpTotalFlowTtl); err != nil {
		log.Printf("redis save active flow error: %v", err.Error())
	}

	// Step 6: Insert audit STARTED (non-sensitive)
	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "PENDING", ip, userAgent, "", flow.Attempts)

	// Step 7: Async send OTP via notification-service (fire and forget)
	if err := sqs.PublishOtpEmail(ctx, email, otp, flowId); err != nil {
		log.Println("Failed to push notification message:", err)
	}

	// Step 8: Always return generic response (do not leak existence)
	c.JSON(http.StatusCreated, gin.H{
		"message": "If an account exists for this email, an OTP has been sent if allowed.",
		"flowId":  flowId,
	})
}

type ResendOTPRequest struct {
	FlowID string `json:"flowId" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

func ResendForgotPassword(c *gin.Context) {
	otpValidityTtlInMinutes := config.AppConfig.OtpValidityTtlInMinutes
	otpExpiryBufferTtlInMinutes := config.AppConfig.OtpExpiryBufferTtlInMinutes
	otpTotalFlowTtlInMinutes := otpValidityTtlInMinutes + otpExpiryBufferTtlInMinutes
	resetFlowCancelBufferTtlInMinutes := config.AppConfig.ResetFlowCancelBufferTtlInMinutes
	cooldownTtlInSeconds := config.AppConfig.CooldownTtlInSeconds

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
		c.JSON(http.StatusGone, gin.H{"message": "The password reset flow has expired. Please initiate the process again."})
		return
	}

	// Step 2: Fetch flow data
	flow, err := cache.GetFlow(ctx, flowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if flow == nil {
		c.JSON(http.StatusGone, gin.H{"message": "The password reset flow has expired. Please initiate the process again."})
		return
	}

	// Step 3: Check status
	validStatuses := []string{"ACTIVE", "EXPIRED"}
	if !slices.Contains(validStatuses, flow.Status) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "The flow is not active or expired. Please initiate the process again."})
		return
	}

	now := time.Now().UTC()

	// Step 4: Cooldown check
	if !flow.CooldownUntil.IsZero() && now.Before(flow.CooldownUntil) {
		remaining := max(int(flow.CooldownUntil.Sub(now).Seconds()), 0)
		c.JSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("Please wait %d seconds requesting another OTP", remaining)})
		return
	}

	// Step 5: Resend count check
	maxOtpResends := config.AppConfig.MaxOtpResends
	if flow.ResendCount >= maxOtpResends {
		cancelBufferTtl := time.Duration(resetFlowCancelBufferTtlInMinutes) * (time.Minute)
		// Mark flow cancelled and delete active pointer
		flow.Status = "CANCELLED"
		// Save flow with small TTL so record remains for audit (keep a few minutes)
		if err := cache.SaveFlow(ctx, flow, cancelBufferTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		// remove active pointer
		if err := cache.DeleteActiveFlow(ctx, email); err != nil {
			log.Printf("redis delete active flow error: %v", err.Error())
		}

		// Insert audit CANCELLED
		failureReason := "reset attempts exceeded"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "CANCELLED", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusTooManyRequests, gin.H{"message": "You have exceeded the maximum number of OTP resends. Please initiate the password reset process again."})
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

	otpHash := utils.HashOTP(otp, salt, config.AppConfig.PwdResetPepper)

	otpValidityTtl := time.Duration(otpValidityTtlInMinutes) * (time.Minute)
	otpTotalFlowTtl := time.Duration(otpTotalFlowTtlInMinutes) * (time.Minute)
	cooldownTtl := time.Duration(cooldownTtlInSeconds) * (time.Second)

	flow.Otp = otp
	flow.OtpHash = otpHash
	flow.OtpSalt = salt
	flow.OtpExpiresAt = now.Add(otpValidityTtl)
	flow.ResendCount = flow.ResendCount + 1
	flow.CooldownUntil = now.Add(cooldownTtl)
	flow.Status = "ACTIVE"

	if err := cache.SaveFlow(ctx, flow, otpTotalFlowTtl); err != nil {
		log.Printf("redis error saving flow: %v", err.Error())
	}

	if err := cache.SetActiveFlow(ctx, email, flowId, otpTotalFlowTtl); err != nil {
		log.Printf("redis save active flow error: %v", err.Error())
	}

	failureReason := "otp missed or expired"
	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "RESENT", ip, userAgent, failureReason, flow.Attempts)

	// Step 8: Async send OTP via notification-service (fire and forget)
	if err := sqs.PublishOtpEmail(ctx, email, otp, flowId); err != nil {
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
	resetFlowBlockBufferTtlInMinutes := config.AppConfig.ResetFlowBlockBufferTtlInMinutes
	pwdResetTtlInMinutes := config.AppConfig.PwdResetTtlInMinutes

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
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid OTP"})
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
		c.JSON(http.StatusGone, gin.H{"message": "The password reset flow has expired. Please initiate the process again."})
		return
	}

	// Step 2: Fetch flow data
	flow, err := cache.GetFlow(ctx, flowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if flow == nil {
		c.JSON(http.StatusGone, gin.H{"message": "The password reset flow has expired. Please initiate the process again."})
		return
	}

	// Step 3: Check status
	if flow.Status != "ACTIVE" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Flow is not active. Please initiate the process again."})
		return
	}

	now := time.Now().UTC()
	pwdResetTtl := time.Duration(pwdResetTtlInMinutes) * (time.Minute)

	// Step 4: Expiry check
	if now.After(flow.OtpExpiresAt) {
		flow.Status = "EXPIRED"

		failureReason := "otp expired"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "EXPIRED", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusUnauthorized, gin.H{"message": "OTP expired"})
		return
	}

	// Step 5: Attempt count check
	maxOtpAttempts := config.AppConfig.MaxOtpAttempts
	if flow.Attempts >= maxOtpAttempts {
		blockBufferTtl := time.Duration(resetFlowBlockBufferTtlInMinutes) * (time.Minute)

		flow.Status = "BLOCKED"

		if err := cache.SaveFlow(ctx, flow, blockBufferTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		if err := cache.DeleteActiveFlow(ctx, email); err != nil {
			log.Printf("redis delete active flow error: %v", err.Error())
		}

		failureReason := "verify attempts exceeded"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "BLOCKED", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusTooManyRequests, gin.H{"message": "You have exceeded the maximum number of OTP verification attempts. Please initiate the password reset process again."})
		return
	}

	calculatedHash := utils.HashOTP(otp, flow.OtpSalt, config.AppConfig.PwdResetPepper)

	if calculatedHash != flow.OtpHash {
		flow.Attempts++
		if flow.Attempts >= maxOtpAttempts {
			blockBufferTtl := time.Duration(resetFlowBlockBufferTtlInMinutes) * (time.Minute)

			flow.Status = "BLOCKED"

			if err := cache.SaveFlow(ctx, flow, blockBufferTtl); err != nil {
				log.Printf("redis save flow error: %v", err.Error())
			}

			if err := cache.DeleteActiveFlow(ctx, email); err != nil {
				log.Printf("redis delete active flow error: %v", err.Error())
			}

			failureReason := "verify attempts limit exceeded"
			_ = db.InsertPasswordResetAudit(ctx, flowId, email, "BLOCKED", ip, userAgent, failureReason, flow.Attempts)

			c.JSON(http.StatusTooManyRequests, gin.H{"message": "You have exceeded the maximum number of OTP verification attempts. Please initiate the password reset process again."})
			return
		}

		if err := cache.SaveFlow(ctx, flow, pwdResetTtl); err != nil {
			log.Printf("redis save flow error: %v", err.Error())
		}

		failureReason := "incorrect otp"
		_ = db.InsertPasswordResetAudit(ctx, flowId, email, "FAILED_ATTEMPT", ip, userAgent, failureReason, flow.Attempts)

		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}

	// Step 7: OTP verified successfully
	flow.Status = "VERIFIED"

	if err := cache.SaveFlow(ctx, flow, pwdResetTtl); err != nil {
		log.Printf("redis save flow error: %v", err.Error())
	}

	if err := cache.SetActiveTtl(ctx, email, pwdResetTtl); err != nil {
		log.Printf("redis save ttl error: %v", err.Error())
	}

	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "VERIFIED", ip, userAgent, "", flow.Attempts)

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

type ResetPasswordRequest struct {
	FlowID      string `json:"flowId" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func ResetPassword(c *gin.Context) {
	pwdResetTtlInMinutes := config.AppConfig.PwdResetTtlInMinutes

	ctx := context.Background()

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	flowId := strings.TrimSpace(req.FlowID)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	newPassword := strings.TrimSpace(req.NewPassword)

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Step 1: Validate active flow
	activeFlowID, err := cache.GetActiveFlow(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if activeFlowID == "" || activeFlowID != flowId {
		c.JSON(http.StatusGone, gin.H{"message": "The password reset flow has expired. Please initiate the process again."})
		return
	}

	// Step 2: Fetch flow data
	flow, err := cache.GetFlow(ctx, flowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	} else if flow == nil {
		c.JSON(http.StatusGone, gin.H{"message": "The password reset flow has expired. Please initiate the process again."})
		return
	}

	// Step 3: Get current hashed password and validate email
	currentHashedPassword, err := db.GetUserHashedPassword(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in checking email", "error": err.Error()})
		return
	}

	// Step 4: Check flow status
	if flow.Status != "VERIFIED" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Flow is not verified. Please initiate the process again."})
		return
	}

	// Step 5: Check password uniqueness
	if hash.CheckPasswordHash(newPassword, currentHashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "New password cannot be the same as the current password"})
		return
	}

	// Step 6: Hash new password
	newHashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password", "error": err.Error()})
		return
	}

	// Step 7: Update password
	if err := db.UpdateUserPassword(ctx, email, string(newHashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating password", "error": err.Error()})
		return
	}

	// Step 8: Insert audit and mark flow completed
	pwdResetTtl := time.Duration(pwdResetTtlInMinutes) * (time.Minute)

	flow.Status = "COMPLETED"

	if err := cache.SaveFlow(ctx, flow, pwdResetTtl); err != nil {
		log.Printf("redis save flow error: %v", err.Error())
	}

	_ = db.InsertPasswordResetAudit(ctx, flowId, email, "COMPLETED", ip, userAgent, "", flow.Attempts)

	// Step 9: Delete active and flow redis key
	if err := cache.DeleteActiveFlow(ctx, email); err != nil {
		log.Printf("redis delete active flow error: %v", err.Error())
	}

	if err := cache.DeleteFlow(ctx, flowId); err != nil {
		log.Printf("redis delete flow error: %v", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
