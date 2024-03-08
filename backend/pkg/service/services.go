package service

import (
	"backend/pkg/repository"
)

// Services contains all the service structs
type Services struct {
	UserService *UserService
}

// InitServices initialises all the services with given repositories with database connection
func InitServices(repositories *repository.Repositories) *Services {
	userService := NewUserService(repositories.UserRepo)
	return &Services{
		UserService: userService,
	}
}
