package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/jwt"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := jwt.Verify(tokenString, config.AppConfig.JwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user still exists in DB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
		err = db.DB.QueryRow(ctx, query, claims.Email).Scan(&exists)
		if err != nil || !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User no longer exists"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)

		c.Next()
	}
}
