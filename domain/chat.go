package domain

// Chat represents a private conversation between two users.
type Chat struct {
	ID             int64 `db:"id" json:"id"`
	Participant1ID int64 `db:"participant1_id" json:"participant1Id"`
	Participant2ID int64 `db:"participant2_id" json:"participant2Id"`
}
