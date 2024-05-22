package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	// app.Static("/", "./web/templates")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World")
	})

	app.Get("/gay", func(c *fiber.Ctx) error {
		return c.SendString("you gay")
	})

	app.Listen(":3000")
}
