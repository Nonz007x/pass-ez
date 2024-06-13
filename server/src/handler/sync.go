package handler

import (
	"time"

	"github.com/Nonz007x/pass-ez/src/database"
	"github.com/Nonz007x/pass-ez/src/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ItemWithLogin struct {
	*model.Item
	*model.Login
}

type ItemResponse struct {
	ID      string         `json:"id"`
	VaultID string         `json:"vault_id"`
	Name    string         `json:"name"`
	Note    string         `json:"note"`
	Type    int16          `json:"type"`
	Created time.Time      `json:"created"`
	Updated time.Time      `json:"updated"`
	Login   *LoginResponse `json:"login,omitempty"`
}

type LoginResponse struct {
	Username             string                    `json:"username"`
	Password             string                    `json:"password"`
	PasswordRevisionDate time.Time                 `json:"password_revision_date"`
	URIs                 []URIResponse             `json:"uris,omitempty"`
	PasswordHistories    []PasswordHistoryResponse `json:"password_history,omitempty"`
}

type URIResponse struct {
	URI string `json:"uri"`
}

type PasswordHistoryResponse struct {
	Password string    `json:"password"`
	LastUsed time.Time `json:"last_used"`
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
              logins.username, logins.password, logins.password_revision_date
				FROM items
				LEFT JOIN logins ON items.id = logins.item_id
				WHERE items.vault_id = ?
				`, v.ID).Scan(&loginItems).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
				}

				for _, res := range loginItems {
					item := ItemResponse{
						ID:      res.Item.ID,
						VaultID: res.Item.VaultID,
						Name:    res.Item.Name,
						Note:    res.Item.Note,
						Type:    res.Item.Type,
						Created: res.Item.Created,
						Updated: res.Item.Updated,
					}

					if res.Login != nil {
						var uris []URIResponse
						var passwordHistories []PasswordHistoryResponse

						err := db.Raw(`
						 	SELECT uri
							FROM uris
							WHERE login_id = ?
						`, res.ID).Scan(&uris).Error
						if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON("model.InternalServerError")
						}

						err = db.Raw(`
						 	SELECT password, last_used
							FROM password_history
							WHERE login_id = ?
						`, res.ID).Scan(&passwordHistories).Error
						if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON("model.InternalServerError")
						}

						item.Login = &LoginResponse{
							Username:             res.Login.Username,
							Password:             res.Login.Password,
							PasswordRevisionDate: res.Login.PasswordRevisionDate,
							URIs:                 uris,
							PasswordHistories:    passwordHistories,
						}
					}

					items = append(items, item)
				}
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(model.InternalServerError)
			}
		}
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"vaults":  vaults,
		"ciphers": items,
	})
}
