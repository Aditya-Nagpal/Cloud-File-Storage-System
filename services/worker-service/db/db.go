package db

import (
	"context"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/db"
	"github.com/pgvector/pgvector-go"
)

func UpdateEmbedding(ctx context.Context, internalID int64, vector pgvector.Vector) error {
	query := "UPDATE entries SET embedding = $1 WHERE id = $2"
	_, err := db.DB.Exec(ctx, query, vector, internalID)
	if err != nil {
		return err
	}
	return nil
}
