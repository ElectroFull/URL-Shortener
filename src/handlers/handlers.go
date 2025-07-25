package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/log"
)

func Welcome(c *fiber.Ctx) error {
	log.Info("Welcome Page have been visited...")
	return c.SendString("Welcome to URL-Shortener!")
}