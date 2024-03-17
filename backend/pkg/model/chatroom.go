package model

import (
	"time"
)

type Chatroom struct {
	ID           uint          `json:"id,omitempty"`
	IsGroup      bool          `json:"isGroup,omitempty"`
	GroupName    string        `json:"groupName,omitempty"`
	CreatedAt    time.Time     `json:"createdAt,omitempty"`
	Messages     []ChatMessage `json:"messages,omitempty"`
	Participants []User        `json:"participants,omitempty"`
}

// ChatroomForUser is the model for chatrooms from the perspective of a user because same chatroom can have different data for different users
type ChatroomForUser struct {
	Chatroom
	UserID             uint   `json:"userID,omitempty"`
	ChatroomName       string `json:"chatroomName,omitempty"`
	ChatroomPictureURL string `json:"chatroomPictureURL,omitempty"`
	UnreadCount        int    `json:"unreadCount,omitempty"`
}
