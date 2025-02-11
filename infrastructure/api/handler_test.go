package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"messaging-app/domain"

	"github.com/go-chi/chi/v5"
)

// dummyService implements the application.MessageService interface for testing.
type dummyService struct{}

func (s *dummyService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error) {
	// Return a dummy message with a generated ID.
	return &domain.Message{
		ID:        1,
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		Timestamp: time.Now().UTC(),
		Status:    domain.MessageStatusSent,
	}, nil
}

func (s *dummyService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error) {
	// Return a dummy list of two messages.
	return []*domain.Message{
		{
			ID:        1,
			ChatID:    chatID,
			SenderID:  1,
			Content:   "First test message",
			Timestamp: time.Now().UTC(),
			Status:    domain.MessageStatusSent,
		},
		{
			ID:        2,
			ChatID:    chatID,
			SenderID:  2,
			Content:   "Second test message",
			Timestamp: time.Now().UTC(),
			Status:    domain.MessageStatusDelivered,
		},
	}, nil
}

func (s *dummyService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	// Return a dummy chat list containing one chat.
	return []*domain.Chat{
		{
			ID:             1,
			Participant1ID: userID,
			Participant2ID: 2,
			Metadata:       "Test chat between user 1 and user 2",
			CreatedAt:      time.Now(),
		},
	}, nil
}

func (s *dummyService) UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) error {
	// For testing, simply return nil (i.e. success).
	return nil
}

func setupTestHandler() *Handler {
	svc := &dummyService{}
	return NewHandler(svc)
}

func newChiContext(paramKey, paramValue string) *chi.Context {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(paramKey, paramValue)
	return rctx
}

func TestSendMessage(t *testing.T) {
	handler := setupTestHandler()

	// Prepare the request.
	reqBody := `{"chatId": 1, "senderId": 1, "content": "Hello, test message."}`
	req := httptest.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()

	handler.SendMessage(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var msg domain.Message
	if err := json.NewDecoder(rr.Body).Decode(&msg); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if msg.ChatID != 1 || msg.SenderID != 1 || msg.Content != "Hello, test message." {
		t.Errorf("unexpected message returned: %+v", msg)
	}
}

func TestGetChatMessages(t *testing.T) {
	handler := setupTestHandler()

	// Create a request for GET /chats/1/messages.
	req := httptest.NewRequest("GET", "/chats/1/messages", nil)
	// Set the URL parameter using chi's context.
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, newChiContext("chatId", "1"))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetChatMessages(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var messages []*domain.Message
	if err := json.NewDecoder(rr.Body).Decode(&messages); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}
}

func TestGetUserChats(t *testing.T) {
	handler := setupTestHandler()

	// Create a request for GET /users/1/chats.
	req := httptest.NewRequest("GET", "/users/1/chats", nil)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, newChiContext("userId", "1"))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetUserChats(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var chats []*domain.Chat
	if err := json.NewDecoder(rr.Body).Decode(&chats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(chats) != 1 {
		t.Errorf("expected 1 chat, got %d", len(chats))
	}
	if chats[0].Participant1ID != 1 {
		t.Errorf("expected Participant1ID to be 1, got %d", chats[0].Participant1ID)
	}
}

func TestUpdateMessageStatus(t *testing.T) {
	handler := setupTestHandler()

	// Create a request for PUT /messages/1/status.
	updatePayload := `{"status": "delivered"}`
	req := httptest.NewRequest("PUT", "/messages/1/status", bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, newChiContext("messageId", "1"))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateMessageStatus(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}
}
