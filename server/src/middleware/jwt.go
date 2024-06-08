package middleware

import (
	"os"
	"time"

	"github.com/Nonz007x/pass-ez/src/model"
	"github.com/Nonz007x/pass-ez/src/util"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const EXPIRE_TIME = time.Minute * 60

var jwtSecret = os.Getenv("JWT_SECRET")

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(jwtSecret)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}

func CreateToken(userId string) (model.Token, error) {
	var msgToken model.Token
	claims := jwt.MapClaims{
		"sub":     util.UUID(),
		"user_id": userId,
		"exp":     time.Now().Add(EXPIRE_TIME).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return msgToken, err
	}

	msgToken.AccessToken = t

	claims = jwt.MapClaims{
		"sub": util.UUID(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	rt, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return msgToken, err
	}
	msgToken.RefreshToken = rt
	return msgToken, nil
}

func ValidateToken(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
	}

	_, subOk := claims["sub"].(string)
	_, expOk := claims["exp"].(float64)
	_, userIdOk := claims["user_id"].(string)

	if !subOk || !expOk || !userIdOk {
		return c.Status(fiber.StatusUnauthorized).JSON(model.InvalidTokenError)
	}

	return c.SendStatus(fiber.StatusOK)
}
