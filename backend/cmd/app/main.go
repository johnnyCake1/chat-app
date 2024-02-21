package main

import (
	"backend/pkg/api"
	"backend/pkg/config"
	"backend/pkg/model"
	"backend/pkg/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// Database connection
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// RabbitMQ connection and message queue declaration for message sending
	conn, ch, err := connectToRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// Start the message receiving service
	go startMessageReceiverService()

	// Initialise all repositories
	repos := repository.InitRepositories(db)

	app := fiber.New()

	// Setup routes and pass dependencies
	api.SetupRoutes(
		app,
		&config.AppDependencies{
			Repos:          repos,
			MessageChannel: ch,
		},
	)

	// Start server on indicated port
	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.ServerPort)))
}

func connectToDB() (*sql.DB, error) {
	connStr := "host=postgres dbname=chatapp_db user=root password=rootuser sslmode=disable"
	return sql.Open("postgres", connStr)
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

// startMessageReceiverService Connects to queue and listens for messages. It's a blocking function, so you should run it in a goroutine
func startMessageReceiverService() {
	conn, err := amqp.Dial("amqp://root:rootuser@rabbitmq/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		config.ChatMessageQueueName, // queue
		"",                          // consumer
		true,                        // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)

	for d := range msgs {
		var message model.Message
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			log.Printf("Error parsing message: %s", err)
			continue
		}

		// Process the message
		log.Printf("Received a message: %s", message.Text)
	}
}
