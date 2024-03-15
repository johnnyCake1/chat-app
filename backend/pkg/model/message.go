package model

import "time"

type Message struct {
	ID            uint      `json:"id"`
	ChatRoomID    uint      `json:"chatRoomID"`
	SenderID      uint      `json:"senderID"`
	Text          string    `json:"text"`
	AttachmentURL string    `json:"attachmentURL"`
	IsRead        bool      `json:"isRead"`
	TimeStamp     time.Time `json:"timeStamp"`
	Viewed        bool      `json:"viewed"`
	Edited        bool      `json:"edited"`
	Deleted       bool      `json:"deleted"`
}

// MessageData used to pass message data from client through websocket
type MessageData struct {
	Message
	MessageOption string `json:"messageOption"`
}

const (
	MessageOptionSend     = "MESSAGE_SEND"
	MessageOptionView     = "MESSAGE_VIEW"
	MessageOptionEdit     = "MESSAGE_EDIT"
	MessageOptionDelete   = "MESSAGE_DELETE"
	MessageOptionReaction = "MESSAGE_REACTION"
)
