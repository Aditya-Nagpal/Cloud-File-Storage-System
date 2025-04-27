package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/handlers"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	r.GET("/auth/protected", middleware.JWTMiddleware(), func(c *gin.Context) {
		email, _ := c.Get("email")
		c.JSON(200, gin.H{
			"message": "Authenticated",
			"user":    email,
		})
	})
}
