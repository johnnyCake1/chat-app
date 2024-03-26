package model

import "time"

// ChatMessage is the model for chat messages
type ChatMessage struct {
	ID            uint      `json:"id,omitempty"`
	ChatroomID    uint      `json:"chatroomID,omitempty"`
	SenderID      uint      `json:"senderID,omitempty"`
	Text          string    `json:"text,omitempty"`
	AttachmentURL string    `json:"attachmentURL,omitempty"`
	TimeStamp     time.Time `json:"timeStamp,omitempty"`
	Viewed        bool      `json:"viewed"`
	Edited        bool      `json:"edited"`
	Deleted       bool      `json:"deleted"`
}
