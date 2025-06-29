package models

type User struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Age   int    `json:"age" db:"age"`
}

type UpdateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
