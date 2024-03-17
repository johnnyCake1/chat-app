package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"fmt"
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

func (cs *ChatroomService) AddMessageToChatroom(chatroomID uint, message model.ChatMessage) (*model.ChatMessage, error) {
	return cs.chatroomRepo.AddMessageToChatroom(chatroomID, message)
}

func (cs *ChatroomService) MarkMessageAsViewed(messageID uint) (*model.ChatMessage, error) {
	return cs.chatroomRepo.MarkMessageAsViewed(messageID)
}

func (cs *ChatroomService) CreatePrivateChatroomWithMessage(chatMessage *model.ChatMessage, chatroomOptions *model.ChatroomOptions) (*model.ChatMessage, error) {
	if chatMessage == nil {
		return nil, fmt.Errorf("cannot create private chatroom: chatMessage is not specified")
	}
	if chatroomOptions == nil {
		return nil, fmt.Errorf("cannot create private chatroom: chatroom options are not specified")
	}
	if chatroomOptions.IsGroup {
		return nil, fmt.Errorf("cannot create private chatroom: chatroom shouldn't be a group")
	}
	if len(chatroomOptions.Participants) != 2 {
		return nil, fmt.Errorf("cannot create private chatroom: chatroom should have 2 participants")
	}

	return cs.chatroomRepo.CreatePrivateChatroomWithMessage(chatroomOptions.Participants[0], chatroomOptions.Participants[1], *chatMessage)
}

func (cs *ChatroomService) CreatePrivateChatroom(createOptions *model.CreatePrivateChatroom) (*model.Chatroom, error) {
	return cs.chatroomRepo.CreatePrivateChatroomWithCreateOptions(createOptions)
}

func (cs *ChatroomService) CreateGroupChatroom(options *model.ChatroomOptions) (*model.Chatroom, error) {
	if !options.IsGroup {
		return nil, fmt.Errorf("cannot create group chatroom: chatroom should be a group")
	}
	if len(options.Participants) < 1 {
		return nil, fmt.Errorf("cannot create group chatroom: chatroom should have at least 1 participant")
	}
	return cs.chatroomRepo.CreateGroupChatroom(options.GroupName, options.Participants)
}

func (cs *ChatroomService) UpdateGroupChatroom(options *model.UpdateGroupChatroom) (*model.Chatroom, error) {
	return cs.chatroomRepo.UpdateGroupChatroom(options)
}
