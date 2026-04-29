package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/utils"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/db"
)

func HandleGenerateEmbedding(ctx context.Context, t *asynq.Task) error {
	var payload models.Payload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	s3Key := payload.S3Key

	content, err := utils.FetchContentFromS3(ctx, s3Key)
	if err != nil {
		return err
	}

	vector, err := GenerateEmbedding(content, config.AppConfig.OpenAiApiKey)
	if err != nil {
		return err
	}

	err = db.UpdateEmbedding(ctx, payload.InternalID, vector)
	if err != nil {
		return err
	}
	return nil
}
