package handlers

import (
	"net/http"
	// "os"
	// "path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "in upload api"})
}

func DownloadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "in download api"})
}
