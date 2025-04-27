package main

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	if err := r.Run(":8001"); err != nil {
		panic(err)
	}
}
