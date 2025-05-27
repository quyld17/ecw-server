package jwt

import (
	"database/sql"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	users "github.com/quyld17/E-Commerce-Website/entities/user"
)

func Generate(email string, db *sql.DB) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	role, err := users.GetRole(email, db)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["role"] = role
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetToken(c echo.Context) string {
	token := c.Request().Header.Get("Authorization")
	return token
}

func GetClaims(token *jwt.Token, key string) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}
