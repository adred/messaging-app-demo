package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"messaging-app/domain"

	"github.com/go-chi/chi/v5"
)

// dummyService is a dummy implementation of application.MessageService for testing.
type dummyService struct{}

// SendMessage returns a dummy message if the chat exists (chatID == 1) and returns an error otherwise.
func (s *dummyService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, error) {
	// For testing, assume that only chat with ID 1 exists.
	if chatID != 1 {
		return nil, errors.New("chat does not exist")
	}
	return &domain.Message{
		ID:        1,
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		Timestamp: time.Now().UTC(),
		Status:    domain.MessageStatusSent,
	}, nil
}

// GetMessages returns a dummy list of two messages.
func (s *dummyService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, error) {
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

// ListChatsForUser returns a dummy chat list containing one chat.
func (s *dummyService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, error) {
	return []*domain.Chat{
		{
			ID:             1,
			Participant1ID: userID,
			Participant2ID: 2,
			Metadata:       "Test chat between user " + strconv.FormatInt(userID, 10) + " and user 2",
			CreatedAt:      time.Now(),
		},
	}, nil
}

// UpdateMessageStatus simulates a successful status update.
func (s *dummyService) UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) error {
	return nil
}

// CreateChat simulates creating a chat.
// It returns an error if the two participant IDs are the same.
// Otherwise, it returns a dummy chat with ID 1.
func (s *dummyService) CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, error) {
	if participant1ID == participant2ID {
		return nil, errors.New("participants must be different")
	}
	// For testing, return a dummy chat with ID 1.
	return &domain.Chat{
		ID:             1,
		Participant1ID: participant1ID,
		Participant2ID: participant2ID,
		Metadata:       "Test Chat",
		CreatedAt:      time.Now(),
	}, nil
}

// setupTestHandler creates an API handler using the dummyService.
func setupTestHandler() *Handler {
	svc := &dummyService{}
	return NewHandler(svc)
}

// newChiContext helps to set URL parameters in the request context.
func newChiContext(paramKey, paramValue string) *chi.Context {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(paramKey, paramValue)
	return rctx
}

//
// Tests for the API endpoints
//

func TestCreateChat(t *testing.T) {
	handler := setupTestHandler()

	// Prepare the request payload to create a chat.
	reqBody := `{"participant1Id": 1, "participant2Id": 2}`
	req := httptest.NewRequest("POST", "/chats", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Assume that CreateChat endpoint is defined on handler.
	// Call the CreateChat handler.
	handler.CreateChat(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", rr.Code)
	}

	var chat domain.Chat
	if err := json.NewDecoder(rr.Body).Decode(&chat); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if chat.ID == 0 {
		t.Error("expected non-zero chat ID")
	}
	if chat.Participant1ID != 1 || chat.Participant2ID != 2 {
		t.Errorf("unexpected chat participants: %+v", chat)
	}
}

func TestSendMessage_ValidChat(t *testing.T) {
	handler := setupTestHandler()

	// Prepare the request for sending a message using an existing chat (chatID == 1).
	reqBody := `{"chatId": 1, "senderId": 1, "content": "Hello, test message."}`
	req := httptest.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
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

func TestSendMessage_InvalidChat(t *testing.T) {
	handler := setupTestHandler()

	// Prepare a request for sending a message to a non-existent chat (chatID != 1).
	reqBody := `{"chatId": 999, "senderId": 1, "content": "This should fail."}`
	req := httptest.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.SendMessage(rr, req)

	// We expect an error response (HTTP 400) because the chat does not exist.
	if rr.Code == http.StatusOK {
		t.Fatalf("expected an error for non-existent chat, but got status 200")
	}
}

func TestGetChatMessages(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest("GET", "/chats/1/messages", nil)
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
