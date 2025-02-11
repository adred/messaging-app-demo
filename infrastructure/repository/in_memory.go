package repository

import (
	"context"
	"errors"
	"messaging-app/domain"
	"sync"
	"time"
)

// ChatRepository defines methods for chat data.
type ChatRepository interface {
	CreateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error)
	GetChatByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	GetChatsByUserID(ctx context.Context, userID int64) ([]*domain.Chat, error)
}

// MessageRepository defines methods for message data.
type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	GetMessagesByChatID(ctx context.Context, chatID int64) ([]*domain.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) error
}

// InMemoryChatRepository implements ChatRepository in memory.
type InMemoryChatRepository struct {
	chats  map[int64]*domain.Chat
	mu     sync.RWMutex
	nextID int64
}

// NewInMemoryChatRepository creates a new repository and seeds a chat.
func NewInMemoryChatRepository() ChatRepository {
	repo := &InMemoryChatRepository{
		chats:  make(map[int64]*domain.Chat),
		nextID: 1,
	}
	return repo
}

func (r *InMemoryChatRepository) CreateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	chat.ID = r.nextID
	r.nextID++
	chat.CreatedAt = time.Now()
	r.chats[chat.ID] = chat
	return chat, nil
}

func (r *InMemoryChatRepository) GetChatByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	chat, exists := r.chats[chatID]
	if !exists {
		return nil, errors.New("chat not found")
	}
	return chat, nil
}

func (r *InMemoryChatRepository) GetChatsByUserID(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Chat
	for _, chat := range r.chats {
		if chat.Participant1ID == userID || chat.Participant2ID == userID {
			result = append(result, chat)
		}
	}
	return result, nil
}

// InMemoryMessageRepository implements MessageRepository in memory.
type InMemoryMessageRepository struct {
	messages map[int64]*domain.Message
	mu       sync.RWMutex
	nextID   int64
}

func NewInMemoryMessageRepository() MessageRepository {
	return &InMemoryMessageRepository{
		messages: make(map[int64]*domain.Message),
		nextID:   1,
	}
}

func (r *InMemoryMessageRepository) CreateMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg.ID = r.nextID
	r.nextID++
	r.messages[msg.ID] = msg
	return msg, nil
}

func (r *InMemoryMessageRepository) GetMessagesByChatID(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Message
	for _, msg := range r.messages {
		if msg.ChatID == chatID {
			result = append(result, msg)
		}
	}
	return result, nil
}

// UpdateMessageStatus updates the status of an existing message.
func (r *InMemoryMessageRepository) UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg, exists := r.messages[messageID]
	if !exists {
		return errors.New("message not found")
	}
	msg.Status = status
	return nil
}
