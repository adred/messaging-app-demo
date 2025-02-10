package repository

import (
	"context"

	"messaging-app/domain"

	"github.com/jmoiron/sqlx"
)

type ChatRepository interface {
	GetChatByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	// Additional methods like listing chats for a user can be added here.
}

type chatRepository struct {
	db *sqlx.DB
}

// NewChatRepository returns a new ChatRepository.
func NewChatRepository(db *sqlx.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) GetChatByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	var chat domain.Chat
	query := `
		SELECT id, participant1_id, participant2_id 
		FROM chats 
		WHERE id = ?
	`
	if err := r.db.GetContext(ctx, &chat, query, chatID); err != nil {
		return nil, err
	}
	return &chat, nil
}
