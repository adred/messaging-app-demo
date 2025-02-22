package application

import (
	"context"
	"encoding/json"
	"log"

	"messaging-app/domain"
	"messaging-app/infrastructure/mq"
	"messaging-app/infrastructure/repository"
	"messaging-app/pkg/apistatus"
)

type MessageService interface {
	SendMessage(ctx context.Context, chatID, senderID int64, content string, attachments []domain.Attachment) (*domain.Message, apistatus.Status)
	GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, apistatus.Status)
	ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, apistatus.Status)
	UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) apistatus.Status
	CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, apistatus.Status)
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

func (s *messageService) SendMessage(ctx context.Context, chatID, senderID int64, content string, attachments []domain.Attachment) (*domain.Message, apistatus.Status) {
	// Retrieve the chat; return error if not found.
	chat, as := s.chatRepo.GetChatByID(ctx, chatID)
	if as != nil {
		return nil, as
	}

	// Use the domain constructor to create the message.
	msg, err := domain.NewMessage(chat, senderID, content, attachments)
	if err != nil {
		return nil, apistatus.New(err.Error()).UnprocessableEntity()
	}

	createdMsg, as := s.messageRepo.CreateMessage(ctx, msg)
	if as != nil {
		return nil, as
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

func (s *messageService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, apistatus.Status) {
	if chatID <= 0 {
		return nil, apistatus.New("unprocessable entity: invalid chatID").UnprocessableEntity()
	}
	// Verify that the chat exists.
	_, as := s.chatRepo.GetChatByID(ctx, chatID)
	if as != nil {
		return nil, as
	}
	return s.messageRepo.GetMessagesByChatID(ctx, chatID)
}

func (s *messageService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, apistatus.Status) {
	if !domain.IsValidUser(userID) {
		return nil, apistatus.New("user does not exist").UnprocessableEntity()
	}
	chats, as := s.chatRepo.GetChatsByUserID(ctx, userID)
	if as != nil {
		return nil, as
	}
	if len(chats) == 0 {
		return nil, apistatus.New("user has no chats").UnprocessableEntity()
	}
	return chats, nil
}

func (s *messageService) UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) apistatus.Status {
	if messageID <= 0 {
		return apistatus.New("invalid messageID").UnprocessableEntity()
	}
	// Check if the message exists.
	msg, as := s.messageRepo.GetMessageByID(ctx, messageID)
	if as != nil {
		return as
	}
	// Validate the status.
	if valid := status.IsValid(); !valid {
		return apistatus.New("invalid message status").UnprocessableEntity()
	}
	// Enforce allowed transitions.
	rules := domain.DefaultTransitionRules()
	if !msg.CanTransitionTo(status, rules) {
		return apistatus.New("invalid status transition").UnprocessableEntity()
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
	}(msg)
	return s.messageRepo.UpdateMessageStatus(ctx, messageID, status)
}

func (s *messageService) CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, apistatus.Status) {
	// Validate that both participants are valid.
	if !domain.IsValidUser(participant1ID) || !domain.IsValidUser(participant2ID) {
		return nil, apistatus.New("one or both participants are invalid").UnprocessableEntity()
	}

	newChat, err := domain.NewChat(participant1ID, participant2ID)
	if err != nil {
		return nil, apistatus.New(err.Error()).UnprocessableEntity()
	}
	return s.chatRepo.CreateChat(ctx, newChat)
}
