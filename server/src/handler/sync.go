package handler

import (
	"fmt"

	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/Nonz007x/pass-ez/src/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ItemWithLogin struct {
	*model.Item
	*model.Login
}

func Sync(c *fiber.Ctx) error {

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

	var vaults []struct {
		ID  string
		Key string
	}

	query := `
		SELECT vaults.id, vaults.key
		FROM user_vault
		INNER JOIN vaults ON user_vault.vault_id = vaults.id
		WHERE user_vault.user_id = ?
	`
	if err := db.Raw(query, userId).Scan(&vaults).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
	}

	var items []interface{}

	for _, v := range vaults {
		var itemTypes []int
		if err := db.Raw(`
			SELECT type
			FROM items
			WHERE items.vault_id = ?
		`, v.ID).Scan(&itemTypes).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
		}

		for _, itemType := range itemTypes {
			switch itemType {
			case 1:
				var loginItems []ItemWithLogin
				if err := db.Raw(`
				SELECT items.id, items.vault_id, items.name, items.note, items.type, items.created, items.updated,
              logins.item_id, logins.username, logins.password, logins.password_revision_date
				FROM items
				INNER JOIN logins ON items.id = logins.item_id
				WHERE items.vault_id = ?
				`, v.ID).Scan(&loginItems).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
				}

				for _, res := range loginItems {
					var item *model.Item
					var login *model.Login
					if res.Item != nil {
						item = &model.Item{
							ID:      res.Item.ID,
							VaultID: res.Item.VaultID,
							Name:    res.Item.Name,
							Note:    res.Item.Note,
							Type:    res.Item.Type,
							Created: res.Item.Created,
							Updated: res.Item.Updated,
						}
					}

					if res.Login != nil {
						login = &model.Login{
							ItemID:               res.Login.ItemID,
							Username:             res.Login.Username,
							Password:             res.Login.Password,
							PasswordRevisionDate: res.Login.PasswordRevisionDate,
						}
					}

					itemWithLogin := &ItemWithLogin{item, login}
					items = append(items, itemWithLogin)
				}
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
			}
		}
	}

	fmt.Println(items)
	fmt.Println(vaults)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"vaults":  vaults,
		"ciphers": items,
	})
}
