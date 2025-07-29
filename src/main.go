package main

import (
	"github.com/electrofull/URL-Shortener/src/db"
	"github.com/electrofull/URL-Shortener/src/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"os"
)

func loadEnvironment() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file!")
	}
}

func getPort() string {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable is not set!")
	}
	return portString
}

func setupMiddleware(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format: "\033[32mINFO:\033[0m 	[${ip}]:${port}   ${status} - ${method} ${path}\n",
	})) // Just logs the basic info (requests)

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	})) // Catches an error and shows where the error happened.

}

func main() {
	loadEnvironment()

	portString := getPort()

	conn := db.SetupConnection()
	defer conn.Close()

	app := fiber.New()

	setupMiddleware(app)

	routers.RegisterHandlers(app, conn)

	log.Fatal(app.Listen(":" + portString))
}
