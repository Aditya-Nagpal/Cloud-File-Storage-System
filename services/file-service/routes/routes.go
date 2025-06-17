package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/handlers"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/utils"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, s3Uploader *utils.S3Uploader) {
	// Initialize the S3Uploader
	r.GET("/list", handlers.ListFilesByPrefix())
	r.POST("/upload", handlers.Upload(s3Uploader))
	r.GET("/download/:filename", handlers.DownloadFile)
	r.DELETE("/delete", handlers.DeleteContent(s3Uploader))
}
