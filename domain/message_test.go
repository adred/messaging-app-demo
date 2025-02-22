package domain

import (
	"testing"
	"time"
)

// TestMessageStatusIsValid verifies that valid statuses return true and invalid ones return false.
func TestMessageStatusIsValid(t *testing.T) {
	validStatuses := []MessageStatus{
		MessageStatusSent,
		MessageStatusDelivered,
		MessageStatusRead,
		MessageStatusFailed,
	}
	for _, s := range validStatuses {
		if !s.IsValid() {
			t.Errorf("expected status %s to be valid", s)
		}
	}

	invalidStatus := MessageStatus("unknown")
	if invalidStatus.IsValid() {
		t.Errorf("expected status %s to be invalid", invalidStatus)
	}
}

// dummyChat is used for testing NewMessage.
var dummyChat = &Chat{
	ID:             1,
	Participant1ID: 1,
	Participant2ID: 2,
	Metadata:       "Test Chat",
	CreatedAt:      time.Now(),
}

// TestNewMessage verifies that NewMessage enforces invariants.
func TestNewMessage(t *testing.T) {
	// Valid case: sender 1 is a participant.
	msg, err := NewMessage(dummyChat, 1, "Hello", nil)
	if err != nil {
		t.Fatalf("expected valid message creation, got error: %v", err)
	}
	if msg.Attachments == nil || len(msg.Attachments) != 0 {
		t.Errorf("expected attachments to be an empty slice, got %v", msg.Attachments)
	}
	if msg.Status != MessageStatusSent {
		t.Errorf("expected initial status 'sent', got %s", msg.Status)
	}

	// Invalid sender: sender 5 is not in HardcodedUsers.
	_, err = NewMessage(dummyChat, 5, "Hello", nil)
	if err == nil || err.Error() != "invalid sender" {
		t.Errorf("expected error 'invalid sender', got %v", err)
	}

	// Sender not a participant: sender 3 is valid but not part of dummyChat.
	_, err = NewMessage(dummyChat, 3, "Hello", nil)
	if err == nil || err.Error() != "sender is not a participant of the chat" {
		t.Errorf("expected error 'sender is not a participant of the chat', got %v", err)
	}
}

// TestCanTransitionTo verifies allowed status transitions.
func TestCanTransitionTo(t *testing.T) {
	// Create a dummy message in "sent" state.
	msg := &Message{
		ID:        1,
		ChatID:    dummyChat.ID,
		SenderID:  1,
		Content:   "Test",
		Timestamp: time.Now(),
		Status:    MessageStatusSent,
	}
	rules := DefaultTransitionRules()
	// From "sent", allowed transitions are "delivered" or "failed".
	if !msg.CanTransitionTo(MessageStatusDelivered, rules) {
		t.Errorf("expected transition from sent to delivered to be allowed")
	}
	if !msg.CanTransitionTo(MessageStatusFailed, rules) {
		t.Errorf("expected transition from sent to failed to be allowed")
	}
	if msg.CanTransitionTo(MessageStatusRead, rules) {
		t.Errorf("expected transition from sent to read to be disallowed")
	}

	// Change status to "delivered" and test.
	msg.Status = MessageStatusDelivered
	if !msg.CanTransitionTo(MessageStatusRead, rules) {
		t.Errorf("expected transition from delivered to read to be allowed")
	}
	if msg.CanTransitionTo(MessageStatusSent, rules) {
		t.Errorf("expected transition from delivered to sent to be disallowed")
	}

	// Once status is "read", no transitions are allowed.
	msg.Status = MessageStatusRead
	if msg.CanTransitionTo(MessageStatusDelivered, rules) || msg.CanTransitionTo(MessageStatusFailed, rules) {
		t.Errorf("expected no transitions allowed from read state")
	}

	// Similarly for "failed".
	msg.Status = MessageStatusFailed
	if msg.CanTransitionTo(MessageStatusDelivered, rules) || msg.CanTransitionTo(MessageStatusRead, rules) {
		t.Errorf("expected no transitions allowed from failed state")
	}
}
