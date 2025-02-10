package domain

import "time"

// MessageStatus defines the possible statuses for a message.
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// Message represents a chat message.
type Message struct {
	ID        int64         `db:"id" json:"id"`
	ChatID    int64         `db:"chat_id" json:"chatId"`
	SenderID  int64         `db:"sender_id" json:"senderId"`
	Content   string        `db:"content" json:"content"`
	Timestamp time.Time     `db:"timestamp" json:"timestamp"`
	Status    MessageStatus `db:"status" json:"status"`
}
