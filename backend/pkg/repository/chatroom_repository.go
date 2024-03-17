package repository

import (
	"backend/pkg/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type ChatroomRepository struct {
	db *sql.DB
}

// NewChatroomRepository creates a new instance of ChatroomRepository with the given database connection.
func NewChatroomRepository(db *sql.DB) *ChatroomRepository {
	return &ChatroomRepository{db: db}
}

// FindByID finds a chatroom by its ID. Returns nil if chatroom is not found.
func (r *ChatroomRepository) FindByID(chatroomID, userID uint, messagesPage, messagesPageSize int) (*model.ChatroomForUser, error) {
	// Calculate offset based on messagesPage number and messagesPage size
	offset := (messagesPage - 1) * messagesPageSize

	query := `
        SELECT c.id, c.is_group, c.group_name, c.created_at,
               m.id, m.chatroom_id, m.sender_user_id, m.text, m.attachment_url, m.timestamp, m.viewed, m.deleted, m.edited,
               u.id, u.nickname, u.email, u.avatar_url, u.created_at
        FROM chatrooms c
        LEFT JOIN (
            SELECT * FROM messages
            WHERE chatroom_id = $1
            ORDER BY timestamp DESC
            LIMIT $2 OFFSET $3
        ) m ON c.id = m.chatroom_id
        LEFT JOIN chatroom_participants cp ON c.id = cp.chatroom_id
        LEFT JOIN users u ON cp.user_id = u.id
        WHERE c.id = $1
    `

	rows, err := r.db.Query(query, chatroomID, messagesPageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find chatroom by id: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		// No rows were returned, so no chatroom with the given ID exists
		return nil, nil
	}

	chatroom := &model.ChatroomForUser{
		Chatroom: model.Chatroom{},
		UserID:   userID,
	}

	var messages []model.ChatMessage
	var participants []model.User

	for rows.Next() {
		var message model.ChatMessage
		var participant model.User
		var groupName sql.NullString
		var attachmentURL sql.NullString

		err := rows.Scan(&chatroom.ID, &chatroom.IsGroup, &groupName, &chatroom.CreatedAt,
			&message.ID, &message.ChatroomID, &message.SenderID, &message.Text, &attachmentURL, &message.TimeStamp, &message.Viewed, &message.Deleted, &message.Edited,
			&participant.ID, &participant.Nickname, &participant.Email, &participant.AvatarURL, &participant.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chatroom data: %v", err)
		}

		if groupName.Valid {
			chatroom.GroupName = groupName.String
		}

		if attachmentURL.Valid {
			message.AttachmentURL = attachmentURL.String
		}

		messages = append(messages, message)
		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to find chatroom by id: %v", err)
	}

	chatroom.Messages = messages
	chatroom.Participants = participants

	return chatroom, nil
}

// GetParticipantsForChatroom Retrieves participants for a chatroom
func (r *ChatroomRepository) GetParticipantsForChatroom(chatroomID uint) ([]model.User, error) {
	// Query to select participants for a chatroom
	query := "SELECT u.id, u.nickname, u.email, u.avatar_url, u.created_at FROM users u JOIN chatroom_participants cp ON u.id = cp.user_id WHERE cp.chatroom_id = $1"

	rows, err := r.db.Query(query, chatroomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants for chatroom: %v", err)
	}
	defer rows.Close()

	var participants []model.User

	// Iterate through the rows and scan participant data into variables
	for rows.Next() {
		var participant model.User
		err := rows.Scan(&participant.ID, &participant.Nickname, &participant.Email, &participant.AvatarURL, &participant.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant data: %v", err)
		}
		// Append participant to the slice
		participants = append(participants, participant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get participants for chatroom: %v", err)
	}

	return participants, nil
}

// FindMessagesByChatroomID fetches messages for a chatroom by chatroom ID with pagination support.
func (r *ChatroomRepository) FindMessagesByChatroomID(chatroomID uint, page, pageSize int) ([]model.ChatMessage, error) {
	// Calculate offset based on page number and page size
	offset := (page - 1) * pageSize

	// Query to select messages for a chatroom with pagination. We first sort messages in desc order and cut the desired part out and sort that part back to ascending order.
	query := `
			SELECT id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited
			FROM (
				SELECT id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited
				FROM messages
				WHERE chatroom_id = $1
				ORDER BY timestamp DESC
				OFFSET $3
				LIMIT $2
				) AS messages
			ORDER BY timestamp ASC
		`
	rows, err := r.db.Query(query, chatroomID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages by chatroom id: %v", err)
	}
	defer rows.Close()

	// Initialize a slice to store messages
	var messages []model.ChatMessage

	// Iterate through the rows and scan message data into variables
	for rows.Next() {
		var message model.ChatMessage
		var attachmentURLNullable sql.NullString

		err := rows.Scan(&message.ID, &message.ChatroomID, &message.SenderID, &message.Text, &attachmentURLNullable, &message.TimeStamp, &message.Viewed, &message.Deleted, &message.Edited)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message data: %v", err)
		}
		var attachmentURLString string
		if attachmentURLNullable.Valid {
			attachmentURLString = attachmentURLNullable.String
		}
		message.AttachmentURL = attachmentURLString
		// Append message to the slice
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to find messages by chatroom id: %v", err)
	}

	return messages, nil
}

