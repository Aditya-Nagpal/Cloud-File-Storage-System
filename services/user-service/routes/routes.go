package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/profile", handlers.GetProfileDetails)
}
