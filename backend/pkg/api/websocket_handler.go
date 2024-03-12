package api

import (
	"backend/pkg/config"
	"backend/pkg/consumer"
	"backend/pkg/model"
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/streadway/amqp"
	"log"
)

// TODO: implement user identification on each request. Also add chatrooms to the clients that they are subscribed to. For that we will also need the connected user identity to see in what chats he is a participant

// ClientWebSocketConnectionHandler handles client connection to the server websocket.
func ClientWebSocketConnectionHandler(messageHub *consumer.MessageHub) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		client := &consumer.Client{Conn: c, ChatIDs: make(map[uint]bool)}
		messageHub.Register <- client

		defer func() {
			messageHub.Unregister <- client
			c.Close()
		}()

		var (
			msg []byte
			err error
		)

		// Handle messages from the client WebSocket and disconnect the client on failure
		for {
			// Listen for message from WebSocket client.
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Printf("Error reading message from a websocket: %v\n", err)
				break
			}

			// Handle message
			var message model.Message
			if err = json.Unmarshal(msg, &message); err != nil {
				log.Printf("Error parsing message: %v\n", err)
				break
			}

			// Publish message to RabbitMQ
			if err := PublishMessageToQueue(message, messageHub.MessageQueueChannel); err != nil {
				log.Printf("Error publishing message: %v\n", err)
				break
			}
		}
	}
}

func PublishMessageToQueue(message model.Message, ch *amqp.Channel) error {
	// Serialize message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// debugging log:
	log.Printf("Message published to message queue:%v\n", message.Text)
	err = ch.Publish(
		"",                           // exchange
		config.ChatMessageRoutingKey, // routing key
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
	return err
}
