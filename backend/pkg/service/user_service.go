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
