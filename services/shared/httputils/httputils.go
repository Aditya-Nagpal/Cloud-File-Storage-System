package httputils

import (
	"fmt"
	"net/http"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/utils"
	"github.com/gin-gonic/gin"
)

func GetUserIdHeader(c *gin.Context) (int64, error) {
	userIdStr := c.GetHeader("X-User-Id")
	if userIdStr == "" {
		return 0, fmt.Errorf("HEADER_EMPTY")
	}

	userId, err := utils.StringToInt64(userIdStr)
	if err != nil {
		return 0, fmt.Errorf("HEADER_INVALID_USER_ID")
	}

	return userId, nil
}

func HandleUserIdHeaderError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch err.Error() {
	case "HEADER_EMPTY":
		c.JSON(http.StatusUnauthorized, gin.H{"message": "X-User-Email header is missing"})
	case "HEADER_INVALID_USER_ID":
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unable to parse X-User-Email", "error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch X-User-Email header", "error": err.Error()})
	}

	return true
}
