package model

type Message struct {
	ID            uint
	ChatID        uint
	SenderID      uint
	Text          string
	AttachmentURL string
	IsRead        bool
}
