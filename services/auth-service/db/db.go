package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	dbURL := os.Getenv("DATABASE_URL")
	var err error

	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	fmt.Println("Connected to database successfully")
}
