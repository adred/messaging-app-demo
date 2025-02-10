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
}

type messageService struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	// You can add an MQ publisher here if you want to publish messages asynchronously.
}

// NewMessageService creates a new instance of MessageService.
func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
	}
}

// SendMessage validates and sends a new message.
func (s *messageService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error) {
	// Validate that the chat exists and that the sender is a participant.
	chat, err := s.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, errors.New("sender is not a participant of the chat")
	}

	// Create a new message.
	msg := &domain.Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		Timestamp: time.Now(),
		Status:    domain.MessageStatusSent,
	}
	return s.messageRepo.CreateMessage(ctx, msg)
}

// GetMessages returns all messages for a given chat.
func (s *messageService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	return s.messageRepo.GetMessagesByChatID(ctx, chatID)
}
