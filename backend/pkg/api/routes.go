package api

import (
	v1 "backend/pkg/api/v1"
	_ "backend/pkg/api/v1/docs"
	"backend/pkg/api/v1/middleware"
	"backend/pkg/config"
	"backend/pkg/consumer"
	"backend/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	"strings"
)

func SetupRoutes(app *fiber.App, services *service.Services, messageHub *consumer.MessageHub) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost,http://localhost:3000,http://frontend:3000",
		AllowHeaders: "Content-Type, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))
	app.Use("/ws", func(c *fiber.Ctx) error {
		requestedProtocol := c.Get("Sec-WebSocket-Protocol")
		// Split the string by the comma separator
		parts := strings.Split(requestedProtocol, ",")
		if len(parts) < 2 {
			fmt.Println("Invalid sub-protocol format. Expected a sub-protocol in websocket handshake connection followed by access token")
			return fiber.ErrBadRequest
		}
		// Extract the JWT token from the second part
		tokenStr := strings.TrimSpace(parts[1])
		userID, err := v1.ValidateAuthToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}
		// save the user id in Locals to be able to fetch it in websocket handler
		c.Locals("userID", userID)
		// check if the client requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Route for serving static files
	app.Static("/static", "./static")
	// WebSocket route
	app.Get("/ws", websocket.New(
		ClientWebSocketConnectionHandler(messageHub, services.ChatroomService),
		websocket.Config{
			Subprotocols: []string{config.WebsocketChatSubProtocol},
		},
	))
	app.Get("/", func(c *fiber.Ctx) error { // solely for server proxy testing
		return c.SendString("Chat app root")
	})
	// Grouping API version 1 prefix
	api := app.Group("/api/v1")
	// Swagger documentation route
	api.Get("/swagger/*", swagger.New(swagger.ConfigDefault))
	// Auth routes
	api.Post("/register", v1.RegisterHandler(services.UserService))
	api.Post("/login", v1.LoginHandler(services.UserService))
	api.Post("/logout", v1.LogoutHandler())
	api.Post("/validateToken", v1.ValidateTokenHandler())
	// User routes
	api.Get("/users", middleware.AuthMiddleware(), v1.GetUsers(services.UserService))
	api.Get("/users/search", middleware.AuthMiddleware(), v1.SearchUsers(services.UserService))
	api.Get("/users/:id", middleware.AuthMiddleware(), v1.GetUser(services.UserService))
	api.Delete("/users/:id", middleware.AuthMiddleware(), v1.GetUser(services.UserService))
	api.Get("/users/:id/chatrooms", middleware.AuthMiddleware(), v1.GetUserChatrooms(services.ChatroomService))
	// Chatrooms routes
	api.Get("/chatrooms/:id", middleware.AuthMiddleware(), v1.GetChatroomById(services.ChatroomService))
}
