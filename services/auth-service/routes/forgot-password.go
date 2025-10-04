package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/handlers"
	"github.com/gin-gonic/gin"
)

func SetupForgotPasswordRoutes(rg *gin.RouterGroup) {
	rg.POST("/start", handlers.StartPasswordReset)
}
