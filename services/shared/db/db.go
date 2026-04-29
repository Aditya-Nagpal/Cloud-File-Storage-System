package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	dbURL := "postgres://adityanagpal:bhaibhai10@localhost:5432/FastFiles"
	var err error

	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	log.Println("Connected to database successfully")
}
