package db

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
)

func InsertPasswordResetAudit(ctx context.Context, flowID, email, event, ip, userAgent string) error {
	query := `INSERT INTO password_reset_audit (flow_id, email, event, ip, user_agent, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := DB.Exec(ctx, query, flowID, email, event, ip, userAgent, time.Now().UTC())
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			_ = pgErr
		}
		return err
	}
	return nil
}
