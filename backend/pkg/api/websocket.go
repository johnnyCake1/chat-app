package api

import (
	v1 "backend/pkg/api/v1"
	"backend/pkg/model"
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/streadway/amqp"
	"log"
)

// WebSocketHandler Upgrade HTTP request to WebSocket
func WebSocketHandler(messageChannel *amqp.Channel) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		var (
			_   int
			msg []byte
			err error
		)

		for {
			// Read message from WebSocket client
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("Error reading message from a websocket:", err)
				break // Exit the loop and close connection on read error
			}

			// Handle message
			var message model.Message
			if err = json.Unmarshal(msg, &message); err != nil {
				log.Printf("Error parsing message: %s", err)
				continue // Skip this message but keep the connection alive
			}

			// Publish message to RabbitMQ
			if err := v1.SendToQueue(message, messageChannel); err != nil {
				log.Printf("Error publishing message: %s", err)
				continue // Skip this message but keep the connection alive
			}
		}
	}
}
