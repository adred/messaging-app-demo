package application

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"messaging-app/domain"
	"messaging-app/infrastructure/mq"
	"messaging-app/infrastructure/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error)
	GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error)
	ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error)
	UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) error
	CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	rabbitMQ    mq.RabbitMQInterface
}

func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository, rabbitMQ mq.RabbitMQInterface) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		rabbitMQ:    rabbitMQ,
	}
}

func (s *messageService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error) {
	// Validate that the sender is one of the hardcoded users.
	if !domain.IsValidUser(senderID) {
		return nil, errors.New("invalid sender")
	}

	// Retrieve the chat; return error if not found.
	chat, err := s.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, errors.New("chat does not exist")
	}

	// Validate that the sender is a participant of the chat.
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, errors.New("sender is not a participant of the chat")
	}

	// Create the message.
	msg := &domain.Message{
		ChatID:    chat.ID,
		SenderID:  senderID,
		Content:   content,
		Timestamp: time.Now(),
		Status:    domain.MessageStatusSent,
	}
	createdMsg, err := s.messageRepo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Publish asynchronously.
	go func(m *domain.Message) {
		eventData, err := json.Marshal(m)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
			return
		}
		if s.rabbitMQ != nil {
			if err := s.rabbitMQ.PublishMessage(eventData); err != nil {
				log.Printf("failed to publish message event: %v", err)
			}
		}
	}(createdMsg)

	return createdMsg, nil
}

func (s *messageService) CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, error) {
	// Validate that both participants are valid.
	if !domain.IsValidUser(participant1ID) || !domain.IsValidUser(participant2ID) {
		return nil, errors.New("one or both participants are invalid")
	}
	// Ensure the participants are not the same.
	if participant1ID == participant2ID {
		return nil, errors.New("participants must be different")
	}

	newChat := &domain.Chat{
		Participant1ID: participant1ID,
		Participant2ID: participant2ID,
		Metadata:       "Created chat",
		CreatedAt:      time.Now(),
	}
	return s.chatRepo.CreateChat(ctx, newChat)
}

func (s *messageService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	return s.messageRepo.GetMessagesByChatID(ctx, chatID)
}

func (s *messageService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	return s.chatRepo.GetChatsByUserID(ctx, userID)
}

func (s *messageService) UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) error {
	return s.messageRepo.UpdateMessageStatus(ctx, messageID, status)
}
