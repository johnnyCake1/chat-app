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

func (cs *ChatroomService) MarkMessageAsViewed(viewMessage *model.ViewMessage) (*model.ChatMessage, error) {
	message, err := cs.chatroomRepo.MarkMessageAsViewed(viewMessage.ChatroomID, viewMessage.MessageID, viewMessage.ViewerID)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (cs *ChatroomService) CreatePrivateChatroom(createOptions *model.CreatePrivateChatroom) (*model.Chatroom, error) {
	return cs.chatroomRepo.CreatePrivateChatroomWithCreateOptions(createOptions)
}

func (cs *ChatroomService) CreateGroupChatroom(createGroupChatroom *model.CreateGroupChatroom) (*model.Chatroom, error) {
	if !createGroupChatroom.IsGroup {
		return nil, fmt.Errorf("cannot create group chatroom: chatroom should be a group")
	}
	if len(createGroupChatroom.Participants) < 1 {
		return nil, fmt.Errorf("cannot create group chatroom: chatroom should have at least 1 participant")
	}
	if createGroupChatroom.GroupName == "" {
		return nil, fmt.Errorf("cannot create group chatroom: chatroom should have a name")
	}
	userIDs := make([]uint, 0, len(createGroupChatroom.Participants))
	for _, user := range createGroupChatroom.Participants {
		userIDs = append(userIDs, user.ID)
	}
	return cs.chatroomRepo.CreateGroupChatroom(createGroupChatroom.GroupName, userIDs)
}

func (cs *ChatroomService) UpdateGroupChatroom(options *model.UpdateGroupChatroom) (*model.Chatroom, error) {
	return cs.chatroomRepo.UpdateGroupChatroom(options)
}

func (cs *ChatroomService) GetChatroomMessages(chatroomID uint, page, pageSize int) ([]model.ChatMessage, error) {
	return cs.chatroomRepo.GetChatroomMessages(chatroomID, page, pageSize)
}
