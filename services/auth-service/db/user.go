package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func RegisterUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (
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
		password,
		terms_and_privacy,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

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
		fmt.Print("error: ", err, err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			fmt.Println("EMAIL_ALREADY_EXISTS", pgErr.Code)
			return fmt.Errorf("EMAIL_ALREADY_EXISTS")
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
