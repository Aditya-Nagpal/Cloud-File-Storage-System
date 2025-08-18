package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/handlers"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/utils"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, s3Uploader *utils.S3Uploader) {
	// Initialize the S3Uploader
	r.GET("/profile", handlers.GetProfileDetails)
	r.PATCH("/profile", handlers.UpdateProfileDetails(s3Uploader))
}
