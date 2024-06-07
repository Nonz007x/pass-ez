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
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonParsingError)
	}

	email := req.Email

	// hash with SHA-256 and encode to base-64 string
	bytePassword := sha256.Sum256([]byte(req.Password))
	hashedPassword := base64.StdEncoding.EncodeToString(bytePassword[:])

	var result struct {
		Id   string
		Salt string
		Key  string
	}
	db := database.DB.Db
	query := `
		SELECT users.id, users.salt, vaults.key 
		FROM users 
		INNER JOIN vaults ON users.id = vaults.owner_id 
		WHERE email = ? AND master_password = ?;
	`
	err := db.Raw(query, email, hashedPassword).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.UserNotFoundError)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.InternalServerError)
	}

	if result.Key == "" || result.Salt == "" {
		return c.Status(fiber.StatusNotFound).JSON(models.UserNotFoundError)
	}

	token, err := middleware.CreateToken(result.Id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.InternalServerError)
	}

	response := fiber.Map{
		"salt":  result.Salt,
		"key":   result.Key,
		"token": token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
