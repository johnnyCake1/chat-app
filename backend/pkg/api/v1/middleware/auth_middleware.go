package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

const jwtSecret = "my_secret_key" // TODO: generate and use own secret key

// GenerateToken generates a JWT token for a given user ID
func GenerateToken(userID uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 3 days
	})
	token, err := claims.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("couldn't generate token: %v", err)
	}
	return token, nil
}

// Protect is a middleware to protect routes from unauthorised access
func Protect() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		//Token validation logic
		token := c.Cookies("jwt")
		_, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorised. Please login/register first",
			})
		}

		//userID := parsedToken.Claims.(*jwt.RegisteredClaims).Issuer // this is how to retrieve the id of the request's user

		return c.Next()
	}
}
