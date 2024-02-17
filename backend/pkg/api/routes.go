package api

import (
	v1 "backend/pkg/api/v1"
	"backend/pkg/api/v1/middleware"
	"backend/pkg/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, repos *repository.Repositories) {
	// Middleware that applies to all requests
	app.Use(logger.New())
	// Enabling all origins only for development purposes!
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	// Grouping API version 1
	api := app.Group("/api/v1")
	// Applying middleware specific to API v1 routes
	api.Use(middleware.AuthMiddleware)
	// Testing route
	api.Get("/", v1.GetHelloWorld)
	// User routes
	api.Get("/user", v1.GetUsers(repos.UserRepo))
	api.Post("/user", v1.CreateUser(repos.UserRepo))
	api.Get("/user/:id", v1.GetUser(repos.UserRepo))

}
