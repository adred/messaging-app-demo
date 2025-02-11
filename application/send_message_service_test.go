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

// TestSendMessageAndUpdateStatus tests sending a message and then updating its status.
func TestSendMessageAndUpdateStatus(t *testing.T) {
	// Create in-memory repositories.
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()

	ctx := context.Background()
	// Create a new chat so that SendMessage can succeed.
	// (Note: repository.CreateChat returns a standard error, so no changes are needed here.)
	chat, apistatus := chatRepo.CreateChat(ctx, &domain.Chat{
		Participant1ID: 1,
		Participant2ID: 2,
		Metadata:       "Test Chat",
		CreatedAt:      time.Now(),
	})
	if apistatus != nil {
		t.Fatalf("failed to create chat: %v", apistatus.GetError())
	}

	rabbitMQ := &dummyRabbitMQ{}
	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)

	// Test sending a message.
	msg, apistatus := service.SendMessage(ctx, chat.ID, 1, "Hello from test")
	if apistatus != nil {
		t.Fatalf("SendMessage failed: %s", apistatus.GetMessage())
	}
	if msg.Status != domain.MessageStatusSent {
		t.Errorf("expected message status 'sent', got %s", msg.Status)
	}

	// Test updating the message status.
	apistatus = service.UpdateMessageStatus(ctx, msg.ID, domain.MessageStatusDelivered)
	if apistatus != nil {
		t.Fatalf("UpdateMessageStatus failed: %s", apistatus.GetMessage())
	}

	// Verify update.
	messages, apistatus := msgRepo.GetMessagesByChatID(ctx, chat.ID)
	if apistatus != nil {
		t.Fatalf("GetMessagesByChatID failed: %v", apistatus.GetError())
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

	// Optionally, check that JSON marshalling works.
	_, err := json.Marshal(updatedMsg)
	if err != nil {
		t.Errorf("failed to marshal updated message: %v", err)
	}
}

// TestListChatsForUser tests ListChatsForUser when chats exist.
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

	// Create a dummy message service that wraps the chatRepo.
	service := NewMessageService(nil, chatRepo, &dummyRabbitMQ{})
	chats, apistatus := service.ListChatsForUser(ctx, 1)
	if apistatus != nil {
		t.Fatalf("ListChatsForUser failed: %s", apistatus.GetMessage())
	}

	// Expect exactly 2 chats for user 1.
	expectedChats := 2
	if len(chats) != expectedChats {
		t.Errorf("expected %d chats for user 1, got %d", expectedChats, len(chats))
	}
}

// TestListChatsForUser_NoChats tests that ListChatsForUser returns an error if the user has no chats.
func TestListChatsForUser_NoChats(t *testing.T) {
	chatRepo := repository.NewInMemoryChatRepository()
	ctx := context.Background()

	// No chats are created here.
	service := NewMessageService(nil, chatRepo, &dummyRabbitMQ{})
	_, apistatus := service.ListChatsForUser(ctx, 1)
	if apistatus == nil {
		t.Error("expected error when listing chats for user with no chats, got nil")
	} else {
		expected := "unprocessable entity: user has no chats"
		if apistatus.GetMessage() != expected {
			t.Errorf("expected error %q, got %q", expected, apistatus.GetMessage())
		}
	}
}

// TestUpdateMessageStatus_NonExistent tests that updating a message that does not exist returns an error.
func TestUpdateMessageStatus_NonExistent(t *testing.T) {
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()
	rabbitMQ := &dummyRabbitMQ{}

	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)
	ctx := context.Background()

	// Attempt to update a message with an ID that doesn't exist.
	apistatus := service.UpdateMessageStatus(ctx, 999, domain.MessageStatusDelivered)
	if apistatus == nil {
		t.Error("expected error when updating non-existent message, got nil")
	} else {
		expected := "unprocessable entity: message does not exist"
		if apistatus.GetMessage() != expected {
			t.Errorf("expected error %q, got %q", expected, apistatus.GetMessage())
		}
	}
}

// TestCreateChatAndSendMessage tests the full flow: create a chat then send a message.
func TestCreateChatAndSendMessage(t *testing.T) {
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()
	rabbitMQ := &dummyRabbitMQ{}

	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)
	ctx := context.Background()

	// Create a chat.
	chat, apistatus := service.CreateChat(ctx, 1, 2)
	if apistatus != nil {
		t.Fatalf("CreateChat failed: %s", apistatus.GetMessage())
	}
	if chat.ID == 0 {
		t.Error("expected non-zero chat ID")
	}

	// Send a message using the created chat.
	msg, apistatus := service.SendMessage(ctx, chat.ID, 1, "Test message in created chat")
	if apistatus != nil {
		t.Fatalf("SendMessage failed: %s", apistatus.GetMessage())
	}
	if msg.ChatID != chat.ID {
		t.Errorf("expected chatID %d, got %d", chat.ID, msg.ChatID)
	}
}

// TestSendMessageInvalidChat tests that SendMessage returns an error for a non-existent chat.
func TestSendMessageInvalidChat(t *testing.T) {
	msgRepo := repository.NewInMemoryMessageRepository()
	chatRepo := repository.NewInMemoryChatRepository()
	rabbitMQ := &dummyRabbitMQ{}

	service := NewMessageService(msgRepo, chatRepo, rabbitMQ)
	ctx := context.Background()

	// Attempt to send a message to a non-existent chat (ID 999).
	_, apistatus := service.SendMessage(ctx, 999, 1, "Test message")
	if apistatus == nil {
		t.Error("expected error when sending message to non-existent chat, got nil")
	} else {
		expected := "unprocessable entity: chat does not exist"
		if apistatus.GetMessage() != expected {
			t.Errorf("expected error %q, got %q", expected, apistatus.GetMessage())
		}
	}
}
