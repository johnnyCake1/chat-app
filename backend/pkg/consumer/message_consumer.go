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
		log.Fatalf("Error consuming messages from ChatMessage Queue: %v", err)
	}
	log.Printf("ChatMessage Consumer Service is running...")
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
			}
		case d := <-msgs:
			messageData := &model.MessageData{}
			if err := json.Unmarshal(d.Body, messageData); err != nil {
				log.Printf("Error unmarshalling messageData: %v\n", err)
				continue
			}
			log.Printf("messageData consumed from message queue: %v\n", messageData)
			// Handle messageData in the server
			_, err := handleMessageNew(messageData, chatroomService, h.Clients)
			if err != nil {
				log.Printf("Error handling messageData with option %v : %v\n", messageData.MessageOption, err)
				continue
			}
			// TODO: This is commented out because broadcasting is done in handleMessageNew because we might need to broadcast different messageData to different clients
			//// Handle messageData broadcasting to clients
			//log.Printf("Broadcasting messageData to %v clients\n", len(h.Clients))
			//for client := range h.Clients {
			//	// If the client is a participant of the chatroom, send the messageData
			//	if (updatedMessageData.ChatMessage != nil && client.ChatIDs[updatedMessageData.ChatMessage.ChatroomID]) ||
			//		updatedMessageData.ChatroomOptions != nil && client.ChatIDs[updatedMessageData.ChatroomOptions.ID] {
			//		if err := client.Conn.WriteJSON(updatedMessageData); err != nil {
			//			log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, messageData.MessageOption, err)
			//			continue
			//		}
			//	}
			//}
		}
	}
}

