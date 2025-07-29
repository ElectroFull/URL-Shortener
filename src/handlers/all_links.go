package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


func All_Links(db *pgxpool.Pool) fiber.Handler {
	return func (c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token) // getting the token
		claims := user.Claims.(jwt.MapClaims) // turning it into the map
		username := claims["username"].(string)


		
		var user_id int
		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE username=$1", username).Scan(&user_id)

		if err != nil {
			log.Debug("Db query failed")
            return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
		}

		rows, err := db.Query(context.Background(), "SELECT original_url, short_url FROM urls WHERE user_id=$1", user_id)
		if err != nil {
            log.Debug("Query failed")
            return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
        }
		defer rows.Close()

		var links = make(map[string]string)

		for rows.Next() {
			var original_url, short_url string
			err := rows.Scan(&original_url, &short_url)
			if err != nil {
				log.Debug("Scan failed")
                return c.Status(fiber.StatusInternalServerError).SendString("Try again later")
			}

			links[short_url] = original_url 
		}

		return c.JSON(links)
	}
}