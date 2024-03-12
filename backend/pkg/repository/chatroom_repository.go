package repository

import (
	"backend/pkg/model"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

type ChatroomRepository struct {
	db *sql.DB
}

// NewChatroomRepository creates a new instance of ChatroomRepository with the given database connection.
func NewChatroomRepository(db *sql.DB) *ChatroomRepository {
	return &ChatroomRepository{db: db}
}

// FindByID finds a chatroom by its ID. Returns nil if chatroom is not found.
func (r *ChatroomRepository) FindByID(id uint, page, pageSize int) (*model.Chatroom, error) {
	query := "SELECT id, is_group, group_name, created_at FROM chatrooms WHERE id = $1"

	row := r.db.QueryRow(query, id)

	var chatroomID int
	var isGroup bool
	var groupName sql.NullString
	var createdAt time.Time

	err := row.Scan(&chatroomID, &isGroup, &groupName, &createdAt)
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
	chatroom := &model.Chatroom{
		ID:        chatroomID,
		IsGroup:   isGroup,
		GroupName: groupNameString,
		CreatedAt: createdAt,
	}

	// Retrieve messages for this chatroom
	messages, err := r.FindMessagesByChatroomID(uint(chatroomID), page, pageSize)
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

	return chatroom, nil
}

// GetParticipantsForChatroom Retrieves participants for a chatroom
func (r *ChatroomRepository) GetParticipantsForChatroom(chatroomID int) ([]model.User, error) {
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

	// Query to select messages for a chatroom with pagination
	query := `
        SELECT id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted
        FROM messages
        WHERE chatroom_id = $1
        ORDER BY timestamp DESC
        LIMIT $2 OFFSET $3
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

		err := rows.Scan(&message.ID, &message.ChatRoomID, &message.SenderID, &message.Text, &attachmentURLNullable, &message.TimeStamp, &message.Viewed, &message.Deleted)
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

// FindChatroomsByUserID retrieves all the chatrooms that a user belongs to.
func (r *ChatroomRepository) FindChatroomsByUserID(userID uint, page, pageSize int) ([]model.Chatroom, error) {
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

	var chatrooms []model.Chatroom

	for rows.Next() {
		var chatroom model.Chatroom
		var groupNameNullable sql.NullString

		if err := rows.Scan(&chatroom.ID, &chatroom.IsGroup, &groupNameNullable, &chatroom.CreatedAt); err != nil {
			return nil, err
		}

		if groupNameNullable.Valid {
			chatroom.GroupName = groupNameNullable.String
		}

		// Retrieve messages for the chatroom
		messages, err := r.FindMessagesByChatroomID(uint(chatroom.ID), page, pageSize)
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

		// Append chatroom to the slice
		chatrooms = append(chatrooms, chatroom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatrooms, nil
}
