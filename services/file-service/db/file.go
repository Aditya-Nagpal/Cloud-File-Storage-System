package db

import (
	"context"
	"fmt"
	"strings"

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

func GetEntityType(ctx context.Context, publicId string, userId int64) (string, error) {
	query := `SELECT type FROM entries WHERE public_id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var entityType string
	err := DB.QueryRow(ctx, query, publicId, userId).Scan(&entityType)
	if err == pgx.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return strings.ToLower(entityType), nil
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

func DeleteFile(ctx context.Context, publicId string, userId int64) error {
	query := `
		UPDATE entries
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE public_id = $1 AND user_id = $2 AND type = 'FILE' AND deleted_at IS NULL
	`

	err := DB.QueryRow(ctx, query, publicId, userId).Scan()
	if err == pgx.ErrNoRows {
		return fmt.Errorf("file not found or already deleted")
	} else if err != nil {
		return err
	}
	return nil
}

func DeleteFolder(ctx context.Context, publicId string, userId int64) error {
	query := `
		WITH RECURSIVE descendants AS (
			SELECT id FROM entries
			WHERE public_id = $1 AND user_id = $2 AND type = 'FOLDER' AND deleted_at IS NULL
			UNION ALL
			SELECT e.id FROM entries e
			INNER JOIN descendants d ON e.parent_id = d.id
			WHERE e.deleted_at IS NULL
		)
		UPDATE entries
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id IN (SELECT id FROM descendants)
	`

	err := DB.QueryRow(ctx, query, publicId, userId).Scan()
	if err == pgx.ErrNoRows {
		return fmt.Errorf("file not found or already deleted")
	} else if err != nil {
		return err
	}
	return nil
}

func GetDeleteFile(ctx context.Context, publicId string, userId int64) (*models.DeleteFile, error) {
	query := `SELECT name, type, s3_key FROM entries WHERE public_id = $1 AND user_id = $2`

	var file models.DeleteFile
	err := DB.QueryRow(ctx, query, publicId, userId).Scan(&file.Name, &file.Type, &file.S3Key)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &file, nil
}
