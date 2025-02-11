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

	// Register Swagger/OpenAPI routes without any authentication.
	r.Get("/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/openapi.yaml")
	})
	fs := http.StripPrefix("/openapi/", http.FileServer(http.Dir("./static/swaggerui")))
	// r.Get("/swagger/*", fs.ServeHTTP)
	r.Get("/openapi/*", fs.ServeHTTP)

	// Create a route group for API endpoints that require basic auth and rate limiting.
	r.Group(func(api chi.Router) {
		api.Use(middleware.BasicAuthMiddleware(conf.AuthUsername, conf.AuthPassword))
		api.Use(httprate.LimitByIP(conf.RateLimit, time.Minute))

		api.Post("/messages", handler.SendMessage)
		api.Get("/chats/{chatId}/messages", handler.GetChatMessages)
		api.Get("/users/{userId}/chats", handler.GetUserChats)
		api.Put("/messages/{messageId}/status", handler.UpdateMessageStatus)
	})

	return r
}
