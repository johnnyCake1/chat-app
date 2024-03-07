package v1

import (
	"backend/pkg/api/v1/middleware"
	"backend/pkg/model"
	"backend/pkg/repository"
	"backend/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

// GetUsers fetches all users
func GetUsers(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userService := service.NewUserService(repo)
		users, err := userService.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		return c.JSON(users)
	}
}

// GetUser gets a specific user with given id
func GetUser(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			// If parsing fails, return a 400 Bad Request response
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
		}

		user, err := repo.FindByID(uint(userID))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		return c.JSON(user)
	}
}

// RegisterHandler handles user registration
func RegisterHandler(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(model.RegistrationRequest)
		if err := c.BodyParser(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't register: Couldn't parse the request: %v", err),
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't register the new user: %v", err),
			})
		}
		// Save new user to DB
		user, err := repo.AddNewUser(model.User{
			Nickname:     request.Nickname,
			PasswordHash: string(hashedPassword),
			Email:        request.Email,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't register the new user: %v", err),
			})
		}

		// generate a jwt token
		token, err := middleware.GenerateToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't register the new user: %v", err),
			})
		}

		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(72 * time.Hour),
			HTTPOnly: true,
		}
		c.Cookie(&cookie)

		return c.JSON(user)
	}

}

// LoginHandler handles user login
func LoginHandler(repo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(model.LoginRequest)
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't parse the request: %v", err),
			})
		}

		// Retrieve user from DB and check password...
		user, err := repo.FindByEmail(request.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't login: %v", err),
			})
		}
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't login: User with email %v not found", request.Email),
			})
		}
		log.Printf("Login username: %v and paswwrod: %v", user.Email, user.PasswordHash)
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't login: %v", err),
			})
		}

		// generate a jwt token
		token, err := middleware.GenerateToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't login: %v", err),
			})
		}

		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(72 * time.Hour),
			HTTPOnly: true,
		}
		c.Cookie(&cookie)

		return c.JSON(user)
	}
}

// LogoutHandler handles user logout
func LogoutHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		}
		c.Cookie(&cookie)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
		})
	}
}
