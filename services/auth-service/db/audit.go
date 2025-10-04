package db

import (
	"context"
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
				attempt_count,
				failure_reason
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	verifiedAt := any(nil)
	if status == "VERIFIED" {
		verifiedAt = time.Now().UTC()
	}

	failureReason := any(nil)
	if failureReasonTemp != "" {
		failureReason = failureReasonTemp
	}

	_, err := DB.Exec(ctx, query, flowID, email, status, ip, userAgent, time.Now().UTC(), verifiedAt, attemptCount, failureReason)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			_ = pgErr
		}
		return err
	}
	return nil
}
