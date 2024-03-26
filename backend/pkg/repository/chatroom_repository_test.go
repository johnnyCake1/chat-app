package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarkMessageAsViewed(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %s", err)
	}
	defer db.Close()

	// Create a new instance of ChatroomRepository
	repo := NewChatroomRepository(db)

	// Prepare the mock database for the expected query
	timestamp := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO message_views").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE chatroom_participants SET unread_count = GREATEST\\(0, unread_count - 1\\) WHERE chatroom_id = \\$1 AND user_id = \\$2").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("UPDATE messages SET viewed = true WHERE id = \\$1 RETURNING id, chatroom_id, sender_user_id, text, attachment_url, timestamp, viewed, deleted, edited").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "chatroom_id", "sender_user_id", "text", "attachment_url", "timestamp", "viewed", "deleted", "edited"},
		).AddRow(1, 1, 1, "Hello world!", nil, timestamp, true, false, false))
	mock.ExpectCommit()

	// Call the method and check the result getting converted to a ChatMessage
	message, err := repo.MarkMessageAsViewed(1, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, uint(1), message.ID)
	assert.Equal(t, uint(1), message.ChatroomID)
	assert.Equal(t, uint(1), message.SenderID)
	assert.Equal(t, timestamp, message.TimeStamp)
	assert.Equal(t, "Hello world!", message.Text)
	assert.True(t, message.Viewed)
}
