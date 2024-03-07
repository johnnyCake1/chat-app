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
	Deleted       bool      `json:"deleted"`
}
