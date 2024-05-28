package handlers

import (
	"crypto/sha256"
	"encoding/base64"

	"errors"

	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/Nonz007x/pass-ez/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot parse JSON",
		})
	}

	db := database.DB.Db

	var existingUser models.User
	err := db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error:            "conflict",
			ErrorDescription: "email_already_in_use",
			Message:          "Email is already in use. Try again with a different email.",
		})
	} else if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:            "database_error",
			ErrorDescription: "internal_error",
			Message:          "Internal server error. Please try again later.",
		})
	}

	// hash with SHA-256 and encode to base-64 string
	bytePassword := sha256.Sum256([]byte(req.Password))
	hashedPassword := base64.StdEncoding.EncodeToString(bytePassword[:])

	user := models.User{
		Id:             uuid.New().String(),
		Email:          req.Email,
		MasterPassword: hashedPassword,
		Salt:           req.Salt,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.ErrorResponse{
				Error:            "internal_server_error",
				ErrorDescription: "database_error",
				Message:          "something went wrong. Try again.",
			},
		)
	}

	vault := models.Vault{
		Id:      uuid.New().String(),
		Key:     req.VaultKey,
		OwnerId: user.Id,
	}

	if err := db.Create(&vault).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.ErrorResponse{
				Error:            "internal_server_error",
				ErrorDescription: "database_error",
				Message:          "something went wrong. Try again.",
			},
		)
	}

	user_vault := models.UserVault{
		VaultId: vault.Id,
		UserId:  user.Id,
	}

	if err := db.Create(&user_vault).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.ErrorResponse{
				Error:            "internal_server_error",
				ErrorDescription: "database_error",
				Message:          "something went wrong. Try again.",
			},
		)
	}

	response := "status: success"
	return c.Status(fiber.StatusOK).JSON(response)
}

func GetSalt(c *fiber.Ctx) error {

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

	response := fiber.Map{
		"salt": result.Salt,
		"key":  result.Key,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
