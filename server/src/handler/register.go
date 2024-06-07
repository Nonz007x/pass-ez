package handler

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/Nonz007x/pass-ez/src/model"
	"github.com/Nonz007x/pass-ez/src/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonParsingError)
	}

	db := database.DB.Db

	var existingUser models.User
	err := db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(models.EmailConflictError)
	} else if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(models.InternalServerError)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// hash with SHA-256 and encode to base-64 string
	bytePassword := sha256.Sum256([]byte(req.Password))
	hashedPassword := base64.StdEncoding.EncodeToString(bytePassword[:])

	user := models.User{
		Id:             util.UUID(),
		Email:          req.Email,
		MasterPassword: hashedPassword,
		Salt:           req.Salt,
	}

	if err := db.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.InternalServerError)
	}

	vault := models.Vault{
		Id:      util.UUID(),
		Key:     req.VaultKey,
		OwnerId: user.Id,
	}

	if err := db.Create(&vault).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.InternalServerError)
	}

	userVault := models.UserVault{
		VaultId: vault.Id,
		UserId:  user.Id,
	}

	if err := db.Create(&userVault).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.InternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.InternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