// FindChatroomsByUserID retrieves all the chatrooms that a user belongs to and calculates the unread count for each chatroom.
func (r *ChatroomRepository) FindChatroomsByUserID(userID uint, page, pageSize int) ([]model.ChatroomForUser, error) {
	query := `
        SELECT c.id, c.is_group, c.group_name, c.created_at,
               (SELECT COUNT(*) FROM messages WHERE chatroom_id = c.id AND sender_user_id != $1 AND viewed = false) AS unread_count
        FROM chatrooms c
        INNER JOIN chatroom_participants cp ON c.id = cp.chatroom_id
        WHERE cp.user_id = $1
    `
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find chatrooms by user id: %v", err)
	}
	defer rows.Close()

	var chatrooms []model.ChatroomForUser

	for rows.Next() {
		var chatroom model.ChatroomForUser
		var groupNameNullable sql.NullString

		if err := rows.Scan(&chatroom.ID, &chatroom.IsGroup, &groupNameNullable, &chatroom.CreatedAt, &chatroom.UnreadCount); err != nil {
			return nil, fmt.Errorf("failed to scan chatroom data: %v", err)
		}

		if groupNameNullable.Valid {
			chatroom.GroupName = groupNameNullable.String
		}

		// Retrieve messages for the chatroom
		messages, err := r.FindMessagesByChatroomID(chatroom.ID, page, pageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to find chatrooms by user id: %v", err)
		}
		chatroom.Messages = messages

		// Retrieve participants for the chatroom
		participants, err := r.GetParticipantsForChatroom(chatroom.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find chatrooms by user id: %v", err)
		}
		chatroom.Participants = participants
		chatroom.UserID = userID

		// Append chatroom to the slice
		chatrooms = append(chatrooms, chatroom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to find chatrooms by user id: %v", err)
	}

	return chatrooms, nil
}

// CalculateUnreadCount calculates the number of unread messages for a given chatroom and user.
func (r *ChatroomRepository) CalculateUnreadCount(chatroomID uint, userID uint) (int, error) {
	query := `
        SELECT COUNT(*) FROM messages
        WHERE chatroom_id = $1 AND sender_user_id != $2 AND viewed = false
    `
	var unreadCount int
	err := r.db.QueryRow(query, chatroomID, userID).Scan(&unreadCount)
	if err != nil {
		return 0, err
	}
	return unreadCount, nil
}

func (r *ChatroomRepository) AddMessageToChatroom(chatroomID uint, message model.ChatMessage) (*model.ChatMessage, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	newMessage, err := r.AddMessageToChatroomTx(tx, chatroomID, message)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newMessage, nil
}

