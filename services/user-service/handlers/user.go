package handlers

import (
	"fmt"
	"net/http"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/utils"

	"github.com/gin-gonic/gin"
)

func GetProfileDetails(c *gin.Context) {
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
		return
	}

	user, err := db.GetProfleByEmail(c.Request.Context(), userEmail)
	if user == nil && err == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found", "error": err.Error()})
		return
	}
	user.Email = userEmail
	user.Age = utils.CalculateAge(user.DOB)
	// dob := user.DOB
	// fmt.Println(dob)
	c.JSON(http.StatusOK, gin.H{"profile": user})
}

func UpdateProfileDetails(uploader *utils.S3Uploader) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := c.GetHeader("X-User-Email")
		if userEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
			return
		}

		removeDp := c.Query("removeDp") == "true"
		dpOnly := c.Query("dp") == "true"

		if removeDp {
			s3Key := "displayPictures/" + userEmail + "/"
			if err := uploader.DeleteDisplayPicture(s3Key); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete display picture", "error": err.Error()})
				return
			}

			if err := db.DeleteDisplayPicture(c.Request.Context(), userEmail); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove display picture from profile", "error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Display picture removed successfully"})
		} else if dpOnly {
			file, fileHeader, err := c.Request.FormFile("displayPicture")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Display picture file is required"})
				return
			}

			s3Key := "displayPictures/" + userEmail + "/" + fileHeader.Filename

			s3Url, err := uploader.UploadDisplayPicture(file, fileHeader, s3Key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload display picture", "error": err.Error()})
				return
			}

			if err := db.UpdateDisplayPicture(c.Request.Context(), userEmail, s3Url); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update display picture", "error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "displayPicture": s3Url})
		} else {
			var update *models.UpdateUser
			if err := c.ShouldBindJSON(&update); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
				return
			}
			update.Email = userEmail

			if err := db.UpdateProfileDetails(c.Request.Context(), update); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update profile", "error": err.Error()})
				return
			}
			fmt.Println("Updating profile details")
			c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "updatedUser": update})
		}
	}
}
