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

func (cs *ChatroomService) GetChatroomById(chatroomID, userID uint, page, pageSize int) (*model.ChatroomForUser, error) {
	return cs.chatroomRepo.FindByID(chatroomID, userID, page, pageSize)
}

func (cs *ChatroomService) GetChatroomsByUserId(userId uint, page, pageSize int) ([]model.ChatroomForUser, error) {
	return cs.chatroomRepo.FindChatroomsByUserID(userId, page, pageSize)
}

func (cs *ChatroomService) AddMessageToChatroom(chatroomID uint, message model.Message) (model.Message, error) {
	return cs.chatroomRepo.AddMessageToChatroom(chatroomID, message)
}

func (cs *ChatroomService) MarkMessageAsViewed(messageID uint) (model.Message, error) {
	return cs.chatroomRepo.MarkMessageAsViewed(messageID)
}
