package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/models"
	"github.com/jackc/pgx/v5"
)

func GetProfleById(ctx context.Context, userId int64) (*models.User, error) {
	query := `SELECT
				name,
				email,
				alternate_email,
				contact_number,
				gender,
				dob,
				country,
				timezone,
				about,
				plan,
				created_at,
				display_picture
			FROM users WHERE id = $1`

	var user models.User

	err := DB.QueryRow(ctx, query, userId).Scan(&user.Name, &user.Email, &user.AlternateEmail, &user.ContactNumber, &user.Gender, &user.DOB, &user.Country, &user.Timezone, &user.About, &user.Plan, &user.CreatedAt, &user.DisplayPicture)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateProfileDetails(ctx context.Context, userId int64, user *models.UpdateUser) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal update user data: %w", err)
	}

	var dataMap map[string]any
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return fmt.Errorf("failed to unmarshal update user data: %w", err)
	}

	delete(dataMap, "email")

	fieldsToUpdate := make(map[string]any)
	for key, value := range dataMap {
		if value != nil && value != 0 && value != "" {
			fieldsToUpdate[key] = value
		}
	}

	if len(fieldsToUpdate) == 0 {
		return fmt.Errorf("no fields to update")
	}

	var setClauses []string
	var args []any
	argIndex := 1

	for column, value := range fieldsToUpdate {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	args = append(args, userId)

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	_, err = DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile details: %w", err)
	}

	return nil
}

func UpdateDisplayPicture(ctx context.Context, userId int64, displayPictureURL string) error {
	query := `UPDATE users SET display_picture = $1 WHERE id = $2`
	_, err := DB.Exec(ctx, query, displayPictureURL, userId)
	if err != nil {
		return fmt.Errorf("failed to update display picture: %w", err)
	}
	return nil
}

func DeleteDisplayPicture(ctx context.Context, userId int64) error {
	_, err := DB.Exec(ctx, `UPDATE users SET display_picture = NULL WHERE id = $1`, userId)
	if err != nil {
		return fmt.Errorf("failed to remove display picture: %w", err)
	}
	return nil
}
