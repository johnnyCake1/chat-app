package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// TimeResponse struct to send the time
type TimeResponse struct {
	ServerTime string `json:"serverTime"`
}

func main() {
	app := fiber.New()

	// Enabling all origins only for development purposes!
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// Database connection parameters
	connStr := "host=postgres dbname=chatapp_db user=root password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Endpoint to retrieve a current time from the database
	app.Get("/api/v1/time", func(c *fiber.Ctx) error {
		var currentTime string
		err := db.QueryRow("SELECT NOW()").Scan(&currentTime)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				SendString("Error querying the database")
		}

		// Respond with the time
		response := TimeResponse{
			ServerTime: currentTime,
		}
		return c.JSON(response)
	})

	log.Fatal(app.Listen(":8080"))
}
