package api

import (
	v1 "backend/pkg/api/v1"
	"backend/pkg/api/v1/middleware"
	"backend/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, appDependencies *config.AppDependencies) {
	// Middleware that applies to all requests
	app.Use(logger.New())
	// Enabling all origins only for development purposes!
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
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
	app.Get("/ws", websocket.New(WebSocketHandler(appDependencies.MessageChannel)))

	app.Get("/", func(c *fiber.Ctx) error { // solely for server proxy testing
		return c.SendString("Chat app root")
	})
	// Grouping API version 1
	api := app.Group("/api/v1")
	// Applying middleware specific to API v1 routes
	api.Use(middleware.AuthMiddleware)
	api.Get("/", func(c *fiber.Ctx) error { // solely for server proxy testing
		return c.SendString("Chat app api v1 root")
	})
	// User routes
	api.Get("/user", v1.GetUsers(appDependencies.Repos.UserRepo))
	api.Post("/user", v1.CreateUser(appDependencies.Repos.UserRepo))
	api.Get("/user/:id", v1.GetUser(appDependencies.Repos.UserRepo))
	// Message routes
	api.Post("/message", v1.SendMessage(appDependencies.MessageChannel))

}
