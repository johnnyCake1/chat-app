package v1

import (
	"backend/pkg/model"
	"backend/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetUsers fetches all users
// @Summary Get all users
// @Description Retrieve all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @Router /api/v1/users [get]
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
// @Summary Get a user
// @Description Retrieve information about a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Router /api/v1/users/{id} [get]
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
// @Summary Search users
// @Description Search for users by email or nickname containing the provided search term
// @Tags Users
// @Accept json
// @Produce json
// @Param searchTerm query string true "Search term"
// @Success 200 {array} model.User
// @Failure 400 {object} map[string]string
// @Router /api/v1/users/search [get]
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
