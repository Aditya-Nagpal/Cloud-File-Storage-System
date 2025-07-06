package handlers

import (
	"net/http"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/models"

	"github.com/gin-gonic/gin"
)

func GetProfileDetails(c *gin.Context) {
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
		return
	}

	user, err := db.GetProfleByEmail(c.Request.Context(), userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": user})
}

func UpdateProfileDetails(c *gin.Context) {
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
