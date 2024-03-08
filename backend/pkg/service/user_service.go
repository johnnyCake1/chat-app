package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) GetAllUsers() ([]model.User, error) {
	return us.repo.FindAll()
}

// GetUsersBySearchTerm finds those users whose nickname or email contains the search term. Returns empty array if no user is found.
func (us *UserService) GetUsersBySearchTerm(searchTerm string) ([]model.User, error) {
	return us.repo.FindUserBySearchTerm(searchTerm)
}
