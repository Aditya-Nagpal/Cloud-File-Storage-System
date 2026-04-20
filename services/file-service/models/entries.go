package models

import (
	"database/sql"
	"time"
)

type ListFileResponse struct {
	PublicId    int64     `json:"public_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FileMetaData struct {
	Id          int            `db:"id"`
	UserEmail   string         `db:"user_email"`
	FileName    string         `db:"filename"`
	ContentType string         `db:"content_type"`
	Size        int64          `db:"size"`
	ParentPath  string         `db:"parent_path"`
	S3URL       sql.NullString `db:"s3_url"`
	UploadedAt  time.Time      `db:"uploaded_at"`
	Type        string         `db:"type"`
}
