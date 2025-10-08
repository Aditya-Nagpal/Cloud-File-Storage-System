package db

import (
	"context"
	"fmt"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/models"
	"github.com/jackc/pgx/v5"
)

func GetProfleByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT email, name, age, display_picture FROM users WHERE email = $1`

	var user models.User

	err := DB.QueryRow(ctx, query, email).Scan(&user.Email, &user.Name, &user.Age, &user.DisplayPicture)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateDisplayPicture(ctx context.Context, userEmail, displayPictureURL string) error {
	query := `UPDATE users SET display_picture = $1 WHERE email = $2`
	_, err := DB.Exec(ctx, query, displayPictureURL, userEmail)
	if err != nil {
		return fmt.Errorf("failed to update display picture: %w", err)
	}
	return nil
}

func UpdateProfileDetails(ctx context.Context, user *models.UpdateUser) error {
	fmt.Println("Updating profile for user:", user)
	query := `UPDATE users SET name = $1, age = $2 WHERE email = $3`
	_, err := DB.Exec(ctx, query, user.Name, user.Age, user.Email)
	if err != nil {
		return fmt.Errorf("failed to update profile details: %w", err)
	}
	return nil
}
