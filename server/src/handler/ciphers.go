package handler

import (
	"errors"

	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/Nonz007x/pass-ez/src/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CreateItem(c *fiber.Ctx) error {
	var req model.Item

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.JsonParsingError)
	}

	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
	}

	userId, userIdOk := claims["user_id"].(string)

	if !userIdOk {
		return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
	}

	db := database.DB.Db

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user_vault model.UserVault
	if err := db.Where("user_id = ? AND vault_id = ?", userId, req.VaultID).First(&user_vault).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
	}

	item := model.Item{
		ID:      utils.UUID(),
		VaultID: req.VaultID,
		Name:    req.Name,
		Note:    req.Note,
		Type:    req.Type,
	}

	if err := db.Create(&item).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
