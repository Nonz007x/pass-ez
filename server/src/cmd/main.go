package main

import (
	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.ConnectDb()

	app := fiber.New()
	app.Use(cors.New())

	setupRoutes(app)

	app.Listen(":3000")
}
