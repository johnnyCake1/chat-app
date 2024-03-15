package config

const ServerPort = 8080

const ChatMessageQueueName = "chat_messages"

const ChatMessageRoutingKey = ChatMessageQueueName

const WebsocketChatSubProtocol = "chat-protocol"

const JwtSecret = "my_secret_key" // TODO: generate a secret key

const MessageHistoryPaginationDefaultSize = 20
