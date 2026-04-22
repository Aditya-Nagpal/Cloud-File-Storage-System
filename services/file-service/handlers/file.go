package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/utils"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/httputils"
	"github.com/gin-gonic/gin"
)

func ListFilesByParentId() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := httputils.GetUserIdHeader(c)
		if httputils.HandleUserIdHeaderError(c, err) {
			return
		}

		publicParentID := c.Query("parent_id")

		var internalParentID *int64
		if publicParentID != "" {
			id, err := db.GetInternalID(c.Request.Context(), publicParentID, userId)
			if id == nil && err == nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "Parent directory not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get internal ID", "error": err.Error()})
				return
			}

			internalParentID = id
		}

		files, err := db.GetFilesByParentId(c.Request.Context(), userId, internalParentID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve files from database", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}

type UploadRequest struct {
	Name           string `json:"name"`
	PublicParentID string `json:"parentId"`
	EntityType     string `json:"entityType" binding:"required"`
}

func Upload(uploader *utils.S3Uploader) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
			return
		}

		userId, err := httputils.GetUserIdHeader(c)
		if httputils.HandleUserIdHeaderError(c, err) {
			return
		}

		entityType := req.EntityType
		publicParentID := req.PublicParentID

		var internalParentID *int64
		if publicParentID != "" {
			id, err := db.GetInternalID(c.Request.Context(), publicParentID, userId)
			if id == nil && err == nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "Parent directory not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get internal ID", "error": err.Error()})
				return
			}

			internalParentID = id
		}

		publicId, err := utils.GenerateUniqueID(12)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate unique ID", "error": err.Error()})
			return
		}

		switch entityType {
		case "file":
			UploadFile(c, uploader, userId, publicId, internalParentID)
		case "folder":
			name := req.Name
			if name == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Folder name is required"})
				return
			}
			UploadFolder(c, userId, publicId, name, internalParentID)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entityType, must be 'file' or 'folder'"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "res"})
	}
}

func UploadFile(c *gin.Context, uploader *utils.S3Uploader, userId int64, publicId string, internalParentID *int64) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get uploaded file", "error": err.Error()})
		return
	}
	defer file.Close()

	baseName, ext := utils.ParseFilename(header.Filename)
	size := header.Size
	s3Key := utils.CreateS3Key(userId, publicId, ext)

	err = uploader.UploadFile(file, header, s3Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file to S3", "error": err.Error()})
		return
	}

	parentId := sql.NullInt64{Valid: false}
	if internalParentID != nil {
		parentId = sql.NullInt64{Int64: *internalParentID, Valid: true}
	}

	entryData := models.EntryData{
		PublicId:    publicId,
		UserId:      userId,
		ParentId:    parentId.Int64,
		Name:        baseName,
		Type:        "FILE",
		ContentType: header.Header.Get("Content-Type"),
		Extension:   ext,
		Size:        size,
		S3Key:       sql.NullString{String: s3Key, Valid: true},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = db.InsertEntryData(c.Request.Context(), &entryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store entry data to db", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func UploadFolder(c *gin.Context, userId int64, publicId string, name string, internalParentID *int64) {
	parentId := sql.NullInt64{Valid: false}
	if internalParentID != nil {
		parentId = sql.NullInt64{Int64: *internalParentID, Valid: true}
	}

	entryData := models.EntryData{
		PublicId:    publicId,
		UserId:      userId,
		ParentId:    parentId.Int64,
		Name:        name,
		Type:        "FOLDER",
		ContentType: "application/x-directory",
		Extension:   "",
		Size:        0,
		S3Key:       sql.NullString{Valid: false},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.InsertEntryData(c.Request.Context(), &entryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to store entry data to db", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder uploaded successfully"})
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
