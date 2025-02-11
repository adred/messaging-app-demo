package api

import "github.com/go-chi/chi/v5"

// NewRouter sets up API routes.
func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/messages", handler.SendMessage)
	r.Get("/chats/{chatId}/messages", handler.GetChatMessages)
	r.Get("/users/{userId}/chats", handler.GetUserChats)
	r.Put("/messages/{messageId}/status", handler.UpdateMessageStatus)
	return r
}
