// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"fmt"
	"messaging-app/application"
	"messaging-app/config"
	"messaging-app/infrastructure/api"
	"messaging-app/infrastructure/mq"
	"messaging-app/infrastructure/repository"
	"net/http"
	"os"
	"time"
)

// Injectors from wire.go:

// InitializeApp sets up and returns an App with all dependencies injected.
func InitializeApp() (*App, error) {
	configConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	messageRepository := repository.NewInMemoryMessageRepository()
	chatRepository := repository.NewInMemoryChatRepository()
	rabbitMQInterface, err := ProvideRabbitMQ(configConfig)
	if err != nil {
		return nil, err
	}
	messageService := application.NewMessageService(messageRepository, chatRepository, rabbitMQInterface)
	handler := api.NewHandler(messageService)
	mux := api.NewRouter(handler)
	app := NewApp(configConfig, mux, rabbitMQInterface)
	return app, nil
}

// wire.go:

// App aggregates the dependencies needed to run the application.
type App struct {
	Config   *config.Config
	Router   http.Handler
	RabbitMQ mq.RabbitMQInterface
}

// NewApp is a constructor for App that requires configuration.
func NewApp(cfg *config.Config, router http.Handler, rabbitMQ mq.RabbitMQInterface) *App {
	return &App{
		Config:   cfg,
		Router:   router,
		RabbitMQ: rabbitMQ,
	}
}

// ProvideRabbitMQ initializes the RabbitMQ connection.
func ProvideRabbitMQ(cfg *config.Config) (mq.RabbitMQInterface, error) {

	if os.Getenv("TEST_MODE") == "true" {
		return &dummyRabbitMQ{}, nil
	}
	var r mq.RabbitMQInterface
	var err error

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
