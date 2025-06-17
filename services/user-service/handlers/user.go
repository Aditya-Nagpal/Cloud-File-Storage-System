package handlers

import (
	"net/http"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/db"

	"github.com/gin-gonic/gin"
)

func GetProfileDetails(c *gin.Context) {
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Email header is missing"})
		return
	}

	user, err := db.GetUserByEmail(c.Request.Context(), userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": user})
}
