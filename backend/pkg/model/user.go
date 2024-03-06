package model

import "time"

type User struct {
	ID           uint      `json:"id"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	AvatarURL    string    `json:"avatarURL"`
	CreatedAt    time.Time `json:"createdAt"`
}
