package repository

import (
	"context"
	"messaging-app/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	GetMessagesByChatID(ctx context.Context, chatID int64) ([]*domain.Message, error)
}

type messageRepository struct {
	db *sqlx.DB
}

// NewMessageRepository returns a new MessageRepository.
func NewMessageRepository(db *sqlx.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	query := `
		INSERT INTO messages (chat_id, sender_id, content, timestamp, status)
		VALUES (?, ?, ?, ?, ?)
	`
	// Set timestamp if not already set.
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	result, err := r.db.ExecContext(ctx, query, msg.ChatID, msg.SenderID, msg.Content, msg.Timestamp, msg.Status)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	msg.ID = id
	return msg, nil
}

func (r *messageRepository) GetMessagesByChatID(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	var messages []*domain.Message
	query := `
		SELECT id, chat_id, sender_id, content, timestamp, status 
		FROM messages 
		WHERE chat_id = ? 
		ORDER BY timestamp ASC
	`
	err := r.db.SelectContext(ctx, &messages, query, chatID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
