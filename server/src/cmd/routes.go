package main

import (
	"github.com/Nonz007x/pass-ez/src/handler"
	"github.com/Nonz007x/pass-ez/src/middleware"
	"github.com/gofiber/fiber/v2"

	"time"

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

	r.Get("/test", Test)
	r.Post("/ciphers", handler.CreateItem) // create an item
	// r.Put("/ciphers/:id")           // edit an item
	// r.Put("/ciphers/:id/delete") // put item into trash
	// r.Delete("/ciphers/:id") // permanently delete an item
	r.Get("/validate-token", middleware.ValidateToken)
	r.Get("/sync", handler.GetItems)
}

func Test(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	_, subOk := claims["sub"].(string)
	exp, expOk := claims["exp"].(float64)
	userId, userIdOk := claims["user_id"].(string)

	if !subOk || !expOk || !userIdOk {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	expStr := time.Unix(int64(exp), 0).Format(time.RFC3339)
	return c.SendString("Welcome " + userId + " Exp: " + expStr)
}
