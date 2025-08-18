package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/config"
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
