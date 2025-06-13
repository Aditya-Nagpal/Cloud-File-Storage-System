package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/db"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/hash"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/jwt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 500, "error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if email already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err := db.DB.QueryRow(ctx, checkQuery, req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check email", "error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password", "error": err.Error()})
		return
	}

	// Insert user in db
	insertQuery := `INSERT INTO users (name, email, age, password, unhashed_password) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.DB.Exec(ctx, insertQuery, req.Name, req.Email, req.Age, hashedPassword, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": req, "message": "User registered successfully"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	var name string
	var age int
	var hashedPassword string

	query := `SELECT id, name, age, password FROM users WHERE email=$1`
	err := db.DB.QueryRow(ctx, query, req.Email).Scan(&id, &name, &age, &hashedPassword)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Error in checking email", "error": err.Error()})
		return
	}

	if !hash.CheckPasswordHash(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Generate access token short-lived
	accessToken, err := jwt.GenerateWithExpiry(req.Email, config.AppConfig.JwtSecret, 2*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not access generate token", "error": err.Error()})
		return
	}

	// Generate refresh token long-lived
	refreshToken, err := jwt.GenerateWithExpiry(req.Email, config.AppConfig.JwtSecret, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not refresh generate token", "error": err.Error()})
		return
	}

	// Set refresh token in secure HttpOnly cookie
	c.SetCookie("refreshToken", refreshToken, 7*24*60*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user": gin.H{
			"id":    id,
			"name":  name,
			"email": req.Email,
			"age":   age,
		},
	})
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No refresh token found"})
		return
	}

	claims, err := jwt.Verify(refreshToken, config.AppConfig.JwtSecret)
	if err != nil || claims == nil || claims.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired refresh token"})
		return
	}

	newAccessToken, err := jwt.GenerateWithExpiry(claims.Email, config.AppConfig.JwtSecret, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})
}

func Logout(c *gin.Context) {
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
