package repository

import (
	"database/sql"
)

// Repositories contains all the repositories
type Repositories struct {
	UserRepo     *UserRepository
	ChatroomRepo *ChatroomRepository
}

// InitRepositories should be called only once when initialising the app
func InitRepositories(db *sql.DB) *Repositories {
	userRepo := NewUserRepository(db)
	chatroomRepo := NewChatroomRepository(db)
	return &Repositories{
		UserRepo:     userRepo,
		ChatroomRepo: chatroomRepo,
	}
}
