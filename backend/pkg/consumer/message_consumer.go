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
}

// MessageHub is for managing clients connections and also publishing, consuming and broadcasting chat messages
type MessageHub struct {
	Clients             map[*Client]bool       // Keeps track of all connected Clients
	Register            chan *Client           // Channel for registering new clients
	Unregister          chan *Client           // Channel for unregistering clients
	Broadcast           chan model.ChatMessage // Channel for broadcasting messages to clients
	MessageQueueChannel *amqp.Channel          // connected RabbitMQ message channel for publishing/consuming chat messages
}

func NewMessageHub(messageQueueChannel *amqp.Channel) *MessageHub {
	return &MessageHub{
		Broadcast:           make(chan model.ChatMessage),
		Register:            make(chan *Client),
		Unregister:          make(chan *Client),
		Clients:             make(map[*Client]bool),
		MessageQueueChannel: messageQueueChannel,
	}
}

type MessageHandler interface {
	// HandleMessage handles the messageData and broadcasts it to the chatroom participants. Returns the updated messageData
	HandleMessage(*model.MessageData, *service.ChatroomService, map[*Client]bool) (*model.MessageData, error)
}

type SendMessageHandler struct{}
type ViewMessageHandler struct{}
type CreateGroupChatroomHandler struct{}
type CreatePrivateChatroomHandler struct{}
type UpdateGroupChatroomHandler struct{}
type DeleteGroupChatroomHandler struct{}
type EditMessageHandler struct{}
type DeleteMessageHandler struct{}
type ReactToMessageHandler struct{}

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
		log.Fatalf("Error consuming message data from message queue: %v", err)
	}
	log.Printf("Message data consumer service is running...")
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client connected (user id: %v). Number of clients: %v", client.UserID, len(h.Clients))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
				log.Printf("Client disconnected (user id: %v). Number of clients: %v", client.UserID, len(h.Clients))
			}
		case d := <-msgs:
			messageData := &model.MessageData{}
			if err := json.Unmarshal(d.Body, messageData); err != nil {
				log.Printf("Error unmarshalling message data: %v\n", err)
				continue
			}
			handler := getHandlerForMessageOption(model.MesssageOption(messageData.MessageOption))
			if handler == nil {
				log.Printf("No handler for message data option: %v\n", messageData.MessageOption)
				continue
			}
			_, err := handler.HandleMessage(messageData, chatroomService, h.Clients)
			if err != nil {
				log.Printf("Error handling message data with option %v : %v\n", messageData.MessageOption, err)
				continue
			}
		}
	}
}

// HandleMessage handles creating a new chat message and broadcasting it to the chatroom participants
func (h *SendMessageHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	if messageData.SendMessage == nil {
		return nil, fmt.Errorf("error sending message: SendMessage is not specified")
	}
	if messageData.SendMessage.ChatroomID == 0 {
		return nil, fmt.Errorf("error sending message: chatroomID is not specified")
	}
	message, err := chatroomService.AddMessageToChatroom(messageData.SendMessage.ChatroomID, model.ChatMessage{
		SenderID:      messageData.SendMessage.SenderID,
		Text:          messageData.SendMessage.Text,
		AttachmentURL: messageData.SendMessage.AttachmentURL,
	})
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}
	// append on response
	messageData.SendMessage.ChatMessage = *message
	for client := range clients {
		// If the client is a participant of the chatroom, send the messageData
		if client.ChatIDs[messageData.SendMessage.ChatroomID] {
			sendMessageToClient(client, messageData, model.MessageDataOptionSendMessage)
		}
	}
	return messageData, nil
}

func (h *ViewMessageHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	if messageData.ViewMessage == nil {
		return nil, fmt.Errorf("error marking message as viewed: ViewMessage is not specified")
	}
	if messageData.ViewMessage.MessageID == 0 {
		return nil, fmt.Errorf("error marking message as viewed: messageID is not specified")
	}
	if messageData.ViewMessage.ChatroomID == 0 {
		return nil, fmt.Errorf("error marking message as viewed: chatroomID is not specified")
	}
	updatedMessage, err := chatroomService.MarkMessageAsViewed(messageData.ViewMessage.MessageID)
	if err != nil {
		return nil, fmt.Errorf("error marking message as viewed: %v", err)
	}
	messageData.ChatMessage = updatedMessage
	for client := range clients {
		// If the client is a participant of the chatroom
		if client.ChatIDs[messageData.ViewMessage.ChatroomID] {
			// send the messageData to the client
			sendMessageToClient(client, messageData, model.MessageDataOptionViewMessage)
		}
	}
	return messageData, nil
}

