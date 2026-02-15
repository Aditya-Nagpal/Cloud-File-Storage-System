package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgconn"
)

func InsertPasswordResetAudit(ctx context.Context, flowID, email, status, ip, userAgent, failureReasonTemp string, attemptCount int) error {
	query := `INSERT INTO password_reset_audit (
		flow_id,
		email,
		status,
		ip_address,
		user_agent,
		created_at,
		verified_at,
		completed_at,
		attempt_count,
		failure_reason
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	verifiedAt := any(nil)
	completedAt := any(nil)
	switch status {
	case "VERIFIED":
		verifiedAt = time.Now().UTC()
	case "COMPLETED":
		completedAt = time.Now().UTC()
	}

	failureReason := any(nil)
	if failureReasonTemp != "" {
		failureReason = failureReasonTemp
	}

	_, err := DB.Exec(ctx, query, flowID, email, status, ip, userAgent, time.Now().UTC(), verifiedAt, completedAt, attemptCount, failureReason)
	if err != nil {
		log.Printf("Error inserting password reset audit: %v, userAgent: %s", err, userAgent)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			_ = pgErr
		}
		return err
	}
	return nil
}
