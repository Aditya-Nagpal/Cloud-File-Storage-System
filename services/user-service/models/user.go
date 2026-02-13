package models

import (
	"time"
)

type User struct {
	Email          string    `json:"email" binding:"omitempty,email"`
	Name           string    `json:"name" binding:"required"`
	AlternateEmail string    `json:"alternate_email" binding:"omitempty,email"`
	ContactNumber  string    `json:"contact_number" binding:"required,min=10"`
	Gender         string    `json:"gender" binding:"required"`
	DOB            time.Time `json:"dob" binding:"required"`
	Age            int       `json:"age" binding:"min=0"`
	Country        string    `json:"country" binding:"required"`
	Timezone       string    `json:"timezone" binding:"required"`
	About          string    `json:"about" binding:"omitempty,max=500"`
	Plan           string    `json:"plan" binding:"required"`
	CreatedAt      time.Time `json:"created_at" binding:"required"`
	DisplayPicture *string   `json:"display_picture" binding:"omitempty"`
}

type UpdateUser struct {
	Email          string `json:"email" binding:"omitempty,email"`
	Name           string `json:"name" binding:"omitempty"`
	AlternateEmail string `json:"alternate_email" binding:"omitempty,email"`
	ContactNumber  string `json:"contact_number" binding:"omitempty,min=10"`
	Gender         string `json:"gender" binding:"omitempty"`
	Country        string `json:"country" binding:"omitempty"`
	Timezone       string `json:"timezone" binding:"omitempty"`
	About          string `json:"about" binding:"max=500,omitempty"`
	Plan           string `json:"plan" binding:"omitempty"`
}
