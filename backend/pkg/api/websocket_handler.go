package api

import (
	"backend/pkg/config"
	"backend/pkg/consumer"
	"backend/pkg/service"
	"github.com/gofiber/websocket/v2"
	"github.com/streadway/amqp"
	"log"
)

// ClientWebSocketConnectionHandler handles client connection to the server websocket.
// @Summary Handle WebSocket connection
// @Description Handle client connections to the WebSocket server
// @Tags WebSocket
// @Success 101
// @Router /ws [get]
func ClientWebSocketConnectionHandler(messageHub *consumer.MessageHub, chatroomService *service.ChatroomService) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		userIDRaw := c.Locals("userID")
		userID, ok := userIDRaw.(uint)
		if !ok {
			log.Printf("failed to connect the client: invalid user id in the context")
			return
		}
		client := &consumer.Client{Conn: c, ChatIDs: make(map[uint]bool)}
		// TODO: send chatrooms to the client on connection through the websocket
		// retrieve chatrooms that the user is subscribed to
		chatrooms, err := chatroomService.GetChatroomsByUserId(userID, 1, 1)
		if err != nil {
			return
		}
		for _, chatroom := range chatrooms {
			client.ChatIDs[chatroom.ID] = true
		}
		client.UserID = userID
		messageHub.Register <- client

		defer func() {
			messageHub.Unregister <- client
			c.Close()
		}()

		// Handle messages from the client WebSocket and disconnect the client on failure
		for {
			var (
				msg []byte
				err error
			)
			// Listen for message from WebSocket client (blocking operation)
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Printf("Error reading message from a websocket: %v. Disconnecting the client\n", err)
				messageHub.Unregister <- client
				return
			}

			// If it's a ping message then ignore it
			if string(msg) == "PING" {
				log.Printf("Received `PING` from client with user %v\n", client.UserID)
				continue
			}

			// Publish message to RabbitMQ
			if err := PublishMessageToQueue(msg, messageHub.MessageQueueChannel); err != nil {
				log.Printf("Error publishing message: %v\n", err)
				continue
			}
		}
	}
}

func PublishMessageToQueue(msg []byte, ch *amqp.Channel) error {
	log.Printf("ChatMessage published to msg queue\n")
	err := ch.Publish(
		"",                           // exchange
		config.ChatMessageRoutingKey, // routing key
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
	return err
}
