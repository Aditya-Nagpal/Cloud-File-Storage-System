package db

import (
	"context"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

func DoesEmailExist(ctx context.Context, email string) (bool, error) {
	var count int

	err := DB.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func RegisterUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (
		name,
		email,
		alternate_email,
		contact_number,
		gender,
		dob,
		age,
		country,
		timezone,
		about,
		plan,
		password,
		terms_and_privacy,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	dob, err := time.Parse("2006-01-02", user.DOB)
	if err != nil {
		return err
	}
	_, err = DB.Exec(ctx, query,
		user.Name,
		user.Email,
		user.AlternateEmail,
		user.ContactNumber,
		user.Gender,
		dob,
		user.Age,
		user.Country,
		user.Timezone,
		user.About,
		user.Plan,
		user.Password,
		user.TermsAndPrivacy,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			_ = pgErr
		}
		return err
	}
	return nil
}

func GetUserHashedPassword(ctx context.Context, email string) (string, error) {
	var hashedPassword string

	err := DB.QueryRow(ctx, `SELECT password FROM users WHERE email=$1`, email).Scan(&hashedPassword)
	if err == pgx.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func UpdateUserPassword(ctx context.Context, email, password string) error {
	_, err := DB.Exec(ctx, `UPDATE users SET password=$1 WHERE email=$2`, password, email)
	if err != nil {
		return err
	}

	return nil
}
