package main

import (
	"backend/pkg/api"
	"backend/pkg/config"
	"backend/pkg/repository"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// Database connection
	connStr := "host=postgres dbname=chatapp_db user=root password=rootuser sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Initialise all repositories
	repos := repository.InitRepositories(db)

	app := fiber.New()

	// Setup routes and pass repositories
	api.SetupRoutes(app, repos)

	// Start server
	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.ServerPort)))
}
