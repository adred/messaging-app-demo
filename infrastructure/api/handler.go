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
	msg, apistatus := h.messageService.SendMessage(r.Context(), req.ChatID, req.SenderID, req.Content)
	if apistatus != nil {
		http.Error(w, apistatus.GetMessage(), apistatus.GetStatus())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chat, apistatus := h.messageService.CreateChat(r.Context(), req.Participant1ID, req.Participant2ID)
	if apistatus != nil {
		http.Error(w, apistatus.GetMessage(), apistatus.GetStatus())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
	messages, apistatus := h.messageService.GetMessages(r.Context(), chatID)
	if apistatus != nil {
		http.Error(w, apistatus.GetMessage(), apistatus.GetStatus())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	chats, apistatus := h.messageService.ListChatsForUser(r.Context(), userID)
	if apistatus != nil {
		http.Error(w, apistatus.GetMessage(), apistatus.GetStatus())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	apistatus := h.messageService.UpdateMessageStatus(r.Context(), messageID, newStatus)
	if apistatus != nil {
		http.Error(w, apistatus.GetMessage(), apistatus.GetStatus())
		return
	}
	w.WriteHeader(http.StatusOK)
}
