package main

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/tasks"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/worker-service/utils"
	"github.com/hibiken/asynq"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	db.ConnectDatabase()

	err := utils.NewS3Uploader()
	if err != nil {
		panic(err)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("task:generate_embedding", tasks.HandleGenerateEmbedding)
	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
