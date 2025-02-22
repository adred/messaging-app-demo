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

// TransitionRules defines valid status transitions (immutable configuration).
type TransitionRules struct {
	rules map[MessageStatus][]MessageStatus
}

// DefaultTransitionRules provides the standard transition configuration.
func DefaultTransitionRules() TransitionRules {
	return TransitionRules{
		rules: map[MessageStatus][]MessageStatus{
			MessageStatusSent:      {MessageStatusDelivered, MessageStatusFailed},
			MessageStatusDelivered: {MessageStatusRead},
		},
	}
}

// Attachment holds metadata for a file attached to a message.
type Attachment struct {
	ID       int64  `json:"id,omitempty"`
	FileURL  string `json:"fileUrl"`
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

// Message represents an immutable chat message.
type Message struct {
	ID          int64         `json:"id"`
	ChatID      int64         `json:"chatId"`
	SenderID    int64         `json:"senderId"`
	Content     string        `json:"content,omitempty"`
	Attachments []Attachment  `json:"attachments"`
	Timestamp   time.Time     `json:"timestamp"`
	Status      MessageStatus `json:"status"`
}

// NewMessage creates a new Message instance with validation.
func NewMessage(chat *Chat, senderID int64, content string, attachments []Attachment) (*Message, error) {
	if !IsValidUser(senderID) {
		return nil, errors.New("invalid sender")
	}
	if chat.Participant1ID != senderID && chat.Participant2ID != senderID {
		return nil, errors.New("sender is not a participant of the chat")
	}
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

// CanTransitionTo is a pure function that checks status transitions.
// Uses value receiver and explicit transition rules parameter.
func (m Message) CanTransitionTo(newStatus MessageStatus, rules TransitionRules) bool {
	allowed, ok := rules.rules[m.Status]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

// WithStatus creates a new Message instance with updated status.
// Preserves immutability by returning a new copy.
func (m Message) WithStatus(newStatus MessageStatus, rules TransitionRules) (Message, error) {
	if !m.CanTransitionTo(newStatus, rules) {
		return Message{}, errors.New("invalid status transition")
	}

	return Message{
		ID:          m.ID,
		ChatID:      m.ChatID,
		SenderID:    m.SenderID,
		Content:     m.Content,
		Attachments: m.Attachments,
		Timestamp:   m.Timestamp,
		Status:      newStatus,
	}, nil
}
