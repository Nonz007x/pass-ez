package main

import (
	"github.com/Nonz007x/pass-ez/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	setupRoutesV1(v1)
}

func setupRoutesV1(r fiber.Router) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"hello":  "greeting",
		})
	})

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.GetSalt)
}
