package routes

import (
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/upload", handlers.UploadFile)
	r.GET("/download/:filename", handlers.DownloadFile)
}
