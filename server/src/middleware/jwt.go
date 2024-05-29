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
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
	})
}

func CreateToken() (models.Token, error) {
	var msgToken models.Token
	claims := jwt.MapClaims{
		"sub": util.UUID(),
		"exp": time.Now().Add(EXPIRE_TIME).Unix(),
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
