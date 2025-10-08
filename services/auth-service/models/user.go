package models

type User struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	AlternateEmail  string `json:"alternateEmail" binding:"omitempty,email"`
	ContactNumber   string `json:"contactNumber" binding:"required,min=10"`
	Country         string `json:"country" binding:"required"`
	Gender          string `json:"gender" binding:"required"`
	DOB             string `json:"dob" binding:"required"`
	Age             int    `json:"age" binding:"required"`
	Timezone        string `json:"timezone" binding:"required"`
	Plan            string `json:"plan" binding:"required"`
	About           string `json:"about"`
	Password        string `json:"password" binding:"required"`
	TermsAndPrivacy bool   `json:"termsAndPrivacy"`
}

var Users = make(map[string]User)
