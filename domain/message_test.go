package domain

import (
	"testing"
	"time"
)

func TestMessageStatus(t *testing.T) {
	msg := Message{
		ID:        1,
		ChatID:    1,
		SenderID:  1,
		Content:   "Hello, World!",
		Timestamp: time.Now(),
		Status:    MessageStatusSent,
	}
	if msg.Status != MessageStatusSent {
		t.Errorf("expected status %s but got %s", MessageStatusSent, msg.Status)
	}
}
