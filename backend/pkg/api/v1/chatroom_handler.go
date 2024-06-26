package v1

import (
	"backend/pkg/config"
	"backend/pkg/service"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetChatroomById gets a chatroom with participants and messages with pagination support
// @Summary Get the chatroom information
// @Description Retrieve information about a chatroom, including participants and messages
// @Tags Chatrooms
// @Accept json
// @Produce json
// @Param id path int true "Chatroom ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} model.Chatroom
// @Failure 400 {object} map[string]string
// @Router /api/v1/chatrooms/{id} [get]
func GetChatroomById(chatroomService *service.ChatroomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chatroomID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid chatroom ID"})
		}
		// Parse pagination parameters and for no valid parameters case replace with default values
		page, err := parseInt(c.Query("page"), 1)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid page query parameter: %v", c.Query("page"))})
		}
		pageSize, err := parseInt(c.Query("pageSize"), config.MessageHistoryPaginationDefaultSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid pageSize query parameter: %v", c.Query("pageSize"))})
		}
		userIDRaw := c.Locals("userID")
		userID, ok := userIDRaw.(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "userID not found in context",
			})
		}

		chatroom, err := chatroomService.GetChatroomById(uint(chatroomID), userID, int(page), int(pageSize))
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
// @Summary Get chatrooms of a user
// @Description Retrieve all chatrooms where a user is a participant
// @Tags Chatrooms
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {array} model.ChatroomForUser
// @Failure 400 {object} map[string]string
// @Router /api/v1/users/{id}/chatrooms [get]
func GetUserChatrooms(chatroomService *service.ChatroomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
		}
		// Parse pagination parameters and for no valid parameters case replace with default values
		page, err := parseInt(c.Query("page"), 1)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid page query parameter: %v", c.Query("page"))})
		}
		pageSize, err := parseInt(c.Query("pageSize"), config.MessageHistoryPaginationDefaultSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid pageSize query parameter: %v", c.Query("pageSize"))})
		}
		chatrooms, err := chatroomService.GetChatroomsByUserId(uint(userId), int(page), int(pageSize))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Couldn't query the chatrooms for user with id %v from database: %v", userId, err)})
		}

		return c.JSON(chatrooms)
	}
}

func parseInt(paramStr string, defaultValue int64) (int64, error) {
	if paramStr == "" {
		return defaultValue, nil
	}
	paramValue, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil || paramValue <= 0 {
		return 0, fmt.Errorf("invalid param string: %v Failed with error: %v", paramStr, err)
	}
	return paramValue, nil
}

func GetChatroomMessages(chatroomService *service.ChatroomService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chatroomID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid chatroom ID"})
		}
		// Parse pagination parameters and for no valid parameters case replace with default values
		page, err := parseInt(c.Query("page"), 1)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid page query parameter: %v", c.Query("page"))})
		}
		pageSize, err := parseInt(c.Query("pageSize"), config.MessageHistoryPaginationDefaultSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Invalid pageSize query parameter: %v", c.Query("pageSize"))})
		}
		messages, err := chatroomService.GetChatroomMessages(uint(chatroomID), int(page), int(pageSize))
		log.Printf("Page number: %v, pageSize: %v, Messages: %v", page, pageSize, messages)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Couldn't query the messages from database: %v", err)})
		}
		if messages == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Messages not found"})
		}

		return c.JSON(messages)
	}
}