func (r *ChatroomRepository) AddMessageToChatroomTx(tx *sql.Tx, chatroomID uint, message model.ChatMessage) (model.ChatMessage, error) {
	// Check if chatroom exists
	var exists bool
	err := tx.QueryRow("SELECT EXISTS (SELECT 1 FROM chatrooms WHERE id = $1)", chatroomID).Scan(&exists)
	if err != nil {
		return model.ChatMessage{}, fmt.Errorf("failed to check if chatroom exists: %v", err)
	}
	if !exists {
		return model.ChatMessage{}, fmt.Errorf("chatroom with id %v does not exist", chatroomID)
	}
	query := `
		INSERT INTO messages (chatroom_id, sender_user_id, text, attachment_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited
	`
	var newMessage model.ChatMessage
	var attachmentURL sql.NullString
	err = tx.QueryRow(query, chatroomID, message.SenderID, message.Text, message.AttachmentURL).Scan(
		&newMessage.ID, &newMessage.ChatroomID, &newMessage.SenderID, &newMessage.Text, &attachmentURL, &newMessage.TimeStamp, &newMessage.Viewed, &newMessage.Deleted, &newMessage.Edited)
	if err != nil {
		return model.ChatMessage{}, err
	}
	if attachmentURL.Valid {
		newMessage.AttachmentURL = attachmentURL.String
	}

	return newMessage, nil
}

func (r *ChatroomRepository) MarkMessageAsViewedTx(tx *sql.Tx, messageID uint) (model.ChatMessage, error) {
	query := `
		UPDATE messages
		SET viewed = true
		WHERE id = $1
		RETURNING id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited
	`
	var message model.ChatMessage
	var attachmentURL sql.NullString
	err := tx.QueryRow(query, messageID).Scan(
		&message.ID, &message.ChatroomID, &message.SenderID, &message.Text, &attachmentURL, &message.TimeStamp, &message.Viewed, &message.Deleted, &message.Edited)
	if err != nil {
		return model.ChatMessage{}, fmt.Errorf("failed to mark the message as viewed: %v", err)
	}
	if attachmentURL.Valid {
		message.AttachmentURL = attachmentURL.String
	}

	return message, nil
}

