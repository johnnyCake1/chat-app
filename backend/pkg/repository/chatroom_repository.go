package repository

import (
	"backend/pkg/model"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type ChatroomRepository struct {
	db *sql.DB
}

// NewChatroomRepository creates a new instance of ChatroomRepository with the given database connection.
func NewChatroomRepository(db *sql.DB) *ChatroomRepository {
	return &ChatroomRepository{db: db}
}

// FindByID finds a chatroom by its ID. Returns nil if chatroom is not found.
func (r *ChatroomRepository) FindByID(chatroomID, userID uint, page, pageSize int) (*model.ChatroomForUser, error) {
	query := "SELECT id, is_group, group_name, created_at FROM chatrooms WHERE id = $1"

	row := r.db.QueryRow(query, chatroomID)

	chatroom := &model.ChatroomForUser{
		Chatroom: model.Chatroom{},
		UserID:   userID,
	}

	var groupName sql.NullString

	err := row.Scan(&chatroom.ID, &chatroom.IsGroup, &groupName, &chatroom.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil if no chatroom found
			return nil, nil
		}
		// Return nil and error if any other error occurs
		return nil, err
	}
	var groupNameString string
	if groupName.Valid {
		groupNameString = groupName.String
	}
	chatroom.GroupName = groupNameString

	// Retrieve messages for this chatroom
	messages, err := r.FindMessagesByChatroomID(chatroomID, page, pageSize)
	if err != nil {
		return nil, err
	}
	chatroom.Messages = messages

	// Retrieve participants for this chatroom
	participants, err := r.GetParticipantsForChatroom(chatroomID)
	if err != nil {
		return nil, err
	}
	chatroom.Participants = participants
	// Calculate unread count for the chatroom
	unreadCount, err := r.CalculateUnreadCount(chatroomID, userID)
	if err != nil {
		return nil, err
	}
	chatroom.UnreadCount = unreadCount

	return chatroom, nil
}

// GetParticipantsForChatroom Retrieves participants for a chatroom
func (r *ChatroomRepository) GetParticipantsForChatroom(chatroomID uint) ([]model.User, error) {
	// Query to select participants for a chatroom
	query := "SELECT u.id, u.nickname, u.email, u.avatar_url, u.created_at FROM users u JOIN chatroom_participants cp ON u.id = cp.user_id WHERE cp.chatroom_id = $1"

	rows, err := r.db.Query(query, chatroomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.User

	// Iterate through the rows and scan participant data into variables
	for rows.Next() {
		var participant model.User
		err := rows.Scan(&participant.ID, &participant.Nickname, &participant.Email, &participant.AvatarURL, &participant.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Append participant to the slice
		participants = append(participants, participant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}

// FindMessagesByChatroomID fetches messages for a chatroom by chatroom ID with pagination support.
func (r *ChatroomRepository) FindMessagesByChatroomID(chatroomID uint, page, pageSize int) ([]model.Message, error) {
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
			LIMIT $2 OFFSET $3
		) AS sub
		ORDER BY timestamp
	`
	rows, err := r.db.Query(query, chatroomID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to store messages
	var messages []model.Message

	// Iterate through the rows and scan message data into variables
	for rows.Next() {
		var message model.Message
		var attachmentURLNullable sql.NullString

		err := rows.Scan(&message.ID, &message.ChatRoomID, &message.SenderID, &message.Text, &attachmentURLNullable, &message.TimeStamp, &message.Viewed, &message.Deleted, &message.Edited)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	return messages, nil
}

// FindChatroomsByUserID retrieves all the chatrooms that a user belongs to and calculates the unread count for each chatroom.
func (r *ChatroomRepository) FindChatroomsByUserID(userID uint, page, pageSize int) ([]model.ChatroomForUser, error) {
	query := `
        SELECT c.id, c.is_group, c.group_name, c.created_at
        FROM chatrooms c
        INNER JOIN chatroom_participants cp ON c.id = cp.chatroom_id
        WHERE cp.user_id = $1
    `
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatrooms []model.ChatroomForUser

	for rows.Next() {
		var chatroom model.ChatroomForUser
		var groupNameNullable sql.NullString

		if err := rows.Scan(&chatroom.ID, &chatroom.IsGroup, &groupNameNullable, &chatroom.CreatedAt); err != nil {
			return nil, err
		}

		if groupNameNullable.Valid {
			chatroom.GroupName = groupNameNullable.String
		}

		// Retrieve messages for the chatroom
		messages, err := r.FindMessagesByChatroomID(chatroom.ID, page, pageSize)
		if err != nil {
			return nil, err
		}
		chatroom.Messages = messages

		// Retrieve participants for the chatroom
		participants, err := r.GetParticipantsForChatroom(chatroom.ID)
		if err != nil {
			return nil, err
		}
		chatroom.Participants = participants

		// Calculate unread count for the chatroom
		unreadCount, err := r.CalculateUnreadCount(chatroom.ID, userID)
		if err != nil {
			return nil, err
		}
		chatroom.UnreadCount = unreadCount
		chatroom.UserID = userID

		// Append chatroom to the slice
		chatrooms = append(chatrooms, chatroom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

func (r *ChatroomRepository) AddMessageToChatroom(chatroomID uint, message model.Message) (model.Message, error) {
	// Prepare SQL query to insert a new message
	query := `
        INSERT INTO messages (chatroom_id, sender_user_id, text, attachment_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited
    `
	var newMessage model.Message
	var attachmentURL sql.NullString
	err := r.db.QueryRow(query, chatroomID, message.SenderID, message.Text, message.AttachmentURL).Scan(
		&newMessage.ID, &newMessage.ChatRoomID, &newMessage.SenderID, &newMessage.Text, &attachmentURL, &newMessage.TimeStamp, &newMessage.Viewed, &newMessage.Deleted, &newMessage.Edited)
	if err != nil {
		return model.Message{}, err
	}
	if attachmentURL.Valid {
		newMessage.AttachmentURL = attachmentURL.String
	}

	return newMessage, nil
}

func (r *ChatroomRepository) MarkMessageAsViewed(messageID uint) (model.Message, error) {
	// Prepare SQL query to update a message and return the updated row
	query := `
        UPDATE messages 
        SET viewed = true 
        WHERE id = $1
        RETURNING id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited
    `
	var message model.Message
	var attachmentURL sql.NullString
	err := r.db.QueryRow(query, messageID).Scan(
		&message.ID, &message.ChatRoomID, &message.SenderID, &message.Text, &attachmentURL, &message.TimeStamp, &message.Viewed, &message.Deleted, &message.Edited)
	if err != nil {
		return model.Message{}, fmt.Errorf("failed to mark the message as viewed: %v", err)
	}
	if attachmentURL.Valid {
		message.AttachmentURL = attachmentURL.String
	}

	return message, nil
}
