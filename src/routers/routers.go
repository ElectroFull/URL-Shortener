package routers

import (
    "github.com/gofiber/fiber/v2"
	"github.com/electrofull/URL-Shortener/src/handlers"
)

func RegisterHandlers(app *fiber.App) {
	app.Get("/", handlers.Welcome)
}