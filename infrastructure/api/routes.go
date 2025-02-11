package api

import (
	"messaging-app/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

// NewRouter sets up API routes.
func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.BasicAuthMiddleware("red", "abc123"))
	r.Use(httprate.LimitByIP(100, time.Minute))

	r.Post("/messages", handler.SendMessage)
	r.Get("/chats/{chatId}/messages", handler.GetChatMessages)
	r.Get("/users/{userId}/chats", handler.GetUserChats)
	r.Put("/messages/{messageId}/status", handler.UpdateMessageStatus)
	return r
}
