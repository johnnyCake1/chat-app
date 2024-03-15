package main

import (
	"backend/pkg/api"
	_ "backend/pkg/api/v1/docs"
	"backend/pkg/config"
	"backend/pkg/consumer"
	"backend/pkg/repository"
	"backend/pkg/service"
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pressly/goose/v3"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// Database connection
	db, err := initAndConnectToDB()
	if err != nil {
		log.Fatal(fmt.Sprintf("couldn't connect to database: %v", err))
	}
	defer db.Close()

	// RabbitMQ connection and message queue declaration for message sending
	messageQueueConn, messageQueueChannel, err := connectToRabbitMQ()
	if err != nil {
		log.Fatal(fmt.Sprintf("couldn't connect to RabbitMQ: %v", err))
	}
	defer messageQueueConn.Close()

	// Initialise all repositories
	repos := repository.InitRepositories(db)
	// Initialise all services
	services := service.InitServices(repos)

	// Start the message consumer service
	MessageHub := consumer.NewMessageHub(messageQueueChannel)
	go MessageHub.StartMessageConsumerService(services.ChatroomService)

	app := fiber.New()

	// Setup routes and inject dependencies
	api.SetupRoutes(
		app,
		services,
		MessageHub,
	)

	// Start server on indicated port
	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.ServerPort)))
}

func initAndConnectToDB() (*sql.DB, error) {
	// connect to postgres
	connStr := "host=postgres dbname=chatapp_db user=root password=rootuser sslmode=disable"
	driverName := "postgres"
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %v", err)
	}
	// apply migrations
	_ = goose.SetDialect(driverName)
	if err := goose.RunContext(context.Background(), "up", db, "./migrations"); err != nil {
		return nil, fmt.Errorf("couldn't run migrations: %v", err)
	}
	log.Println("Migrations successfully applied")
	return db, nil
}

func connectToRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://root:rootuser@rabbitmq/")
	if err != nil {
		return nil, nil, err
	}
	// message channel and queue creation
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	_, err = ch.QueueDeclare(
		config.ChatMessageQueueName, // queue name
		false,                       // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)

	return conn, ch, nil
}
