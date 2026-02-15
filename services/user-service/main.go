package main

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/routes"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/utils"

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

	// Setup routes
	routes.SetupRoutes(r, s3Uploader)

	// Start the server
	if err := r.Run(config.AppConfig.Port); err != nil {
		panic(err)
	}
}
