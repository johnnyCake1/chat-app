package model

import (
	"time"
)

type Chatroom struct {
	ID           int
	IsGroup      bool
	GroupName    string
	CreatedAt    time.Time
	Messages     []Message
	Participants []User
}
