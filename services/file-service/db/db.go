package db

import (
	"context"
	"fmt"
	"log"

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
	query := `INSERT INTO file_metadata (user_email, filename, content_type, size, s3_key, s3_url, uploaded_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := DB.Exec(ctx, query, meta.UserEmail, meta.Filename, meta.ContentType, meta.Size, meta.S3Key, meta.S3URL, meta.UploadedAt)
	return err
}
