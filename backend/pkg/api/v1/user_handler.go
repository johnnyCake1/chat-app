package v1

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

// GetUsers fetches all users
func GetUsers(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userService := service.NewUserService(repo)
		users, err := userService.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(users)
	}
}

// GetUser gets a specific user with given id
func GetUser(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("user")
	}
}

// CreateUser creates a user
func CreateUser(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(model.User)
		err := c.BodyParser(body)
		if err != nil {
			_ = c.Status(fiber.StatusBadRequest).SendString(err.Error())
			return err
		}
		return c.Status(fiber.StatusOK).JSON(body)
	}
}
