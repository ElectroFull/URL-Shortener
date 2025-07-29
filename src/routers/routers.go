package routers

import (
    "github.com/gofiber/fiber/v2"
	"github.com/electrofull/URL-Shortener/src/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	jwtware "github.com/gofiber/contrib/jwt"
	"os"
)

func RegisterHandlers(app *fiber.App, conn *pgxpool.Pool) {
	jwt_secret := os.Getenv("JWT_SECRET")

	// Public routes
	app.Get("/", handlers.Welcome)
	app.Post("/register", handlers.Register(conn, jwt_secret))
	app.Post("/login", handlers.Login(conn, jwt_secret))

	// Protected with JWT
	app.Post("/shorten", jwtware.New(jwtware.Config{ // POST /shorten
		SigningKey: jwtware.SigningKey{Key: []byte(jwt_secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	}), handlers.Shorten(conn))
	app.Get("/all", jwtware.New(jwtware.Config{           // GET /all
		SigningKey: jwtware.SigningKey{Key: []byte(jwt_secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	}), handlers.All_Links(conn))

	// public
	app.Get("/:shortcode", handlers.Redirect(conn))
}