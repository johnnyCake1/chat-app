package config

import (
	"backend/pkg/repository"
	"github.com/streadway/amqp"
)

type AppDependencies struct {
	Repos             *repository.Repositories // Repositories with database connection
	MessageConnection *amqp.Connection         // RabbitMQ connection
	MessageChannel    *amqp.Channel            // connected RabbitMQ message channel for chat messages
}

const ServerPort = 8080

const ChatMessageQueueName = "chat_messages"

const ChatMessageRoutingKey = ChatMessageQueueName
