package main

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/routes"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/services/tasks"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	db.ConnectDatabase()

	r := gin.Default()
	s3Uploader, err := utils.NewS3Uploader()
	if err != nil {
		panic(err)
	}

	taskService := tasks.NewTaskService("localhost" + config.AppConfig.RedisPort)

	// Setup routes
	routes.SetupRoutes(r, s3Uploader, taskService)

	// Start the server
	if err := r.Run(config.AppConfig.Port); err != nil {
		panic(err)
	}
}
