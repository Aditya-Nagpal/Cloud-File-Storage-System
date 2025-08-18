package models

type User struct {
	Name           string  `json:"name" db:"name"`
	Email          string  `json:"email" db:"email"`
	Age            int     `json:"age" db:"age"`
	DisplayPicture *string `json:"displayPicture" db:"display_picture"`
}

type UpdateUser struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Age   int    `json:"age" db:"age"`
}
