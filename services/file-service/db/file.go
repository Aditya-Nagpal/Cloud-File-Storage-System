package db

import (
	"context"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/models"
)

func InsertFileMetadata(ctx context.Context, meta *models.FileMetaData) error {
	query := `
		INSERT INTO file_metadata (user_email, filename, content_type, size, parent_path, s3_url, uploaded_at, type)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := DB.Exec(ctx, query, meta.UserEmail, meta.FileName, meta.ContentType, meta.Size, meta.ParentPath, meta.S3URL, meta.UploadedAt, meta.Type)
	return err
}

type ListFileResponse struct {
	Id          int       `json:"id"`
	UserEmail   string    `json:"user_email"`
	FileName    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Type        string    `json:"type"`
}

func GetFilesByPrefix(ctx context.Context, userEmail string, prefix string) ([]ListFileResponse, error) {
	query := `
		SELECT id, user_email, filename, content_type, size, uploaded_at, type FROM file_metadata
		WHERE user_email = $1 AND parent_path = $2
		ORDER BY uploaded_at DESC
	`
	rows, err := DB.Query(ctx, query, userEmail, prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []ListFileResponse
	for rows.Next() {
		var file ListFileResponse
		err := rows.Scan(&file.Id, &file.UserEmail, &file.FileName, &file.ContentType, &file.Size, &file.UploadedAt, &file.Type)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func DeleteFileMetadata(ctx context.Context, userEmail string, parentPath string, fileName string, isFolder bool) error {
	key := parentPath + fileName
	if isFolder {
		key += "/"
		query := `DELETE FROM file_metadata WHERE user_email = $1 AND (parent_path = $2 OR parent_path LIKE $3)`
		likePattern := key + "%"
		_, err := DB.Exec(ctx, query, userEmail, key, likePattern)
		if err != nil {
			return err
		}

		delSelf := `DELETE FROM file_metadata WHERE user_email = $1 AND parent_path = $2 AND filename = $3`
		_, err = DB.Exec(ctx, delSelf, userEmail, parentPath, fileName)
		return err
	}
	query := `DELETE FROM file_metadata WHERE user_email = $1 AND parent_path = $2 AND filename = $3`
	_, err := DB.Exec(ctx, query, userEmail, parentPath, fileName)
	return err
}

type FileRecord struct {
	UserEmail  string `json:"user_email"`
	FileName   string `json:"filename"`
	ParentPath string `json:"parent_path"`
	Type       string `json:"type"`
}

func GetFileRecordByID(ctx context.Context, id int) (*FileRecord, error) {
	query := `SELECT user_email, filename, parent_path, type FROM file_metadata WHERE id = $1`
	row := DB.QueryRow(ctx, query, id)

	var record FileRecord
	err := row.Scan(&record.UserEmail, &record.FileName, &record.ParentPath, &record.Type)
	if err != nil {
		return nil, err
	}

	return &record, nil
}
