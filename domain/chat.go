package domain

import "time"

// Chat represents a private conversation between two users.
type Chat struct {
	ID             int64     `json:"id"`
	Participant1ID int64     `json:"participant1Id"`
	Participant2ID int64     `json:"participant2Id"`
	Metadata       string    `json:"metadata"`
	CreatedAt      time.Time `json:"createdAt"`
}
