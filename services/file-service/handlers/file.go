package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/utils"
	"github.com/gin-gonic/gin"
)

func Upload(uploader *utils.S3Uploader) gin.HandlerFunc {
	return func(c *gin.Context) {
		uploadType := c.PostForm("uploadType")

		userEmail := c.GetHeader("X-User-Email")
		log.Println(uploadType)
		switch uploadType {
		case "file":
			UploadFile(c, uploader, userEmail)
		case "folder":
			CreateFolder(c, uploader, userEmail)
			c.JSON(http.StatusOK, gin.H{"response": "res"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid uploadType, must be 'file' or 'folder'"})
		}
	}
}

func UploadFile(c *gin.Context, uploader *utils.S3Uploader, userEmail string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user email from JWT (set by api-gateway middleware)
		userEmail := c.GetHeader("X-User-Email")

		// Get uploaded file
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get uploaded file"})
			return
		}
		defer file.Close()

		// S3 key
		s3Key := time.Now().Format("20060102_150405") + "_" + fileHeader.Filename

		// Upload file to S3
		s3URL, err := uploader.UploadFile(file, fileHeader, s3Key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file to S3"})
			return
		}

		// Save metadata to database
		metadata := models.FileMetaData{
			UserEmail:   userEmail,
			Filename:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
			S3Key:       s3Key,
			S3URL:       s3URL,
			UploadedAt:  time.Now(),
		}

		err = db.InsertFileMetadata(c.Request.Context(), &metadata)
		if err != nil {
			fmt.Println("Insert to db error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store metadata to database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "metadata": "metadata", "a": fileHeader, "b": userEmail})
	}
}

func CreateFolder(c *gin.Context, uploader *utils.S3Uploader, userEmail string) {
	folderKey := c.PostForm("folderKey")
	if folderKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing folderKey"})
		return
	}

	if folderKey[len(folderKey)-1] != '/' {
		folderKey += "/"
	}

	err := uploader.CreateFolder(folderKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create folder"})
		return
	}

	metadata := models.FileMetaData{
		UserEmail:   userEmail,
		Filename:    folderKey,
		ContentType: "application/x-directory",
		Size:        0,
		S3Key:       folderKey,
		S3URL:       uploader.GetS3URL(folderKey),
		UploadedAt:  time.Now(),
	}

	// if err := db.InsertFileMetadata(c.Request.Context(), &metadata); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Folder created", "metadata": metadata})
}

func DownloadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "in download api"})
}
