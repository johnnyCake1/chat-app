package v1

import (
	"backend/pkg/model"
	"backend/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetUsers fetches all users
func GetUsers(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := userService.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		return c.JSON(users)
	}
}

// GetUser gets a specific user with given id
func GetUser(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
		}

		user, err := userService.GetUserByID(uint(userID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Couldn't query the user from database: %v", err)})
		}
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		return c.JSON(user)
	}
}

// SearchUsers fetches all users containing the searchTerm in their email or nickname
func SearchUsers(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the search term from the query parameter
		searchTerm := c.Query("searchTerm")

		// Check if searchTerm is empty
		if searchTerm == "" {
			return c.JSON([]model.User{})
		}

		users, err := userService.GetUsersBySearchTerm(searchTerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error searching users: %v", err)})
		}

		return c.JSON(users)
	}
}
