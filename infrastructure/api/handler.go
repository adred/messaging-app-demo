package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"messaging-app/application"
	"messaging-app/config"
	"messaging-app/domain"
	"messaging-app/pkg/apistatus"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	messageService application.MessageService
	config         *config.Config
}

func NewHandler(msgService application.MessageService, cfg *config.Config) *Handler {
	return &Handler{messageService: msgService, config: cfg}
}

type SendMessageRequest struct {
	ChatID      int64               `json:"chatId"`
	SenderID    int64               `json:"senderId"`
	Content     string              `json:"content"`
	Attachments []domain.Attachment `json:"attachments,omitempty"`
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

// processAttachment extracts an attachment from the provided file reader and saves the file to disk.
// It returns a domain.Attachment containing file metadata.
func (h *Handler) processAttachment(file io.Reader, originalFilename, contentType string) (domain.Attachment, error) {
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
	}
	if !allowedTypes[contentType] {
		return domain.Attachment{}, errors.New("invalid file type")
	}
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		return domain.Attachment{}, err
	}
	uniqueName := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + originalFilename
	dstPath := filepath.Join(uploadsDir, uniqueName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return domain.Attachment{}, err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return domain.Attachment{}, err
	}
	stat, err := dst.Stat()
	if err != nil {
		return domain.Attachment{}, err
	}
	att := domain.Attachment{
		FileURL:  "http://localhost:" + h.config.HTTPPort + "/uploads/" + uniqueName,
		FileName: originalFilename,
		MimeType: contentType,
		Size:     stat.Size(),
	}
	return att, nil
}

// SendMessage handles POST /messages.
// It accepts both JSON and multipart/form-data payloads.
// When using multipart/form-data, it extracts attachments and passes them to the service.
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Create a context with a 5-second timeout.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	contentType := r.Header.Get("Content-Type")
	var req SendMessageRequest
	if contentType == "application/json" {
		// Decode JSON payload.
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			as := apistatus.New(err).BadRequest()
			http.Error(w, as.GetMessage(), as.GetStatus())
			return
		}
	} else if len(contentType) >= len("multipart/form-data") &&
		strings.HasPrefix(contentType, "multipart/form-data") {
		// Decode multipart/form-data payload.
		chatIDStr := r.FormValue("chatId")
		senderIDStr := r.FormValue("senderId")
		req.Content = r.FormValue("content")
		var err error
		req.ChatID, err = strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			as := apistatus.New("invalid chatId").BadRequest()
			http.Error(w, as.GetMessage(), as.GetStatus())
			return
		}
		req.SenderID, err = strconv.ParseInt(senderIDStr, 10, 64)
		if err != nil {
			as := apistatus.New("invalid senderId").BadRequest()
			http.Error(w, as.GetMessage(), as.GetStatus())
			return
		}
		// Process attachments.
		files := r.MultipartForm.File["files"]
		for _, fh := range files {
			file, err := fh.Open()
			if err != nil {
				as := apistatus.New("failed to open file").InternalServerError()
				http.Error(w, as.GetMessage(), as.GetStatus())
				return
			}
			att, err := h.processAttachment(file, fh.Filename, fh.Header.Get("Content-Type"))
			file.Close()
			if err != nil {
				as := apistatus.New(err).InternalServerError()
				http.Error(w, as.GetMessage(), as.GetStatus())
				return
			}
			req.Attachments = append(req.Attachments, att)
		}
	} else {
		as := apistatus.New("unsupported content type").BadRequest()
		http.Error(w, as.GetMessage(), as.GetStatus())
		return
	}
	msg, apistatus := h.messageService.SendMessage(ctx, req.ChatID, req.SenderID, req.Content, req.Attachments)
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
