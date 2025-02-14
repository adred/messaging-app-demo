package domain

import (
	"errors"
	"time"
)

// MessageStatus represents the state of a message.
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// IsValid checks if the MessageStatus is one of the allowed statuses.
func (ms MessageStatus) IsValid() bool {
	switch ms {
	case MessageStatusSent, MessageStatusDelivered, MessageStatusRead, MessageStatusFailed:
		return true
	}
	return false
}

// Attachment holds metadata for a file attached to a message.
type Attachment struct {
	ID       int64  `json:"id,omitempty"`
	FileURL  string `json:"fileUrl"`
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

// Message represents a chat message.
type Message struct {
	ID          int64         `json:"id"`
	ChatID      int64         `json:"chatId"`
	SenderID    int64         `json:"senderId"`
	Content     string        `json:"content,omitempty"`
	Attachments []Attachment  `json:"attachments"`
	Timestamp   time.Time     `json:"timestamp"`
	Status      MessageStatus `json:"status"`
}

// AllowedTransitions defines the valid status transitions.
var AllowedTransitions = map[MessageStatus][]MessageStatus{
	MessageStatusSent:      {MessageStatusDelivered, MessageStatusFailed},
	MessageStatusDelivered: {MessageStatusRead},
	// Once a message is "read" or "failed", no transitions are allowed.
}

// NewMessage creates a new Message instance and enforces invariants:
// - The sender must be valid (using IsValidUser).
// - The sender must be a participant of the provided chat.
func NewMessage(chat *Chat, senderID int64, content string, attachments []Attachment) (*Message, error) {
	if !IsValidUser(senderID) {
		return nil, errors.New("invalid sender")
	}
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, errors.New("sender is not a participant of the chat")
	}
	// Ensure attachments is never nil.
	if attachments == nil {
		attachments = []Attachment{}
	}
	return &Message{
		ChatID:      chat.ID,
		SenderID:    senderID,
		Content:     content,
		Attachments: attachments,
		Timestamp:   time.Now(),
		Status:      MessageStatusSent,
	}, nil
}

// CanTransitionTo returns true if the message can transition from its current status to newStatus.
func (m *Message) CanTransitionTo(newStatus MessageStatus) bool {
	allowed, ok := AllowedTransitions[m.Status]
	if !ok {
		// If there's no allowed transition from current status, disallow any change.
		return false
	}
	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}
