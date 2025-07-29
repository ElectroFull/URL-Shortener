package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


func Register(db *pgxpool.Pool, jwt_secret string) fiber.Handler{
	return func(c *fiber.Ctx) error {
		username, password := c.FormValue("username"), c.FormValue("password")
		if username == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Username and password required")
		}
		hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil{
			log.Debug("Bcrypt function crashed")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		// DB logic
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Debug("Can not get a transaction")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		defer tx.Rollback(context.Background())
		_, err = tx.Exec(context.Background(), "INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hashed_pass))

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				return c.Status(fiber.StatusConflict).SendString("Username already exists")
			}
			log.Debug("Can not register a user into db")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		err = tx.Commit(context.Background())

		if err != nil {
			log.Debug("DB Commit failed")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		// JWT logic
		claims := jwt.MapClaims{
			"username": username,
			"exp":   time.Now().Add(time.Hour * 48).Unix(),  // 2 days
		}

		header_payload := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := header_payload.SignedString([]byte(jwt_secret))

		if err != nil {
			log.Debug("Could not sign JWT")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		return c.JSON(fiber.Map{"token": token})
	}
}

func Login(db *pgxpool.Pool, jwt_secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username, password := c.FormValue("username"), c.FormValue("password")
		if username == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Username and password required")
		}

		// DB LOGIC
		var hash_pass string
		err := db.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE username=$1", username).Scan(&hash_pass)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
			}
			log.Debug("Db query failed")
            return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		
		if bcrypt.CompareHashAndPassword([]byte(hash_pass), []byte(password)) != nil{
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
		}

		// JWT logic
		claims := jwt.MapClaims{
			"username": username,
			"exp":   time.Now().Add(time.Hour * 48).Unix(),  // 2 days
		}

		header_payload := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := header_payload.SignedString([]byte(jwt_secret))

		if err != nil {
			log.Debug("Could not sign JWT")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		return c.JSON(fiber.Map{"token": token})
	}
}