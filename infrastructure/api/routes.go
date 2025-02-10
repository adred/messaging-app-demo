package api

import "github.com/go-chi/chi/v5"

// NewRouter sets up all API routes.
func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/messages", handler.SendMessage)
	r.Get("/chats/{chatId}/messages", handler.GetChatMessages)
	return r
}
