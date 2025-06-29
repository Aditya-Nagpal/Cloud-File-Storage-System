package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	dbURL := config.AppConfig.DatabaseURL
	var err error

	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Connected to database successfully")
}

func GetProfleByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	query := `SELECT email, name, age FROM users WHERE email = $1`
	row := DB.QueryRow(ctx, query, userEmail)

	var user models.User
	err := row.Scan(&user.Email, &user.Name, &user.Age)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateProfileDetails(ctx context.Context, user *models.UpdateUser) error {
	query := `UPDATE users SET name = $1, age = $2 WHERE email = $3`

	_, err := DB.Exec(ctx, query, user.Name, user.Age, user.Email)
	return err
}
