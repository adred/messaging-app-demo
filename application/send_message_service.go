package application

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
	// Validate that the sender is one of the hardcoded users.
	if !domain.IsValidUser(senderID) {
		return nil, apistatus.New("invalid sender").UnprocessableEntity()
	}

	// Retrieve the chat; return error if not found.
	chat, as := s.chatRepo.GetChatByID(ctx, chatID)
	if as != nil {
		return nil, as
	}

	// Validate that the sender is part of the chat.
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, apistatus.New("sender is not a participant of the chat").UnprocessableEntity()
	}
	// Ensure attachments slice is non-nil.
	if attachments == nil {
		attachments = []domain.Attachment{}
	}
	msg := &domain.Message{
		ChatID:      chat.ID,
		SenderID:    senderID,
		Content:     content,
		Attachments: attachments,
		Timestamp:   time.Now(),
		Status:      domain.MessageStatusSent,
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
	if userID <= 0 {
		return nil, apistatus.New("invalid userID").UnprocessableEntity()
	}
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
	// Validate the status.
	switch status {
	case domain.MessageStatusSent, domain.MessageStatusDelivered, domain.MessageStatusRead, domain.MessageStatusFailed:
		// valid
	default:
		return apistatus.New("invalid message status").UnprocessableEntity()
	}
	// Check if the message exists.
	_, as := s.messageRepo.GetMessageByID(ctx, messageID)
	if as != nil {
		return as
	}
	return s.messageRepo.UpdateMessageStatus(ctx, messageID, status)
}

func (s *messageService) CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, apistatus.Status) {
	// Validate that both participants are valid.
	if !domain.IsValidUser(participant1ID) || !domain.IsValidUser(participant2ID) {
		return nil, apistatus.New("one or both participants are invalid").UnprocessableEntity()
	}
	// Ensure the participants are not the same.
	if participant1ID == participant2ID {
		return nil, apistatus.New("participants must be different").UnprocessableEntity()
	}

	newChat := &domain.Chat{
		Participant1ID: participant1ID,
		Participant2ID: participant2ID,
		Metadata:       "Created chat",
		CreatedAt:      time.Now(),
	}
	return s.chatRepo.CreateChat(ctx, newChat)
}
