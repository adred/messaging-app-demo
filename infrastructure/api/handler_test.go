package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"messaging-app/domain"
	"messaging-app/pkg/apistatus"

	"github.com/go-chi/chi/v5"
)

// dummyService is a dummy implementation of the MessageService interface for testing.
type dummyService struct{}

// SendMessage now returns an apistatus.Status instead of error.
func (s *dummyService) SendMessage(ctx context.Context, chatID, senderID int64, content string) (*domain.Message, apistatus.Status) {
	// For testing, assume that only chat with ID 1 exists.
	if chatID != 1 {
		return nil, apistatus.New("chat does not exist").NotFound()
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

// GetMessages returns a dummy list of messages.
func (s *dummyService) GetMessages(ctx context.Context, chatID int64) ([]*domain.Message, apistatus.Status) {
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

// ListChatsForUser returns chats for the given user.
func (s *dummyService) ListChatsForUser(ctx context.Context, userID int64) ([]*domain.Chat, apistatus.Status) {
	// Simulate that user with ID 999 has no chats.
	if userID == 999 {
		return nil, apistatus.New("user has no chats").NotFound()
	}
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

// UpdateMessageStatus updates a message's status.
func (s *dummyService) UpdateMessageStatus(ctx context.Context, messageID int64, status domain.MessageStatus) apistatus.Status {
	// If the message ID is not 1, simulate that it doesn't exist.
	if messageID != 1 {
		return apistatus.New("message does not exist").NotFound()
	}
	// Otherwise, simulate success.
	return nil
}

// CreateChat creates a chat if the participants are different.
func (s *dummyService) CreateChat(ctx context.Context, participant1ID, participant2ID int64) (*domain.Chat, apistatus.Status) {
	if participant1ID == participant2ID {
		return nil, apistatus.New("participants must be different").BadRequest()
	}
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

// newChiContext helps set URL parameters in the request context.
func newChiContext(paramKey, paramValue string) *chi.Context {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(paramKey, paramValue)
	return rctx
}

// --- Tests for the API endpoints ---

// TestCreateChat verifies that the CreateChat endpoint returns a valid chat.
func TestCreateChat(t *testing.T) {
	handler := setupTestHandler()

	reqBody := `{"participant1Id": 1, "participant2Id": 2}`
	req := httptest.NewRequest("POST", "/chats", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

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

// TestSendMessage_ValidChat verifies sending a message with a valid chat.
func TestSendMessage_ValidChat(t *testing.T) {
	handler := setupTestHandler()

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

// TestSendMessage_InvalidChat verifies that sending a message with a non-existent chat returns an error.
func TestSendMessage_InvalidChat(t *testing.T) {
	handler := setupTestHandler()

	reqBody := `{"chatId": 999, "senderId": 1, "content": "This should fail."}`
	req := httptest.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.SendMessage(rr, req)

	// Expect an error response (status code not 200).
	if rr.Code == http.StatusOK {
		t.Fatalf("expected error for non-existent chat, but got status 200")
	}
}

// TestGetChatMessages verifies that the GetChatMessages endpoint returns the expected messages.
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

// TestGetUserChats verifies that GetUserChats returns a valid chat for a user.
func TestGetUserChats(t *testing.T) {
	handler := setupTestHandler()

	// Valid case: user 1 has a chat.
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

	// Error case: user 999 has no chats.
	reqNoChats := httptest.NewRequest("GET", "/users/999/chats", nil)
	ctxNoChats := context.WithValue(reqNoChats.Context(), chi.RouteCtxKey, newChiContext("userId", "999"))
	reqNoChats = reqNoChats.WithContext(ctxNoChats)

	rrNoChats := httptest.NewRecorder()
	handler.GetUserChats(rrNoChats, reqNoChats)
	if rrNoChats.Code == http.StatusOK {
		t.Errorf("expected error when user has no chats, got status %d", rrNoChats.Code)
	}
}

// TestUpdateMessageStatus verifies that UpdateMessageStatus updates a message's status properly.
func TestUpdateMessageStatus(t *testing.T) {
	handler := setupTestHandler()

	// Valid update: update message with ID 1.
	updatePayload := `{"status": "delivered"}`
	req := httptest.NewRequest("PUT", "/messages/1/status", bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, newChiContext("messageId", "1"))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateMessageStatus(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status code %d for valid update, got %d", http.StatusNoContent, rr.Code)
	}

	// Error case: non-existent message (ID != 1).
	updatePayload2 := `{"status": "delivered"}`
	req2 := httptest.NewRequest("PUT", "/messages/2/status", bytes.NewBufferString(updatePayload2))
	req2.Header.Set("Content-Type", "application/json")
	ctx2 := context.WithValue(req2.Context(), chi.RouteCtxKey, newChiContext("messageId", "2"))
	req2 = req2.WithContext(ctx2)

	rr2 := httptest.NewRecorder()
	handler.UpdateMessageStatus(rr2, req2)
	if rr2.Code == http.StatusNoContent {
		t.Errorf("expected error when updating non-existent message, but got status %d", rr2.Code)
	}
}
