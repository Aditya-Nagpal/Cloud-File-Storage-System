package models

import "time"

type FileMetaData struct {
	Id          int       `db:"id"`
	UserEmail   string    `db:"user_email"`
	Filename    string    `db:"filename"`
	ContentType string    `db:"content_type"`
	Size        int64     `db:"size"`
	S3Key       string    `db:"s3_key"`
	S3URL       string    `db:"s3_url"`
	UploadedAt  time.Time `db:"uploaded_at"`
}
