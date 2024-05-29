package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/Nonz007x/pass-ez/src/middleware"
	"github.com/Nonz007x/pass-ez/src/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse JSON",
		})
	}

	email := req.Email

	// hash with SHA-256 and encode to base-64 string
	bytePassword := sha256.Sum256([]byte(req.Password))
	hashedPassword := base64.StdEncoding.EncodeToString(bytePassword[:])

	var result models.SaltKey
	db := database.DB.Db
	query := `
		SELECT users.salt, vaults.key 
		FROM users 
		INNER JOIN vaults ON users.id = vaults.owner_id 
		WHERE email = ? AND master_password = ?;
	`
	err := db.Raw(query, email, hashedPassword).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				models.ErrorResponse{
					Error:            "not_found",
					ErrorDescription: "user not found",
					Message:          "Wrong email or password. Try again.",
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.DatabaseError)
	}

	if result.Key == "" || result.Salt == "" {
		return c.Status(fiber.StatusNotFound).JSON(
			models.ErrorResponse{
				Error:            "not_found",
				ErrorDescription: "user not found",
				Message:          "Wrong email or password. Try again.",
			},
		)
	}

	token, err := middleware.CreateToken()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.ErrorResponse{
				Error:            "internal_server_error",
				ErrorDescription: "error creating token",
				Message:          "Something went wrong. Try again.",
			},
		)
	}

	response := fiber.Map{
		"salt":  result.Salt,
		"key":   result.Key,
		"token": token,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
