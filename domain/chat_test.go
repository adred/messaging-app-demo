package domain

import (
	"testing"
	"time"
)

func TestNewChat_Success(t *testing.T) {
	chat, err := NewChat(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if chat.Participant1ID != 1 {
		t.Errorf("expected Participant1ID to be 1, got: %d", chat.Participant1ID)
	}
	if chat.Participant2ID != 2 {
		t.Errorf("expected Participant2ID to be 2, got: %d", chat.Participant2ID)
	}
	if chat.Metadata != "Test chat" {
		t.Errorf("expected Metadata to be 'Test chat', got: %s", chat.Metadata)
	}
	if time.Since(chat.CreatedAt) > 2*time.Second {
		t.Errorf("expected CreatedAt to be recent, got: %v", chat.CreatedAt)
	}
}

func TestNewChat_SameParticipants(t *testing.T) {
	_, err := NewChat(1, 1)
	if err == nil {
		t.Fatal("expected error when both participants are the same, got nil")
	}
	expectedError := "participants must be different"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
