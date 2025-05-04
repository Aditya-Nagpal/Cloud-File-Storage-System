package middleware

import (
	"net/http"
	"strings"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/jwt"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing or invalid"})
			return
		}

		accessTokenStrings := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.Verify(accessTokenStrings, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired access token"})
			return
		}

		// Add user info to headers
		c.Request.Header.Set("X-User-Email", claims.Email)
	}
}
