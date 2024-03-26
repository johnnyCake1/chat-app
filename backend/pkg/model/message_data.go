package model

// MessageData used to pass message data between client and server through websocket
type MessageData struct {
	// new implementation
	MessageOption string `json:"messageOption,omitempty"`
	// chatroom actions:
	// CreateGroupChatroom is used to create a group chatroom
	CreateGroupChatroom *CreateGroupChatroom `json:"createGroupChatroom,omitempty"`
	// UpdateGroupChatroom is used to update a group chatroom
	UpdateGroupChatroom *UpdateGroupChatroom `json:"updateGroupChatroom,omitempty"`
	// DeleteGroupChatroom is used to delete a group chatroom
	DeleteGroupChatroom *DeleteGroupChatroom `json:"deleteGroupChatroom,omitempty"`
	// CreatePrivateChatroom is used to create a private 1-to-1 chatroom
	CreatePrivateChatroom *CreatePrivateChatroom `json:"createPrivateChatroom,omitempty"`
	// message actions:
	// SendMessage is used to send a message
	SendMessage *SendMessage `json:"sendMessage,omitempty"`
	// ViewMessage is used to mark a message as viewed
	ViewMessage *ViewMessage `json:"viewMessage,omitempty"`
	// EditMessage is used to edit a message
	EditMessage *EditMessage `json:"editMessage,omitempty"`
	// DeleteMessage is used to delete a message
	DeleteMessage *DeleteMessage `json:"deleteMessage,omitempty"`
	// ReactToMessage is used to react to a message
	ReactToMessage *ReactToMessage `json:"reactToMessage,omitempty"`
}

// TODO: currently we are using MessageData for both listening for actions and broadcasting. Instead use MessageData only for actions (requests from client to server) and implement a separate struct for broadcasting notifications (response from server to clients)

type SendMessage struct {
	ChatMessage
}

type ViewMessage struct {
	ViewerID   uint `json:"viewerID,omitempty"`
	MessageID  uint `json:"messageID,omitempty"`
	ChatroomID uint `json:"chatroomID,omitempty"`
	// append on response:
	ChatMessage `json:"chatMessage,omitempty"`
}

type EditMessage struct {
	MessageID uint   `json:"messageID,omitempty"`
	Text      string `json:"text,omitempty"`
}

type DeleteMessage struct {
	MessageID uint `json:"messageID,omitempty"`
}

type ReactToMessage struct {
	MessageID uint   `json:"messageID,omitempty"`
	Reaction  string `json:"reaction,omitempty"` // TODO: Implement reactions
}

type CreatePrivateChatroom struct {
	// Participants is a list of user IDs of the participants in the private chatroom (should be exactly 2 participants)
	Participants [2]User `json:"participants,omitempty"`
	// ChatMessage is the first message to be sent to the chatroom to initialise the private chatroom
	ChatMessage ChatMessage `json:"chatMessage,omitempty"`
	// append on response:
	ChatroomForUser
}

type CreateGroupChatroom struct {
	Chatroom
}

type UpdateGroupChatroom struct {
	Chatroom
}

type DeleteGroupChatroom struct {
	ChatroomID uint `json:"chatroomID,omitempty"`
}

type MesssageOption string

const (
	// MessageDataOptionSendMessage is used to create a chat message or a chatroom (depending on what is provided in the messageData)
	MessageDataOptionSendMessage = "SEND_MESSAGE"
	// MessageDataOptionViewMessage is used to mark a message as viewed
	MessageDataOptionViewMessage = "VIEW_MESSAGE"
	// MessageDataOptionEditMessage is used to edit a message
	MessageDataOptionEditMessage = "EDIT_MESSAGE"
	// MessageDataOptionDeleteMessage is used to delete a message
	MessageDataOptionDeleteMessage = "DELETE_MESSAGE"
	// MessageDataOptionReactToMessage is used to react to a message
	MessageDataOptionReactToMessage = "REACT_TO_MESSAGE"
	// MessageDataOptionCreatePrivateChatroom is used to create a private chatroom
	MessageDataOptionCreatePrivateChatroom = "CREATE_PRIVATE_CHATROOM"
	MessageDataOptionUpdatePrivateChatroom = "UPDATE_PRIVATE_CHATROOM"
	// MessageDataOptionCreateGroupChatroom is used to create a group chatroom
	MessageDataOptionCreateGroupChatroom = "CREATE_GROUP_CHATROOM"
	// MessageDataOptionUpdateGroupChatroom is used to update a group chatroom
	MessageDataOptionUpdateGroupChatroom = "UPDATE_GROUP_CHATROOM"
	// MessageDataOptionDeleteGroupChatroom is used to delete a group chatroom
	MessageDataOptionDeleteGroupChatroom = "DELETE_GROUP_CHATROOM"
)
