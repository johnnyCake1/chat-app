package api

import (
	v1 "backend/pkg/api/v1"
	"backend/pkg/api/v1/middleware"
	"backend/pkg/consumer"
	"backend/pkg/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, services *service.Services, messageHub *consumer.MessageHub) {
	app.Use(logger.New())
	// Enabling all origins only for development purposes!
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost,http://localhost:3000,http://frontend:3000",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
		AllowHeaders:     "Content-Type, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	// WebSocket route
	app.Get("/ws", websocket.New(ClientWebSocketConnectionHandler(messageHub)))
	app.Get("/", func(c *fiber.Ctx) error { // solely for server proxy testing
		return c.SendString("Chat app root")
	})
	// Grouping API version 1
	api := app.Group("/api/v1")
	// Auth routes
	api.Post("/register", v1.RegisterHandler(services.UserService))
	api.Post("/login", v1.LoginHandler(services.UserService))
	api.Post("/logout", v1.LogoutHandler())
	// User routes
	api.Get("/users", middleware.Protect(), v1.GetUsers(services.UserService))
	api.Get("/users/search", middleware.Protect(), v1.SearchUsers(services.UserService))
	api.Get("/users/:id", middleware.Protect(), v1.GetUser(services.UserService))
	api.Delete("/users/:id", middleware.Protect(), v1.GetUser(services.UserService))
	// Chatrooms
	api.Get("chatrooms", middleware.Protect(), v1.GetUser(services.UserService))
}
