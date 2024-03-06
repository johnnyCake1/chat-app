package v1

import (
	"backend/pkg/config"
	"backend/pkg/model"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

// SendMessage Sends a message to message queue
func SendMessage(messageChannel *amqp.Channel) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(model.Message)
		err := c.BodyParser(body)
		if err != nil {
			_ = c.Status(fiber.StatusBadRequest).SendString(err.Error())
			return err
		}
		err = SendToQueue(*body, messageChannel)
		if err != nil {
			_ = c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			return err
		}
		return c.Status(fiber.StatusOK).JSON(body)
	}
}

func SendToQueue(message model.Message, ch *amqp.Channel) error {
	// Serialize message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",                           // exchange
		config.ChatMessageRoutingKey, // routing key
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
	return err
}
