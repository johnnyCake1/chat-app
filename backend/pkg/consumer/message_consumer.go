package consumer

import (
	"backend/pkg/config"
	"backend/pkg/model"
	"backend/pkg/service"
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/streadway/amqp"
	"log"
)

// Client is for storing clients connections and keeping track of chats that the client is subscribed to
type Client struct {
	Conn    *websocket.Conn
	ChatIDs map[uint]bool // Maps to keep track of which chat IDs the client is subscribed to
	UserID  uint          // User ID of the client
	//Send    chan []byte   // Channel (buffered) for sending messages to the client // We actually use MessageHub.MessageQueueChannel RMQ channel to publish the message and use the Conn websocket connection for broadcasting
}

// MessageHub is for managing clients connections and also publishing, consuming and broadcasting chat messages
type MessageHub struct {
	Clients             map[*Client]bool   // Keeps track of all connected Clients
	Register            chan *Client       // Channel for registering new clients
	Unregister          chan *Client       // Channel for unregistering clients
	Broadcast           chan model.Message // Channel for broadcasting messages to clients
	MessageQueueChannel *amqp.Channel      // connected RabbitMQ message channel for publishing/consuming chat messages
}

func NewMessageHub(messageQueueChannel *amqp.Channel) *MessageHub {
	return &MessageHub{
		Broadcast:           make(chan model.Message),
		Register:            make(chan *Client),
		Unregister:          make(chan *Client),
		Clients:             make(map[*Client]bool),
		MessageQueueChannel: messageQueueChannel,
	}
}

// StartMessageConsumerService Connects to message queue and consumes messages to broadcast them. It also listens for client registration/unregistration to add/delete the clients to broadcast.  It's a blocking function, so you should run it in a goroutine
func (h *MessageHub) StartMessageConsumerService(chatroomService *service.ChatroomService) {
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
	if err != nil {
		log.Fatalf("Error consuming messages from Message Queue: %v", err)
	}
	log.Printf("Message Consumer Service is running...")
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
			}
		//case message := <-h.broadcast: {log.Printf("Message received for broadcasting: %v", message)} // We could also directly consume from a broadcast channel instead of consuming from a message queue
		case d := <-msgs:
			var message model.MessageData
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Error unmarshalling message: %v\n", err)
				continue
			}
			log.Printf("Message consumed from message queue: %v", message)
			// Handle message in the server
			if message.Message, err = handleMessage(message, chatroomService); err != nil {
				log.Printf("Error handling message with option %v : %v\n", message.MessageOption, err)
				continue
			}
			// Handle message broadcasting to clients
			log.Printf("Broadcasting message to %v clients", len(h.Clients))
			for client := range h.Clients {
				// If the client is subscribed to the chatroom of the message
				if client.ChatIDs[message.ChatRoomID] {
					log.Printf("sending message to client %v with option %v", client.UserID, message.MessageOption)
					if err := client.Conn.WriteJSON(message); err != nil {
						log.Printf("Error sending message to client with option %v: %v\n", message.MessageOption, err)
						continue
					}
				}
			}
		}
	}
}

func handleMessage(message model.MessageData, chatroomService *service.ChatroomService) (model.Message, error) {
	switch message.MessageOption {
	case model.MessageOptionView:
		return chatroomService.MarkMessageAsViewed(message.ID)
	case model.MessageOptionSend:
		return chatroomService.AddMessageToChatroom(message.ChatRoomID, message.Message)
	case model.MessageOptionEdit:
		// TODO: implement message editing
		return model.Message{}, nil
	case model.MessageOptionDelete:
		// TODO: implement message deletion
		return model.Message{}, nil
	case model.MessageOptionReaction:
		// TODO: implement message reaction
		return model.Message{}, nil
	}
	return model.Message{}, fmt.Errorf("unknown message option from message: %v", message)
}
