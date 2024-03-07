package model

import "time"

type User struct {
	ID           uint      `json:"id"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // omit in serialisation
	AvatarURL    string    `json:"avatarURL"`
	CreatedAt    time.Time `json:"createdAt"`
}
