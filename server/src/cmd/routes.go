package main

import (
	"github.com/Nonz007x/pass-ez/src/handler"
	"github.com/Nonz007x/pass-ez/src/middleware"
	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
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

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	restricted := r.Group("")
	restricted.Use(middleware.AuthRequired())

	setupRestrictedRoutesV1(restricted)
}

func setupRestrictedRoutesV1(r fiber.Router) {
	r.Get("/restricted", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		return c.SendString("Welcome " + name)
	})

	r.Get("/test", Test)
}

func Test(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	exp := claims["exp"].(string)
	return c.SendString("Welcome " + sub + " Exp: " + exp)
}
