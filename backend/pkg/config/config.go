package config

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/streadway/amqp"
	"log"
)

type Client struct {
	Conn    *websocket.Conn
	ChatIDs map[uint]bool // Maps to keep track of which chat IDs the client is subscribed to
	Send    chan []byte   // Channel (buffered) for sending messages to the client
}

type Hub struct {
	Clients             map[*Client]bool   // Keeps track of all connected Clients
	Register            chan *Client       // Channel for registering new clients
	Unregister          chan *Client       // Channel for unregistering clients
	Broadcast           chan model.Message // Channel for broadcasting messages to clients
	MessageQueueChannel *amqp.Channel      // connected RabbitMQ message channel for publishing/consuming chat messages
}

type AppDependencies struct {
	Repos               *repository.Repositories // Repositories with database connection
	MessageQueueChannel *amqp.Channel            // connected RabbitMQ message channel for chat messages
	MessageHub          *Hub                     // MessageHub is for managing clients connections and also publishing, consuming and broadcasting chat messages
}

// StartMessageConsumerService Connects to message queue and consumes messages to broadcast them. It also listens for client registration/unregistration to add/delete the clients to broadcast.  It's a blocking function, so you should run it in a goroutine
func (h *Hub) StartMessageConsumerService() {
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
		ChatMessageQueueName, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Fatalf("Error consuming messages from Message Queue: %v", err)
	}
	// debug log:
	log.Printf("Message Consumer Service is running...")
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		//case message := <-h.broadcast: {log.Printf("Message received for broadcasting: %v", message)} // We could also directly use broadcast channel instead of consuming from a message queue
		case d := <-msgs:
			// debug log
			log.Printf("Message consumed: %v", d.Body)
			var message model.Message
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Error parsing message: %v\n", err)
				continue
			}
			for client := range h.Clients {
				//if client.ChatIDs[message.ChatRoomID] {
				messageBytes, err := json.Marshal(message)
				if err != nil {
					log.Printf("couldn't serialise the message to JSON: %v\n", err)
					continue
				}
				// Process the message
				log.Printf("Server received a message: %v\tSending to client:%v\n", message.Text, client.Conn)
				// Write the message to the WebSocket connection
				if err := client.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
					log.Printf("Error writing message to client: %v\n", err)
					continue
				}
				//}
			}
		}
	}
}

const ServerPort = 8080

const ChatMessageQueueName = "chat_messages"

const ChatMessageRoutingKey = ChatMessageQueueName
