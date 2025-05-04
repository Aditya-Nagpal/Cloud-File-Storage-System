package main

import (
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/api-gateway/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/api-gateway/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.AppConfig.FrontendBaseUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	routes.SetupRoutes(r)

	// Start gateway
	if err := r.Run(config.AppConfig.Port); err != nil {
		panic(err)
	}

}
