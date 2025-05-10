package main

import (
	"github.com/sai-zack-dev/FlatSync-API/database"
	"github.com/sai-zack-dev/FlatSync-API/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	// Auth routes
	app.Post("/api/register", handlers.Register)
	app.Post("/api/login", handlers.Login)

	app.Get("/api/protected", handlers.Protected)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber!")
	})

	app.Listen(":3000")
}
