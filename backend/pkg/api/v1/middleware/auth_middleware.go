package middleware

import (
	v1 "backend/pkg/api/v1"
	"github.com/gofiber/fiber/v2"
)

// Protect is a middleware to protect routes from unauthorised access
func Protect() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("jwt")
		if err := v1.ValidateAuthToken(token); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorised. Please login/register first",
			})
		}

		return c.Next()
	}
}
