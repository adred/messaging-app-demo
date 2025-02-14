package domain

// User represents a user in the system.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// HardcodedUsers holds the four recipients.
// This should be replaced with a database lookup.
var HardcodedUsers = []User{
	{ID: 1, Name: "Red"},
	{ID: 2, Name: "Jrue"},
	{ID: 3, Name: "Miro"},
	{ID: 4, Name: "Joann"},
}

// IsValidUser returns true if the provided userID is found in HardcodedUsers.
// This should be replaced with a database lookup.
func IsValidUser(userID int64) bool {
	for _, u := range HardcodedUsers {
		if u.ID == userID {
			return true
		}
	}
	return false
}
