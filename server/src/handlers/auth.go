package handlers

import (
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

	user := models.User{
		Id:             uuid.New().String(),
		Email:          req.Email,
		MasterPassword: req.Password,
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

	response := "status: success"
	return c.Status(fiber.StatusOK).JSON(response)
}
