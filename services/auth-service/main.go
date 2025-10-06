package main

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/cache"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/routes"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/services/sqs"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	db.ConnectDatabase()
	cache.InitRedis()
	sqs.InitSQS()

	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	if err := r.Run(config.AppConfig.Port); err != nil {
		panic(err)
	}
}
