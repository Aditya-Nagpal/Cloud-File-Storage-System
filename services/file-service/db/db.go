package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	dbURL := config.AppConfig.DatabaseURL
	var err error

	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Connected to database successfully")
}

func InsertFileMetadata(ctx context.Context, meta *models.FileMetaData) error {
	query := `
		INSERT INTO file_metadata (user_email, filename, content_type, size, parent_path, s3_url, uploaded_at, type)
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := DB.Exec(ctx, query, meta.UserEmail, meta.FileName, meta.ContentType, meta.Size, meta.ParentPath, meta.S3URL, meta.UploadedAt, meta.Type)
	return err
}

type ListFileResponse struct {
	FileName    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	S3URL       string    `json:"s3_url"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Type        string    `json:"type"`
}

func GetFilesByPrefix(ctx context.Context, userEmail string, prefix string) ([]ListFileResponse, error) {
	query := `
		SELECT filename, content_type, size, s3_url, uploaded_at, type FROM file_metadata
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
		err := rows.Scan(&file.FileName, &file.ContentType, &file.Size, &file.S3URL, &file.UploadedAt, &file.Type)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func DeleteFileMetadata(ctx context.Context, userEmail string, parentPath string, fileName string, isFolder bool) error {
	key := userEmail + "/" + parentPath + fileName
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
