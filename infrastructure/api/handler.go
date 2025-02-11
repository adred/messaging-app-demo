package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"messaging-app/application"
	"messaging-app/domain"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	messageService application.MessageService
}

func NewHandler(msgService application.MessageService) *Handler {
	return &Handler{messageService: msgService}
}

type SendMessageRequest struct {
	ChatID   int64  `json:"chatId"`
	SenderID int64  `json:"senderId"`
	Content  string `json:"content"`
}

// CreateChatRequest defines the payload to create a chat.
type CreateChatRequest struct {
	Participant1ID int64 `json:"participant1Id"`
	Participant2ID int64 `json:"participant2Id"`
}

// UpdateStatusRequest is the payload for updating a message status.
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// SendMessage handles POST /messages.
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := h.messageService.SendMessage(r.Context(), req.ChatID, req.SenderID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chat, err := h.messageService.CreateChat(r.Context(), req.Participant1ID, req.Participant2ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// GetChatMessages handles GET /chats/{chatId}/messages.
func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	chatIDStr := chi.URLParam(r, "chatId")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid chatId", http.StatusBadRequest)
		return
	}
	messages, err := h.messageService.GetMessages(r.Context(), chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetUserChats handles GET /users/{userId}/chats.
func (h *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}
	chats, err := h.messageService.ListChatsForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

// UpdateMessageStatus handles PUT /messages/{messageId}/status.
func (h *Handler) UpdateMessageStatus(w http.ResponseWriter, r *http.Request) {
	messageIDStr := chi.URLParam(r, "messageId")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newStatus := domain.MessageStatus(req.Status)
	switch newStatus {
	case domain.MessageStatusDelivered, domain.MessageStatusRead, domain.MessageStatusFailed:
		// valid statuses
	default:
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	if err := h.messageService.UpdateMessageStatus(r.Context(), messageID, newStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
