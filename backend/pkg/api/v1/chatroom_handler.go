package v1

import (
	"backend/pkg/config"
	"backend/pkg/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetPrivateChatroom gets a private chatroom with participants and messages
func GetPrivateChatroom(chatroomService *service.ChatroomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chatRoomId, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid chatroom ID"})
		}
		// Parse pagination parameters and for no valid parameters case replace with default values
		var page, pageSize int64
		page, err = strconv.ParseInt(c.Params("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1
		}
		pageSize, err = strconv.ParseInt(c.Params("pageSize"), 10, 64)
		if err != nil || pageSize <= 0 {
			pageSize = config.MessageHistoryPaginationDefaultSize
		}

		chatroom, err := chatroomService.GetChatroomById(uint(chatRoomId), int(page), int(pageSize))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Couldn't query the chatroom from database: %v", err)})
		}
		if chatroom == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Chatroom not found"})
		}

		return c.JSON(chatroom)
	}
}

// GetUserChatrooms gets all chatrooms where given user is a participant
func GetUserChatrooms(chatroomService *service.ChatroomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
		}
		// Parse pagination parameters and for no valid parameters case replace with default values
		var page, pageSize int64
		page, err = strconv.ParseInt(c.Params("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1
		}
		pageSize, err = strconv.ParseInt(c.Params("pageSize"), 10, 64)
		if err != nil || pageSize <= 0 {
			pageSize = config.MessageHistoryPaginationDefaultSize
		}

		chatroom, err := chatroomService.GetChatroomsByUserId(uint(userId), int(page), int(pageSize))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Couldn't query the chatrooms for user with id %v from database: %v", userId, err)})
		}
		if chatroom == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Chatroom not found"})
		}

		return c.JSON(chatroom)
	}
}
