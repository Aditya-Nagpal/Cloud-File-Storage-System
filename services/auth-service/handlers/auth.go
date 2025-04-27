package handlers

import (
	"net/http"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/models"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/hash"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/shared/jwt"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists
	if _, exists := models.Users[input.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hashedPassword, _ := hash.HashPassword(input.Password)
	input.Password = string(hashedPassword)
	models.Users[input.Email] = input

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	storedUser, exists := models.Users[input.Email]
	if !exists || !hash.CheckPasswordHash(input.Password, storedUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := jwt.Generate(storedUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
