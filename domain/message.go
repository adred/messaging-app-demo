package domain

import "time"

// MessageStatus defines the status of a message.
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// Message represents a chat message.
type Message struct {
	ID        int64         `json:"id"`
	ChatID    int64         `json:"chatId"`
	SenderID  int64         `json:"senderId"`
	Content   string        `json:"content"`
	Timestamp time.Time     `json:"timestamp"`
	Status    MessageStatus `json:"status"`
}
