package domain

// User represents a user in the system.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// HardcodedUsers holds the four recipients.
var HardcodedUsers = []User{
	{ID: 1, Name: "Red"},
	{ID: 2, Name: "Jrue"},
	{ID: 3, Name: "Miro"},
	{ID: 4, Name: "Joann"},
}
