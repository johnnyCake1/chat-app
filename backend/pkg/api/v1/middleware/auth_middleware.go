package middleware

import (
	v1 "backend/pkg/api/v1"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// AuthMiddleware middleware to check JWT token in Authorization header
func AuthMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header missing",
			})
		}
		token := strings.Split(authHeader, " ")[1]
		userID, err := v1.ValidateAuthToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}
		// Set user ID in context for subsequent handlers
		c.Locals("userID", userID)
		return c.Next()
	}
}
