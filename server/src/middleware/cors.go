package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	})
}

func RequireJSONContentType(c *fiber.Ctx) error {
	if c.Method() == "POST" || c.Method() == "PUT" {
		if c.Get("Content-Type") != "application/json" {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
				"error": "Unsupported Media Type",
				"message": "Content-Type must be application/json",
			})
		}
	}
	return c.Next()
}