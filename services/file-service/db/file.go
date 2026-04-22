package db

import (
	"context"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/models"
	"github.com/jackc/pgx/v5"
)

func GetInternalID(ctx context.Context, publicId string, userId int64) (*int64, error) {
	query := `SELECT id FROM entries WHERE public_id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var internalId int64
	err := DB.QueryRow(ctx, query, publicId, userId).Scan(&internalId)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &internalId, nil
}

func GetFilesByParentId(ctx context.Context, userId int64, internalParentID *int64) ([]models.ListFileResponse, error) {
	query := `
		SELECT
			public_id,
			name,
			type,
			content_type,
			size,
			created_at,
			updated_at
		FROM entries
		WHERE user_id = $1
			AND (
				($2::BIGINT IS NULL AND parent_id IS NULL)
				OR
				(parent_id = $2)
			)
			AND deleted_at IS NULL
		ORDER BY
			type DESC,
			updated_at DESC
	`
	rows, err := DB.Query(ctx, query, userId, internalParentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]models.ListFileResponse, 0)
	for rows.Next() {
		var file models.ListFileResponse
		err := rows.Scan(&file.PublicId, &file.Name, &file.Type, &file.ContentType, &file.Size, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func InsertEntryData(ctx context.Context, data *models.EntryData) error {
	query := `
		INSERT INTO file_metadata (public_id, user_id, parent_id, name, type, content_type, extension, size, s3_key, created_at, updated_at)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := DB.Exec(ctx, query, data.PublicId, data.UserId, data.ParentId, data.Name, data.Type, data.ContentType, data.Extension, data.Size, data.S3Key, data.CreatedAt, data.UpdatedAt)
	return err
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
