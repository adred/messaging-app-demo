package domain

import "testing"

func TestHardcodedUsers(t *testing.T) {
	expectedUsers := []User{
		{ID: 1, Name: "Red"},
		{ID: 2, Name: "Jrue"},
		{ID: 3, Name: "Miro"},
		{ID: 4, Name: "Joann"},
	}

	if len(HardcodedUsers) != len(expectedUsers) {
		t.Fatalf("expected %d hardcoded users, got %d", len(expectedUsers), len(HardcodedUsers))
	}

	for i, expected := range expectedUsers {
		actual := HardcodedUsers[i]
		if actual.ID != expected.ID || actual.Name != expected.Name {
			t.Errorf("expected user at index %d to be %+v, got %+v", i, expected, actual)
		}
	}
}

func TestIsValidUser(t *testing.T) {
	// Test valid users.
	validIDs := []int64{1, 2, 3, 4}
	for _, id := range validIDs {
		if !IsValidUser(id) {
			t.Errorf("expected userID %d to be valid", id)
		}
	}

	// Test invalid user IDs.
	invalidIDs := []int64{0, 5, -1, 100}
	for _, id := range invalidIDs {
		if IsValidUser(id) {
			t.Errorf("expected userID %d to be invalid", id)
		}
	}
}
