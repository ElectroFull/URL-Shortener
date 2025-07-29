package db

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func SetupConnection() *pgxpool.Pool {
	URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = runMigrations(conn)

	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	return conn
}

func runMigrations(db *pgxpool.Pool) error {
	sql_lines, err := os.ReadFile("../tables.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), string(sql_lines))

	if err != nil {
		return err
	}

	log.Info("Database schema initialized successfully")

	return nil
}
