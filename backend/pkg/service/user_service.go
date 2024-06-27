package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (us *UserService) GetAllUsers() ([]model.User, error) {
	return us.userRepo.FindAll()
}

// GetUserByID finds a user by their ID and returns the user details. Returns nil if user is not found.
func (us *UserService) GetUserByID(userID uint) (*model.User, error) {
	return us.userRepo.FindByID(userID)
}

// GetUserByEmail finds a user by their email address and returns the user details. Returns nil if user is not found.
func (us *UserService) GetUserByEmail(email string) (*model.User, error) {
	return us.userRepo.FindByEmail(email)
}

// GetUsersBySearchTerm finds those users whose nickname or email contains the search term. Returns empty array if no user is found.
func (us *UserService) GetUsersBySearchTerm(searchTerm string, excludedUsers []uint) ([]model.User, error) {
	return us.userRepo.FindUserBySearchTerm(searchTerm, excludedUsers)
}

func (us *UserService) CreateNewUser(user model.User) (model.User, error) {
	return us.userRepo.AddNewUser(user)
}
