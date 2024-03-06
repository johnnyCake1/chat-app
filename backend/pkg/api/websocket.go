package api

import (
	"backend/pkg/config"
	"backend/pkg/model"
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/streadway/amqp"
	"log"
)

type Client struct {
	conn    *websocket.Conn
	chatIDs map[uint]bool // Maps to keep track of which chat IDs the client is subscribed to
	send    chan []byte   // Channel for sending messages to the client
}

type Hub struct {
	clients    map[*Client]bool   // Keeps track of all connected clients
	register   chan *Client       // Channel for registering new clients
	unregister chan *Client       // Channel for unregistering clients
	broadcast  chan model.Message // Channel for broadcasting messages to clients
}

var MessageHub = Hub{
	broadcast:  make(chan model.Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("couldn't serialise the message to JSON: %v\n", err)
				break
			}
			for client := range h.clients {
				if client.chatIDs[message.ChatID] { // Check if client is subscribed to the chat ID
					select {
					case client.send <- messageBytes:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

// WebSocketHandler Upgrade HTTP request to WebSocket
func WebSocketHandler(messageChannel *amqp.Channel) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {

		client := &Client{conn: c, chatIDs: make(map[uint]bool), send: make(chan []byte, 256)}
		MessageHub.register <- client

		defer func() {
			MessageHub.unregister <- client
			c.Close()
		}()

		var (
			msg []byte
			err error
		)

		// Handle messages from the client WebSocket and disconnect the client on failure
		for {
			// Read message from WebSocket client.
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
			if err := PublishMessage(message, messageChannel); err != nil {
				log.Printf("Error publishing message: %v\n", err)
				break
			}
		}
	}
}

func PublishMessage(message model.Message, ch *amqp.Channel) error {
	// Serialize message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

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