// handleMessageNew is a new implementation of handleMessage where we use actions based on the messageData option provided
func handleMessageNew(messageData *model.MessageData, chatroomService *service.ChatroomService, clients map[*Client]bool) (*model.MessageData, error) {
	if messageData == nil {
		return nil, fmt.Errorf("messageData is nil")
	}
	switch messageData.MessageOption {
	case model.MessageDataOptionSendMessage:
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
				if err := client.Conn.WriteJSON(messageData); err != nil {
					log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, model.MessageDataOptionSendMessage, err)
					continue
				}
			}
		}
		return messageData, nil
	case model.MessageDataOptionViewMessage:
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
			// If the client is a participant of the chatroom, send the messageData
			if client.ChatIDs[messageData.ViewMessage.ChatroomID] {
				if err := client.Conn.WriteJSON(messageData); err != nil {
					log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, model.MessageDataOptionViewMessage, err)
					continue
				}
			}
		}
		return messageData, nil
	case model.MessageDataOptionCreateGroupChatroom:
		if messageData.CreateGroupChatroom == nil {
			return nil, fmt.Errorf("error creating group chatroom: CreateGroupChatroom is not specified")
		}
		if len(messageData.CreateGroupChatroom.Participants) == 0 {
			return nil, fmt.Errorf("error creating group chatroom: participants should be specified")
		}
		chatroom, err := chatroomService.CreateGroupChatroom(&model.ChatroomOptions{
			IsGroup:      true,
			GroupName:    messageData.CreateGroupChatroom.GroupName,
			Participants: messageData.CreateGroupChatroom.Participants,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating group chatroom: %v", err)
		}
		// update clients' chatroom subscriptions
		for client := range clients {
			if exists(messageData.CreateGroupChatroom.Participants, client.UserID) {
				client.ChatIDs[chatroom.ID] = true
				if err := client.Conn.WriteJSON(messageData); err != nil {
					log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, model.MessageDataOptionSendMessage, err)
					continue
				}
			}
		}
		return messageData, nil
	case model.MessageDataOptionCreatePrivateChatroom:
		log.Printf("creating private chatroom: %v\n", messageData.CreatePrivateChatroom)
		if messageData.CreatePrivateChatroom == nil {
			return nil, fmt.Errorf("error creating private chatroom: CreatePrivateChatroom is not specified")
		}
		firstParticipant := messageData.CreatePrivateChatroom.Participants[0]
		secondParticipant := messageData.CreatePrivateChatroom.Participants[1]
		if firstParticipant.ID == 0 || secondParticipant.ID == 0 {
			return nil, fmt.Errorf("error creating private chatroom: both participants should be specified and have valid IDs")
		}
		if firstParticipant.ID != messageData.CreatePrivateChatroom.ChatMessage.SenderID && secondParticipant.ID != messageData.CreatePrivateChatroom.ChatMessage.SenderID {
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
		messageData.CreatePrivateChatroom.CreatedAt = chatroom.CreatedAt
		// update clients' chatroom subscriptions and send the messageData to the participants
		for client := range clients {
			isUserParticipant := client.UserID == firstParticipant.ID || client.UserID == secondParticipant.ID
			if isUserParticipant {
				client.ChatIDs[chatroom.ID] = true
				// resolve chatroom name and picture for the client because the chatroom name and picture look different for each participant of a private chatroom
				isClientFirstParticipant := client.UserID == firstParticipant.ID
				if isClientFirstParticipant {
					otherParticipant := firstParticipant
					if otherParticipant.Nickname != "" {
						messageData.CreatePrivateChatroom.ChatroomName = otherParticipant.Nickname
					} else {
						messageData.CreatePrivateChatroom.ChatroomName = otherParticipant.Email
					}
					messageData.CreatePrivateChatroom.ChatroomPictureURL = otherParticipant.AvatarURL
				} else {
					otherParticipant := secondParticipant
					if otherParticipant.Nickname != "" {
						messageData.CreatePrivateChatroom.ChatroomName = otherParticipant.Nickname
					} else {
						messageData.CreatePrivateChatroom.ChatroomName = otherParticipant.Email
					}
					messageData.CreatePrivateChatroom.ChatroomName = otherParticipant.Email
					messageData.CreatePrivateChatroom.ChatroomPictureURL = otherParticipant.AvatarURL
				}

				err := client.Conn.WriteJSON(messageData)
				if err != nil {
					log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, model.MessageDataOptionCreatePrivateChatroom, err)
					continue
				}
			}
		}
		messageData.CreatePrivateChatroom.ChatroomName = ""
		messageData.CreatePrivateChatroom.ChatroomPictureURL = ""
		return messageData, nil
	case model.MessageDataOptionUpdateGroupChatroom:
		if messageData.UpdateGroupChatroom == nil {
			return nil, fmt.Errorf("error updating group chatroom: UpdateGroupChatroom is not specified")
		}
		if messageData.UpdateGroupChatroom.ChatroomID == 0 {
			return nil, fmt.Errorf("error updating group chatroom: chatroomID should be specified")
		}
		if messageData.UpdateGroupChatroom.GroupName == "" || len(messageData.UpdateGroupChatroom.Participants) == 0 {
			return nil, fmt.Errorf("error updating group chatroom: nothing to update")
		}
		chatroom, err := chatroomService.UpdateGroupChatroom(messageData.UpdateGroupChatroom)
		if err != nil {
			return nil, fmt.Errorf("error updating group chatroom: %v", err)

		}
		for client := range clients {
			if exists(messageData.UpdateGroupChatroom.Participants, client.UserID) {
				client.ChatIDs[chatroom.ID] = true
				if err := client.Conn.WriteJSON(messageData); err != nil {
					log.Printf("Error sending messageData to client %v with option %v: %v\n", client.UserID, model.MessageDataOptionSendMessage, err)
					continue
				}
			}
		}
	case model.MessageDataOptionDeleteGroupChatroom:
		// TODO: implement chatroom deletion
	case model.MessageDataOptionEditMessage:
		// TODO: implement messageData editing
	case model.MessageDataOptionDeleteMessage:
		// TODO: implement messageData deletion
	case model.MessageDataOptionReactToMessage:
		// TODO: implement messageData reaction
	}
	return nil, fmt.Errorf("unknown messageData option from messageData: %v", messageData)
}

func validateMessageData(messageData *model.MessageData) error {
	if messageData == nil {
		return fmt.Errorf("messageData is nil")
	}
	switch messageData.MessageOption {
	case model.MessageDataOptionSendMessage:
		if messageData.ChatMessage == nil {
			if messageData.ChatroomOptions == nil || !messageData.ChatroomOptions.IsGroup {
				return fmt.Errorf("chatMessage is not specified and either chatroomOptions is not specified or it's not a group chatroom")
			}
		} else if messageData.ChatMessage.ChatroomID == 0 {
			if messageData.ChatroomOptions == nil || len(messageData.ChatroomOptions.Participants) == 0 {
				return fmt.Errorf("for a new chatroom, chatroomOptions and participants should be specified")
			}
		}
	default:
		return fmt.Errorf("unknown messageOption: %v", messageData.MessageOption)
	}

	return nil
}

func exists(slice []uint, item uint) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
