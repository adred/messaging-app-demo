package application

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"messaging-app/domain"
	"messaging-app/infrastructure/repository"
)

// dummyRabbitMQ is a stub implementation of the RabbitMQ interface for testing.
type dummyRabbitMQ struct{}

func (d *dummyRabbitMQ) PublishMessage(body []byte) error {
	// Simply return nil to simulate a successful publish.
	return nil
}

func (d *dummyRabbitMQ) Close() {}

func TestSendMessageAndUpdateStatus(t *testing.T) {
	// Create in-memory repositories.
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()

	// For testing, create a new chat so that SendMessage can succeed.
	ctx := context.Background()
	chat, err := chatRepo.CreateChat(ctx, &domain.Chat{
		Participant1ID: 1,
		Participant2ID: 2,
		Metadata:       "Test Chat",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create chat: %v", err)
	}

	// Use dummyRabbitMQ so that we don’t make a real connection.
	rabbitMQ := &dummyRabbitMQ{}

	// Create the message service.
	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)

	// Test sending a message.
	msg, err := service.SendMessage(ctx, chat.ID, 1, "Hello from test")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if msg.Status != domain.MessageStatusSent {
		t.Errorf("expected message status 'sent', got %s", msg.Status)
	}

	// Test updating the message status.
	err = service.UpdateMessageStatus(ctx, msg.ID, domain.MessageStatusDelivered)
	if err != nil {
		t.Fatalf("UpdateMessageStatus failed: %v", err)
	}

	// Verify update.
	messages, err := msgRepo.GetMessagesByChatID(ctx, chat.ID)
	if err != nil {
		t.Fatalf("GetMessagesByChatID failed: %v", err)
	}
	var updatedMsg *domain.Message
	for _, m := range messages {
		if m.ID == msg.ID {
			updatedMsg = m
			break
		}
	}
	if updatedMsg == nil {
		t.Fatal("message not found after update")
	}
	if updatedMsg.Status != domain.MessageStatusDelivered {
		t.Errorf("expected message status 'delivered', got %s", updatedMsg.Status)
	}
}

func TestListChatsForUser(t *testing.T) {
	chatRepo := repository.NewInMemoryChatRepository()
	ctx := context.Background()

	// Create two chats for user 1.
	_, err := chatRepo.CreateChat(ctx, &domain.Chat{
		Participant1ID: 1,
		Participant2ID: 2,
		Metadata:       "Chat 1",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create chat: %v", err)
	}
	_, err = chatRepo.CreateChat(ctx, &domain.Chat{
		Participant1ID: 3,
		Participant2ID: 1,
		Metadata:       "Chat 2",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create chat: %v", err)
	}

	chats, err := chatRepo.GetChatsByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetChatsByUserID failed: %v", err)
	}

	// Since we no longer pre-seed any chat, we expect exactly 2 chats for user 1.
	expectedChats := 2
	if len(chats) != expectedChats {
		t.Errorf("expected %d chats for user 1, got %d", expectedChats, len(chats))
	}
}

func TestUpdateMessageStatusDirectly(t *testing.T) {
	// Create in-memory repositories.
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()

	// Use dummyRabbitMQ so that we don't make a real connection.
	rabbitMQ := &dummyRabbitMQ{}

	// Create the message service.
	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)
	ctx := context.Background()

	// Create a new chat so that SendMessage can succeed.
	chat, err := chatRepo.CreateChat(ctx, &domain.Chat{
		Participant1ID: 1,
		Participant2ID: 2,
		Metadata:       "Update Status Test Chat",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create chat: %v", err)
	}

	// Create a message.
	msg, err := service.SendMessage(ctx, chat.ID, 1, "Message to update")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	// Check initial status.
	if msg.Status != domain.MessageStatusSent {
		t.Errorf("expected initial status 'sent', got %s", msg.Status)
	}

	// Update message status to "read".
	err = service.UpdateMessageStatus(ctx, msg.ID, domain.MessageStatusRead)
	if err != nil {
		t.Fatalf("UpdateMessageStatus failed: %v", err)
	}

	// Retrieve messages for the chat and verify the updated status.
	messages, err := msgRepo.GetMessagesByChatID(ctx, chat.ID)
	if err != nil {
		t.Fatalf("GetMessagesByChatID failed: %v", err)
	}
	var updatedMsg *domain.Message
	for _, m := range messages {
		if m.ID == msg.ID {
			updatedMsg = m
			break
		}
	}
	if updatedMsg == nil {
		t.Fatal("message not found after update")
	}
	if updatedMsg.Status != domain.MessageStatusRead {
		t.Errorf("expected status 'read', got %s", updatedMsg.Status)
	}

	// Optionally, check that the JSON marshalling of the updated message works as expected.
	_, err = json.Marshal(updatedMsg)
	if err != nil {
		t.Errorf("failed to marshal updated message: %v", err)
	}
}

func TestCreateChatAndSendMessage(t *testing.T) {
	// Create in-memory repositories.
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()
	rabbitMQ := &dummyRabbitMQ{} // Use your dummy RabbitMQ (defined in main or a dedicated file).

	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)
	ctx := context.Background()

	// Create a chat.
	chat, err := service.CreateChat(ctx, 1, 2)
	if err != nil {
		t.Fatalf("CreateChat failed: %v", err)
	}
	if chat.ID == 0 {
		t.Error("expected non-zero chat ID")
	}

	// Send a message using the created chat.
	msg, err := service.SendMessage(ctx, chat.ID, 1, "Test message in created chat")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if msg.ChatID != chat.ID {
		t.Errorf("expected chatID %d, got %d", chat.ID, msg.ChatID)
	}
}

func TestSendMessageInvalidChat(t *testing.T) {
	// Create in-memory repositories.
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()
	rabbitMQ := &dummyRabbitMQ{}

	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)
	ctx := context.Background()

	// Attempt to send a message to a non-existent chat.
	_, err := service.SendMessage(ctx, 999, 1, "Test message")
	if err == nil {
		t.Error("expected error when sending message to non-existent chat")
	}
}
