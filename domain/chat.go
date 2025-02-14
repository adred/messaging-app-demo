package domain

import (
	"errors"
	"time"
)

// Chat represents a private conversation between two users.
type Chat struct {
	ID             int64     `json:"id"`
	Participant1ID int64     `json:"participant1Id"`
	Participant2ID int64     `json:"participant2Id"`
	Metadata       string    `json:"metadata"`
	CreatedAt      time.Time `json:"createdAt"`
}

func NewChat(participant1ID, participant2ID int64) (*Chat, error) {
	// Ensure the participants are not the same.
	if participant1ID == participant2ID {
		return nil, errors.New("participants must be different")
	}

	return &Chat{
		Participant1ID: participant1ID,
		Participant2ID: participant2ID,
		Metadata:       "Test chat",
		CreatedAt:      time.Now(),
	}, nil
}
