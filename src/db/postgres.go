package db

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupConnection() *pgxpool.Pool{
	URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), URL)
	if err != nil {
		log.Fatalf("Failed to connect to database", err)
	}
	return conn
}