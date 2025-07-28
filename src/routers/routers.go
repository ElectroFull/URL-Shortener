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
	app.Get("/", handlers.Welcome)
	app.Post("/register", handlers.Register(conn, jwt_secret))
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwt_secret)},
	})) // JWT Authentication
	app.Post("/shorten", handlers.Shorten(conn))
	
}