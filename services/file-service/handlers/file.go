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

func ListFilesByPrefix() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := c.GetHeader("X-User-Email")
		if userEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
			return
		}
		parentPath := c.Query("parentPath")
		parentPath = userEmail + "/" + parentPath
		files, err := db.GetFilesByPrefix(c.Request.Context(), userEmail, parentPath)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve files from database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": files})
	}
}

func Upload(uploader *utils.S3Uploader) gin.HandlerFunc {
	return func(c *gin.Context) {
		uploadType := c.PostForm("uploadType")

		userEmail := c.GetHeader("X-User-Email")
		if userEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
			return
		}
		log.Println("userEmail: ", userEmail)
		log.Println("uploadType: ", uploadType)
		switch uploadType {
		case "file":
			UploadFile(c, uploader, userEmail)
		case "folder":
			uploadFolder(c, uploader, userEmail)
			c.JSON(http.StatusOK, gin.H{"response": "res"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid uploadType, must be 'file' or 'folder'"})
		}
	}
}

func UploadFile(c *gin.Context, uploader *utils.S3Uploader, userEmail string) {
	// Get uploaded file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	parentPath := c.PostForm("folderKey")
	parentPath = userEmail + "/" + parentPath
	key := parentPath + fileHeader.Filename
	log.Println((key))
	log.Println((parentPath))

	err = uploader.UploadFile(file, fileHeader, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file to S3"})
		return
	}

	metadata := models.FileMetaData{
		UserEmail:   userEmail,
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        fileHeader.Size,
		ParentPath:  parentPath,
		S3URL:       uploader.GetS3URL(key),
		UploadedAt:  time.Now(),
		Type:        "file",
	}

	err = db.InsertFileMetadata(c.Request.Context(), &metadata)
	if err != nil {
		fmt.Println("Insert to db error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store metadata to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "metadata": "metadata"})
}

func uploadFolder(c *gin.Context, uploader *utils.S3Uploader, userEmail string) {
	parentPath := c.PostForm("folderKey")
	folderName := c.PostForm("folderName")
	parentPath = userEmail + "/" + parentPath
	key := parentPath + folderName + "/"
	log.Println((key))

	err := uploader.UploadFolder(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create folder"})
		return
	}

	metadata := models.FileMetaData{
		UserEmail:   userEmail,
		FileName:    folderName,
		ContentType: "application/x-directory",
		Size:        0,
		ParentPath:  parentPath,
		S3URL:       uploader.GetS3URL(key),
		UploadedAt:  time.Now(),
		Type:        "folder",
	}

	if err := db.InsertFileMetadata(c.Request.Context(), &metadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder created", "metadata": metadata})
}

func DownloadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "in download api"})
}
