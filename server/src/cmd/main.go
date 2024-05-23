package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Nonz007x/pass-ez/src/database"
)

func main() {

	database.ConnectDb()

	app := fiber.New()

	app.Get("api/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World")
	})

	app.Get("api/gay", func(c *fiber.Ctx) error {
		return c.SendString("you gay")
	})

	app.Listen(":3000")
}
