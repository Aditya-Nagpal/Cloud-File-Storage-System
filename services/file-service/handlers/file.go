package handlers

import (
	"log"
	"net/http"
	"strconv"
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve files from database", "error": err.Error()})
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
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get uploaded file", "error": err.Error()})
		return
	}
	defer file.Close()

	parentPath := c.PostForm("folderKey")
	parentPath = utils.GetParentPath(userEmail, parentPath)
	key := parentPath + fileHeader.Filename

	err = uploader.UploadFile(file, fileHeader, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file to S3", "error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store metadata to database", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "metadata": "metadata"})
}

func uploadFolder(c *gin.Context, uploader *utils.S3Uploader, userEmail string) {
	parentPath := c.PostForm("folderKey")
	folderName := c.PostForm("folderName")
	parentPath = utils.GetParentPath(userEmail, parentPath)
	key := parentPath + folderName + "/"

	err := uploader.UploadFolder(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create folder", "error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store metadata to database", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder created", "metadata": metadata})
}

type DeleteRequest struct {
	ParentPath string `json:"parentPath"`
	FileName   string `json:"fileName"`
	Type       string `json:"type"`
}

func DeleteContent(uploader *utils.S3Uploader) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.FileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
			return
		}

		userEmail := c.GetHeader("X-User-Email")
		if userEmail == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing user email"})
			return
		}

		parentPath := req.ParentPath
		parentPath = utils.GetParentPath(userEmail, parentPath)
		key := parentPath + req.FileName
		isFolder := false
		if req.Type == "folder" {
			key += "/"
			isFolder = true
		}

		if err := uploader.DeleteObject(key, isFolder); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete object from S3", "error": err.Error()})
			return
		}

		if err := db.DeleteFileMetadata(c.Request.Context(), userEmail, parentPath, req.FileName, isFolder); err != nil {
			log.Printf("Failed to delete metadata for %s: %v", key, err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete metadata from database", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted successfully", "a": req, "key": key, "isFolder": isFolder})
	}
}

func DownloadFile(uploader *utils.S3Uploader) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := c.GetHeader("X-User-Email")
		if userEmail == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing user email"})
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file ID"})
			return
		}

		fileRecode, err := db.GetFileRecordByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "File not found", "error": err.Error()})
			return
		}

		if fileRecode.UserEmail != userEmail {
			c.JSON(http.StatusForbidden, gin.H{"message": "Access denied"})
			return
		}

		if fileRecode.Type == "folder" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot download a folder"})
			return
		}

		parentPath := fileRecode.ParentPath
		key := parentPath + fileRecode.FileName

		url, err := uploader.GeneratePresignedURL(key, 2*time.Minute, fileRecode.FileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate download URL", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"downloadURL": url})
	}
}
