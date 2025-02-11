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
	Router   http.Handler
	RabbitMQ *mq.RabbitMQ
}

// NewApp constructs an App.
func NewApp(router http.Handler, rabbitMQ *mq.RabbitMQ) *App {
	return &App{
		Router:   router,
		RabbitMQ: rabbitMQ,
	}
}

// ProvideRabbitMQ initializes the RabbitMQ connection.
func ProvideRabbitMQ(cfg *config.Config) (*mq.RabbitMQ, error) {
	var r *mq.RabbitMQ
	var err error
	// Retry for up to 10 attempts with a 5-second delay between tries.
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
		config.LoadConfig,
		ProvideRabbitMQ,
		repository.NewInMemoryMessageRepository,
		repository.NewInMemoryChatRepository,
		application.NewMessageService,
		api.NewHandler,
		api.NewRouter,
		// Bind *chi.Mux (returned by api.NewRouter) to http.Handler.
		wire.Bind(new(http.Handler), new(*chi.Mux)),
		NewApp,
	)
	return &App{}, nil
}
