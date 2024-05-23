package main

import (
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	setupRoutesV1(v1)
}

func setupRoutesV1(r fiber.Router) {
	r.Get("/home", func(c *fiber.Ctx) error {
		return c.SendString("Home")
	})
}
