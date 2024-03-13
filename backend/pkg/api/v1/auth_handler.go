package v1

import (
	"backend/pkg/config"
	"backend/pkg/model"
	"backend/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body model.RegistrationRequest true "User registration request"
// @Success 201 {object} model.User
// @Router /api/v1/register [post]
func RegisterHandler(userService *service.UserService) fiber.Handler {
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
		user, err := userService.CreateNewUser(model.User{
			Nickname:     request.Nickname,
			PasswordHash: string(hashedPassword),
			Email:        request.Email,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't register the new user: %v", err),
			})
		}

		// Generate a jwt token
		token, err := GenerateAuthToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't register the new user: %v", err),
			})
		}

		// Store the token into client's Cookie
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
// @Summary Log in a user
// @Description Log in a user with the provided credentials
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "User login request"
// @Success 200 {object} model.User
// @Router /api/v1/login [post]
func LoginHandler(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(model.LoginRequest)
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't parse the request: %v", err),
			})
		}

		// Retrieve user from DB and check password...
		user, err := userService.GetUserByEmail(request.Email)
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
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Couldn't login: %v", err),
			})
		}

		// generate a jwt token
		token, err := GenerateAuthToken(user.ID)
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
// @Summary Log out a user
// @Description Log out the currently authenticated user
// @Tags Authentication
// @Success 200 {object} map[string]string
// @Router /api/v1/logout [post]
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

// ValidateTokenHandler handles token validation. It validates the token provided in the query parameter. Otherwise, validates the token provided in the cookie.
// @Summary Validate a JWT token
// @Description Validate the JWT token provided in the query parameter or cookie
// @Tags Authentication
// @Param token query string false "JWT token"
// @Success 200
// @Failure 401 {object} map[string]string
// @Router /api/v1/validateToken [get]
func ValidateTokenHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			token = c.Cookies("jwt")
		}
		if err := ValidateAuthToken(token); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		return nil
	}
}

// GenerateAuthToken generates a JWT token for a given user ID
func GenerateAuthToken(userID uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 3 days
	})
	token, err := claims.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("couldn't generate token: %v", err)
	}
	return token, nil
}

func ValidateAuthToken(token string) error {
	if token == "" {
		return fmt.Errorf("empty token provided")
	}
	//Token validation logic
	_, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})
	//userID := parsedToken.Claims.(*jwt.RegisteredClaims).Issuer // this is how to retrieve the id of the request's user

	return err
}
