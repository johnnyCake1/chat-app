package model

import (
	"time"
)

type Chatroom struct {
	ID           int       `json:"id"`
	IsGroup      bool      `json:"isGroup"`
	GroupName    string    `json:"groupName"`
	CreatedAt    time.Time `json:"createdAt"`
	Messages     []Message `json:"messages"`
	Participants []User    `json:"participants"`
}
