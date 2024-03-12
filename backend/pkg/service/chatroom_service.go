package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
)

type ChatroomService struct {
	chatroomRepo *repository.ChatroomRepository
}

func NewChatroomService(repo *repository.ChatroomRepository) *ChatroomService {
	return &ChatroomService{chatroomRepo: repo}
}

func (cs *ChatroomService) GetChatroomById(chatroomId uint, page, pageSize int) (*model.Chatroom, error) {
	return cs.chatroomRepo.FindByID(chatroomId, page, pageSize)
}

func (cs *ChatroomService) GetChatroomsByUserId(userId uint, page, pageSize int) ([]model.Chatroom, error) {
	return cs.chatroomRepo.FindChatroomsByUserID(userId, page, pageSize)
}
