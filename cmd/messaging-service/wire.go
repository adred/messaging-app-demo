//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"

	"messaging-app/application"
	"messaging-app/config"
	"messaging-app/infrastructure/api"
	"messaging-app/infrastructure/mq"
	"messaging-app/infrastructure/repository"
)

// App aggregates the dependencies needed to run the application.
type App struct {
	Config   *config.Config
	Router   http.Handler
	RabbitMQ *mq.RabbitMQ
}

// NewApp is a constructor for App that requires configuration.
func NewApp(cfg *config.Config, router http.Handler, rabbitMQ *mq.RabbitMQ) *App {
	return &App{
		Config:   cfg,
		Router:   router,
		RabbitMQ: rabbitMQ,
	}
}

// ProvideRabbitMQ initializes the RabbitMQ connection.
func ProvideRabbitMQ(cfg *config.Config) (*mq.RabbitMQ, error) {
	var r *mq.RabbitMQ
	var err error
	// Increase to 10 attempts with a 5-second delay between attempts.
	for i := 0; i < 10; i++ {
		r, err = mq.NewRabbitMQ(cfg.RabbitMQDSN, cfg.RabbitMQQueue)
		if err == nil {
			return r, nil
		}
		fmt.Printf("Attempt %d: RabbitMQ not ready: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

// InitializeApp sets up and returns an App with all dependencies injected.
func InitializeApp() (*App, error) {
	wire.Build(
		// Load the configuration from environment variables.
		config.LoadConfig,
		// Provide RabbitMQ using the config.
		ProvideRabbitMQ,
		// In-memory repository implementations.
		repository.NewInMemoryMessageRepository,
		repository.NewInMemoryChatRepository,
		// Application service.
		application.NewMessageService,
		// API handler and router.
		api.NewHandler,
		api.NewRouter,
		// Bind *chi.Mux (returned by api.NewRouter) to http.Handler.
		wire.Bind(new(http.Handler), new(*chi.Mux)),
		// Construct the App using NewApp.
		NewApp,
	)
	return &App{}, nil
}
