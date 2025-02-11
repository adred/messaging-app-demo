// application/send_message_service.go
package application

import (
	"context"
	"errors"
	"time"

	"messaging-app/domain"
	"messaging-app/infrastructure/repository"
)

// MessageService defines the use-case methods.
type MessageService interface {
	SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error)
	ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
}

// NewMessageService creates a new instance of MessageService.
func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
	}
}

func (s *messageService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error) {
	chat, err := s.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	// Ensure the sender is a participant.
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, errors.New("sender is not a participant of the chat")
	}
	msg := &domain.Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		Timestamp: time.Now(),
		Status:    domain.MessageStatusSent,
	}
	return s.messageRepo.CreateMessage(ctx, msg)
}

func (s *messageService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	return s.messageRepo.GetMessagesByChatID(ctx, chatID)
}

func (s *messageService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	return s.chatRepo.GetChatsByUserID(ctx, userID)
}
