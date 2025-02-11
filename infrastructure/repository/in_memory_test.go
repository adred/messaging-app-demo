package repository

import (
	"context"
	"testing"
	"time"

	"messaging-app/domain"
)

func TestInMemoryMessageRepository(t *testing.T) {
	repo := NewInMemoryMessageRepository()
	ctx := context.Background()

	msg := &domain.Message{
		ChatID:    1,
		SenderID:  1,
		Content:   "Test message",
		Timestamp: time.Now(),
		Status:    domain.MessageStatusSent,
	}
	createdMsg, err := repo.CreateMessage(ctx, msg)
	if err != nil {
		t.Fatalf("CreateMessage failed: %v", err)
	}
	if createdMsg.ID == 0 {
		t.Error("expected non-zero message ID")
	}

	messages, err := repo.GetMessagesByChatID(ctx, 1)
	if err != nil {
		t.Fatalf("GetMessagesByChatID failed: %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	// Test updating message status.
	err = repo.UpdateMessageStatus(ctx, createdMsg.ID, domain.MessageStatusDelivered)
	if err != nil {
		t.Fatalf("UpdateMessageStatus failed: %v", err)
	}

	// Verify update.
	messages, err = repo.GetMessagesByChatID(ctx, 1)
	if err != nil {
		t.Fatalf("GetMessagesByChatID failed: %v", err)
	}
	if messages[0].Status != domain.MessageStatusDelivered {
		t.Errorf("expected status delivered, got %s", messages[0].Status)
	}
}

func TestInMemoryChatRepository(t *testing.T) {
	repo := NewInMemoryChatRepository()
	ctx := context.Background()

	chat := &domain.Chat{
		Participant1ID: 1,
		Participant2ID: 2,
		Metadata:       "Test chat",
		CreatedAt:      time.Now(),
	}
	createdChat, err := repo.CreateChat(ctx, chat)
	if err != nil {
		t.Fatalf("CreateChat failed: %v", err)
	}
	if createdChat.ID == 0 {
		t.Error("expected non-zero chat ID")
	}

	fetchedChat, err := repo.GetChatByID(ctx, createdChat.ID)
	if err != nil {
		t.Fatalf("GetChatByID failed: %v", err)
	}
	if fetchedChat.Metadata != "Test chat" {
		t.Errorf("expected metadata 'Test chat', got '%s'", fetchedChat.Metadata)
	}

	// Test GetChatsByUserID.
	chats, err := repo.GetChatsByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetChatsByUserID failed: %v", err)
	}
	if len(chats) < 1 {
		t.Errorf("expected at least 1 chat for user 1, got %d", len(chats))
	}
}
