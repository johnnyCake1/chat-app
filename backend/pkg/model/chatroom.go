package model

import (
	"time"
)

type Chatroom struct {
	ID           uint      `json:"id"`
	IsGroup      bool      `json:"isGroup"`
	GroupName    string    `json:"groupName"`
	CreatedAt    time.Time `json:"createdAt"`
	Messages     []Message `json:"messages"`
	Participants []User    `json:"participants"`
}

type ChatroomForUser struct {
	Chatroom
	UserID      uint `json:"userID"`
	UnreadCount int  `json:"unreadCount"`
}
