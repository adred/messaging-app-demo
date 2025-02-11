package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"messaging-app/application"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	MessageService application.MessageService
}

func NewHandler(msgService application.MessageService) *Handler {
	return &Handler{MessageService: msgService}
}

type SendMessageRequest struct {
	ChatID   int64  `json:"chatId"`
	SenderID int64  `json:"senderId"`
	Content  string `json:"content"`
}

// SendMessage handles POST /messages.
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := h.MessageService.SendMessage(r.Context(), req.ChatID, req.SenderID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// GetChatMessages handles GET /chats/{chatId}/messages.
func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	chatIDStr := chi.URLParam(r, "chatId")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid chatId", http.StatusBadRequest)
		return
	}
	messages, err := h.MessageService.GetMessages(r.Context(), chatID)
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
	chats, err := h.MessageService.ListChatsForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
