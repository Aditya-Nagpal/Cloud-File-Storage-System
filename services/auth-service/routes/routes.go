package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/handlers"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/refresh", handlers.RefreshToken)
	r.GET("/logout", handlers.Logout)

	r.GET("/protected", middleware.JWTMiddleware(), func(c *gin.Context) {
		email, _ := c.Get("email")
		c.JSON(200, gin.H{
			"message": "Authenticated",
			"user":    email,
		})
	})

	forgotGroup := r.Group("/forgot-password")
	SetupForgotPasswordRoutes(forgotGroup)
}
