package service

import (
	"backend/pkg/repository"
)

// Services contains all the service structs
type Services struct {
	UserService     *UserService
	ChatroomService *ChatroomService
}

// InitServices initialises all the services with given repositories with database connection
func InitServices(repositories *repository.Repositories) *Services {
	userService := NewUserService(repositories.UserRepo)
	chatroomService := NewChatroomService(repositories.ChatroomRepo)
	return &Services{
		UserService:     userService,
		ChatroomService: chatroomService,
	}
}
