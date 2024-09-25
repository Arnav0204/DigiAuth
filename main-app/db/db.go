package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func InitDB() error {
	// Load .env file
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}

	// Read environment variables
	user := os.Getenv("user")
	password := os.Getenv("password")
	host := os.Getenv("host")
	port := os.Getenv("port")
	dbname := os.Getenv("dbname")

	connString := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname

	// Connect to the database
	DB, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return err
	}

	return nil
}

func CloseDB() {
	if err := DB.Close(context.Background()); err != nil {
		log.Fatalf("Error closing database connection: %v", err)
	}
}
