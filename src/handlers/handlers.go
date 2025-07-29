package handlers

import (
	"context"
	"github.com/electrofull/URL-Shortener/src/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"os"
	"strings"
)

func Welcome(c *fiber.Ctx) error {
	log.Info("Welcome Page have been visited...")
	return c.SendString("Welcome to URL-Shortener!")
}

type URL struct {
	Url string `json:"url"`
}

func isValidURL(rawURL string) bool {
	parsed, err := url.ParseRequestURI(rawURL)

	if err != nil {
		return false
	}

	domain_parts := strings.Split(parsed.Hostname(), ".") // ["www" ... "google", "com"]

	if len(domain_parts) < 2 || len(domain_parts[len(domain_parts)-1]) < 2 {
		return false
	}

	return parsed.Scheme != "" && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func Shorten(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !strings.HasPrefix(c.Get("Content-Type"), "application/json") {
			return c.Status(fiber.StatusBadRequest).SendString("Content-Type must be application/json")
		}

		Url := new(URL)
		if err := c.BodyParser(Url); err != nil {
			return err
		}

		if !isValidURL(Url.Url) {
			return c.Status(fiber.StatusBadRequest).SendString("Wrong URL format")
		}

		log.Info(Url) // to check whether it is correct

		user := c.Locals("user").(*jwt.Token) // getting the token
		claims := user.Claims.(jwt.MapClaims) // turning it into the map
		username := claims["username"].(string)
		// getting user_id
		var user_id, url_id int
		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE username=$1", username).Scan(&user_id)

		if err != nil {
			log.Debug("Db query failed")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		tx, err := db.Begin(context.Background()) // Transaction opens
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}
		defer tx.Rollback(context.Background())

		err = tx.QueryRow(context.Background(), "INSERT INTO urls (original_url, user_id) VALUES ($1, $2) RETURNING id", Url.Url, user_id).Scan(&url_id)
		if err != nil {
			log.Debug("Could not insert into urls")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		encoding := helpers.Base62Encode(url_id)

		base_url := os.Getenv("BASE_URL")

		if base_url == "" {
			log.Debug("BASE_URL environment variable not set")
			return c.Status(fiber.StatusInternalServerError).SendString("Server configuration error")
		}

		short_url := base_url + "/" + encoding

		_, err = tx.Exec(context.Background(), "UPDATE urls SET short_url = $1 WHERE id = $2", short_url, url_id)

		if err != nil {
			log.Debug("Can not insert into urls")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		err = tx.Commit(context.Background())

		if err != nil {
			log.Debug("DB Commit failed")
			return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		return c.JSON(fiber.Map{"short_url": short_url, "original_url": Url.Url})
	}
}
