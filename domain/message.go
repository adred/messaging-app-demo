package domain

import "time"

// MessageStatus represents the state of a message.
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

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