func (h *CreateGroupChatroomHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	if messageData.CreateGroupChatroom == nil {
		return nil, fmt.Errorf("error creating group chatroom: CreateGroupChatroom is not specified")
	}
	if len(messageData.CreateGroupChatroom.Participants) == 0 {
		return nil, fmt.Errorf("error creating group chatroom: participants should be specified")
	}
	participantsIDs := getUsersIDs(messageData.CreateGroupChatroom.Participants)
	chatroom, err := chatroomService.CreateGroupChatroom(&model.ChatroomOptions{
		IsGroup:      true,
		GroupName:    messageData.CreateGroupChatroom.GroupName,
		Participants: participantsIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating group chatroom: %v", err)
	}
	messageData.CreateGroupChatroom.Chatroom = *chatroom
	for client := range clients {
		// If the client is a participant of the chatroom
		if exists(participantsIDs, client.UserID) {
			// update clients' chatroom subscriptions
			client.ChatIDs[chatroom.ID] = true
			// send the messageData to the client
			sendMessageToClient(client, messageData, model.MessageDataOptionCreateGroupChatroom)
		}
	}
	return messageData, nil
}

func (h *CreatePrivateChatroomHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	log.Printf("creating private chatroom: %v\n", messageData.CreatePrivateChatroom)
	if messageData.CreatePrivateChatroom == nil {
		return nil, fmt.Errorf("error creating private chatroom: CreatePrivateChatroom is not specified")
	}
	participant1 := messageData.CreatePrivateChatroom.Participants[0]
	participant2 := messageData.CreatePrivateChatroom.Participants[1]
	if participant1.ID == 0 || participant2.ID == 0 {
		return nil, fmt.Errorf("error creating private chatroom: both participants should be specified and have valid IDs")
	}
	if participant1.ID != messageData.CreatePrivateChatroom.ChatMessage.SenderID && participant2.ID != messageData.CreatePrivateChatroom.ChatMessage.SenderID {
		return nil, fmt.Errorf("error creating private chatroom: sender should be one of the participants")
	}
	if messageData.CreatePrivateChatroom.ChatMessage.Text == "" && messageData.CreatePrivateChatroom.ChatMessage.AttachmentURL == "" {
		return nil, fmt.Errorf("error creating private chatroom: chat message should be specified")
	}
	chatroom, err := chatroomService.CreatePrivateChatroom(messageData.CreatePrivateChatroom)
	if err != nil {
		return nil, fmt.Errorf("error creating private chatroom: %v", err)
	}

	log.Printf("created private chatroom object: %v\n", chatroom)
	messageData.CreatePrivateChatroom.Chatroom = *chatroom
	// update clients' chatroom subscriptions and send the messageData to the participants
	for client := range clients {
		isUserParticipant := client.UserID == participant1.ID || client.UserID == participant2.ID
		if isUserParticipant {
			client.ChatIDs[chatroom.ID] = true
			// resolve chatroom name and picture for the client because the name and picture look different for each participant of a private (1 to 1) chatroom
			otherParticipant := getOtherParticipant(client.UserID, participant1, participant2)
			messageData.CreatePrivateChatroom.ChatroomName = extractUserName(otherParticipant)
			messageData.CreatePrivateChatroom.ChatroomPictureURL = otherParticipant.AvatarURL
			sendMessageToClient(client, messageData, model.MessageDataOptionCreatePrivateChatroom)
		}
	}
	messageData.CreatePrivateChatroom.ChatroomName = ""
	messageData.CreatePrivateChatroom.ChatroomPictureURL = ""
	return messageData, nil
}

func (h *UpdateGroupChatroomHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	if messageData.UpdateGroupChatroom == nil {
		return nil, fmt.Errorf("error updating group chatroom: UpdateGroupChatroom is not specified")
	}
	if messageData.UpdateGroupChatroom.ID == 0 {
		return nil, fmt.Errorf("error updating group chatroom: chatroomID should be specified")
	}
	if messageData.UpdateGroupChatroom.GroupName == "" || len(messageData.UpdateGroupChatroom.Participants) == 0 {
		return nil, fmt.Errorf("error updating group chatroom: nothing to update")
	}
	chatroom, err := chatroomService.UpdateGroupChatroom(messageData.UpdateGroupChatroom)
	if err != nil {
		return nil, fmt.Errorf("error updating group chatroom: %v", err)

	}
	participantsIDs := getUsersIDs(messageData.UpdateGroupChatroom.Participants)
	for client := range clients {
		// Only if the client is a participant of the chatroom, send the messageData to that client
		if exists(participantsIDs, client.UserID) {
			client.ChatIDs[chatroom.ID] = true
			if err := client.Conn.WriteJSON(messageData); err != nil {
				log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, model.MessageDataOptionSendMessage, err)
				continue
			}
		}
	}
	return messageData, nil
}

func (h *DeleteGroupChatroomHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	// TODO: implement chatroom deletion
	return nil, nil
}

func (h *EditMessageHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	// TODO: implement messageData editing
	return nil, nil
}

func (h *DeleteMessageHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	// TODO: implement messageData deletion
	return nil, nil
}

func (h *ReactToMessageHandler) HandleMessage(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	// TODO: implement messageData reaction
	return nil, nil
}

func getHandlerForMessageOption(option model.MesssageOption) MessageHandler {
	switch option {
	case model.MessageDataOptionSendMessage:
		return &SendMessageHandler{}
	case model.MessageDataOptionViewMessage:
		return &ViewMessageHandler{}
	case model.MessageDataOptionCreateGroupChatroom:
		return &CreateGroupChatroomHandler{}
	case model.MessageDataOptionCreatePrivateChatroom:
		return &CreatePrivateChatroomHandler{}
	case model.MessageDataOptionUpdateGroupChatroom:
		return &UpdateGroupChatroomHandler{}
	default:
		return nil
	}
}

func getOtherParticipant(userID uint, participant1 model.User, participant2 model.User) model.User {
	if userID == participant1.ID {
		return participant2
	}
	return participant1
}

func sendMessageToClient(client *Client, messageData *model.MessageData, option model.MesssageOption) {
	err := client.Conn.WriteJSON(messageData)
	if err != nil {
		log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, option, err)
	}
}

func extractUserName(user model.User) string {
	if user.Nickname != "" {
		return user.Nickname
	}
	return user.Email
}

func getUsersIDs(users []model.User) []uint {
	ids := make([]uint, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}

func exists(slice []uint, item uint) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
