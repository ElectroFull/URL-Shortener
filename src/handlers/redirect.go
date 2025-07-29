package handlers

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Redirect(db *pgxpool.Pool) fiber.Handler{
	return func (c *fiber.Ctx) error{
		shortcode := c.Params("shortcode", "")
		if shortcode == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var original_url string 

		base_url := os.Getenv("BASE_URL")
		
		if base_url == "" {
            log.Debug("BASE_URL environment variable not set")
            return c.Status(fiber.StatusInternalServerError).SendString("Server configuration error")
        }

		shortcode = base_url + "/" + shortcode

		err := db.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_url=$1", shortcode).Scan(&original_url)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.SendStatus(fiber.StatusNotFound)
			}
            log.Debug("Database query failed:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		
		// Planning to add reusability of the shortcodes in the future

		return c.Redirect(original_url, fiber.StatusTemporaryRedirect)
	}
}