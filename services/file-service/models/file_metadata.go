package models

import "time"

type FileMetaData struct {
	Id          int       `db:"id"`
	UserEmail   string    `db:"user_email"`
	FileName    string    `db:"filename"`
	ContentType string    `db:"content_type"`
	Size        int64     `db:"size"`
	ParentPath  string    `db:"parent_path"`
	S3URL       string    `db:"s3_url"`
	UploadedAt  time.Time `db:"uploaded_at"`
	Type        string    `db:"type"`
}