func (r *ChatroomRepository) MarkMessageAsViewed(messageID uint) (*model.ChatMessage, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	message, err := r.MarkMessageAsViewedTx(tx, messageID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", err, err)
		}
		return nil, fmt.Errorf("failed to mark the message as viewed: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// CreatePrivateChatroom Create private chatroom between two users returning the newly created chatroom
func (r *ChatroomRepository) CreatePrivateChatroom(user1ID, user2ID uint) (*model.Chatroom, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	chatroom, err := r.CreatePrivateChatroomTx(tx, user1ID, user2ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}

func (r *ChatroomRepository) CreatePrivateChatroomTx(tx *sql.Tx, user1ID, user2ID uint) (*model.Chatroom, error) {
	query := `
		INSERT INTO chatrooms (is_group)
		VALUES (false)
		RETURNING id, is_group, created_at
	`
	var chatroom model.Chatroom
	err := tx.QueryRow(query).Scan(&chatroom.ID, &chatroom.IsGroup, &chatroom.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Add participants to the chatroom
	if err := r.AddParticipantsToChatroom(tx, chatroom.ID, []uint{user1ID, user2ID}); err != nil {
		return nil, err
	}

	return &chatroom, nil
}

// AddParticipantsToChatroom adds participants to a chatroom in a transaction
func (r *ChatroomRepository) AddParticipantsToChatroom(tx *sql.Tx, id uint, uints []uint) error {
	// Prepare the base of the query
	query := "INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES "

	// Prepare the values placeholder and arguments
	var values []interface{}
	for i, userID := range uints {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("($1, $%d)", i+2) // $1 is for chatroom_id, $2 onwards are for user_id
		values = append(values, userID)
	}

	// Add chatroom_id to the beginning of values
	values = append([]interface{}{id}, values...)

	// Execute the query with the arguments
	_, err := tx.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to add participants to chatroom: %v", err)
	}

	return nil
}

// CreatePrivateChatroomWithMessage creates a private chatroom between two users and adds a message to it.
func (r *ChatroomRepository) CreatePrivateChatroomWithMessage(user1ID, user2ID uint, message model.ChatMessage) (*model.ChatMessage, error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	// Create a private chatroom and add a message to it
	newMessage, err := r.CreatePrivateChatroomWithMessageTx(tx, user1ID, user2ID, message)
	if err != nil {
		// If there is an error, rollback the transaction
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
		}
		return nil, err
	}

	// If everything goes well, commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newMessage, nil
}

// CreatePrivateChatroomWithCreateOptions creates a private chatroom between two users and adds a message to it.
func (r *ChatroomRepository) CreatePrivateChatroomWithCreateOptions(createOptions *model.CreatePrivateChatroom) (*model.Chatroom, error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	// Create a private chatroom and add a message to it
	newChatroom, err := r.CreatePrivateChatroomWithCreateOptionsTx(tx, createOptions)
	if err != nil {
		// If there is an error, rollback the transaction
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
		}
		return nil, err
	}

	// If everything goes well, commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newChatroom, nil
}

func (r *ChatroomRepository) CreatePrivateChatroomWithCreateOptionsTx(tx *sql.Tx, createOptions *model.CreatePrivateChatroom) (model.Chatroom, error) {
	// Create a private chatroom between two users
	chatroom, err := r.CreatePrivateChatroomTx(tx, createOptions.Participants[0].ID, createOptions.Participants[1].ID)
	if err != nil {
		return model.Chatroom{}, fmt.Errorf("failed to create private chatroom: %v", err)
	}

	// Add message to the chatroom
	newMessage, err := r.AddMessageToChatroomTx(
		tx,
		chatroom.ID,
		model.ChatMessage{
			ChatroomID:    chatroom.ID,
			SenderID:      createOptions.ChatMessage.SenderID,
			Text:          createOptions.ChatMessage.Text,
			AttachmentURL: createOptions.ChatMessage.AttachmentURL,
		},
	)
	if err != nil {
		return model.Chatroom{}, err
	}
	// Return the chatroom with the initial message
	chatroom.Messages = append(chatroom.Messages, newMessage)
	return *chatroom, nil
}

func (r *ChatroomRepository) CreatePrivateChatroomWithMessageTx(tx *sql.Tx, user1ID, user2ID uint, message model.ChatMessage) (model.ChatMessage, error) {
	// Create a private chatroom between two users
	chatroom, err := r.CreatePrivateChatroomTx(tx, user1ID, user2ID)
	if err != nil {
		return model.ChatMessage{}, fmt.Errorf("failed to create private chatroom: %v", err)
	}
	log.Printf("got the chatroom: %v", chatroom)

	// Add message to the chatroom
	newMessage, err := r.AddMessageToChatroomTx(tx, chatroom.ID, message)
	if err != nil {
		return model.ChatMessage{}, err
	}

	return newMessage, nil
}

// CreateGroupChatroom creates a group chatroom with the given name and participants.
func (r *ChatroomRepository) CreateGroupChatroom(groupName string, participants []uint) (*model.Chatroom, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	chatroom, err := r.CreateGroupChatroomTx(tx, groupName, participants)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}

func (r *ChatroomRepository) CreateGroupChatroomTx(tx *sql.Tx, groupName string, participants []uint) (*model.Chatroom, error) {
	query := `
		INSERT INTO chatrooms (is_group, group_name)
		VALUES (true, $1)
		RETURNING id, is_group, created_at
	`
	var chatroom model.Chatroom
	err := tx.QueryRow(query, groupName).Scan(&chatroom.ID, &chatroom.IsGroup, &chatroom.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Add participants to the chatroom
	if err := r.AddParticipantsToChatroom(tx, chatroom.ID, participants); err != nil {
		return nil, err
	}

	return &chatroom, nil
}

// UpdateGroupChatroom updates a group chatroom with the given options.
func (r *ChatroomRepository) UpdateGroupChatroom(options *model.UpdateGroupChatroom) (*model.Chatroom, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	chatroom, err := r.UpdateGroupChatroomTx(tx, options)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("rollback failed: %v, original error: %v", rollbackErr, err)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}

func (r *ChatroomRepository) UpdateGroupChatroomTx(tx *sql.Tx, options *model.UpdateGroupChatroom) (*model.Chatroom, error) {
	query := `
		UPDATE chatrooms
		SET group_name = $1
		WHERE id = $2
		RETURNING id, is_group, created_at
	`
	var chatroom model.Chatroom
	err := tx.QueryRow(query, options.GroupName, options.ChatroomID).Scan(&chatroom.ID, &chatroom.IsGroup, &chatroom.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Add participants to the chatroom
	if err := r.AddParticipantsToChatroom(tx, chatroom.ID, options.Participants); err != nil {
		return nil, err
	}

	return &chatroom, nil
}
