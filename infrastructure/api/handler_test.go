package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"messaging-app/domain"
)

// dummyService implements MessageService for testing.
type dummyService struct{}

func (s *dummyService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error) {
	return nil, nil
}
func (s *dummyService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	return nil, nil
}
func (s *dummyService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	return []*domain.Chat{
		{ID: 1, Participant1ID: userID, Participant2ID: 2, Metadata: "Test Chat"},
	}, nil
}

func TestGetUserChats(t *testing.T) {
	svc := &dummyService{}
	handler := NewHandler(svc)
	req := httptest.NewRequest("GET", "/users/1/chats", nil)
	// Set the URL parameter using chi's URLParam helper via the router.
	rr := httptest.NewRecorder()
	r := NewRouter(handler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var chats []*domain.Chat
	if err := json.NewDecoder(rr.Body).Decode(&chats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(chats) != 1 {
		t.Errorf("expected 1 chat, got %d", len(chats))
	}
}
