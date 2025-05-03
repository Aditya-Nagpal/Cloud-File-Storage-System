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
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Error in checking email", "error": err.Error()})
		return
	}

	if !hash.CheckPasswordHash(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := jwt.Generate(req.Email, config.AppConfig.JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    id,
			"name":  name,
			"email": req.Email,
			"age":   age,
		},
	})
}
