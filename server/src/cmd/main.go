package main

import (
	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectDb()

	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
