package api

import (
	"messaging-app/middleware"
	"net/http"
	"time"

	"messaging-app/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

// NewRouter sets up API routes.
func NewRouter(handler *Handler, conf *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.BasicAuthMiddleware(conf.AuthUsername, conf.AuthPassword))
	r.Use(httprate.LimitByIP(conf.RateLimit, time.Minute))

	r.Post("/messages", handler.SendMessage)
	r.Post("/chats", handler.CreateChat)
	r.Get("/chats/{chatId}/messages", handler.GetChatMessages)
	r.Get("/users/{userId}/chats", handler.GetUserChats)
	r.Put("/messages/{messageId}/status", handler.UpdateMessageStatus)

	// Register Swagger/OpenAPI routes without any authentication.
	r.Get("/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/openapi.yaml")
	})
	fs := http.StripPrefix("/openapi/", http.FileServer(http.Dir("./static/swaggerui")))
	r.Get("/openapi/*", fs.ServeHTTP)

	return r
}
