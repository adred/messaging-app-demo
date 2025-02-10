//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"

	"messaging-app/application"
	"messaging-app/config"
	"messaging-app/infrastructure/api"
	"messaging-app/infrastructure/mq"
	"messaging-app/infrastructure/repository"
)

// App aggregates the dependencies needed to run the application.
type App struct {
	Router   http.Handler
	DB       *sqlx.DB
	RabbitMQ *mq.RabbitMQ
}

// NewApp is a constructor for App.
func NewApp(router http.Handler, db *sqlx.DB, rabbitMQ *mq.RabbitMQ) *App {
	return &App{
		Router:   router,
		DB:       db,
		RabbitMQ: rabbitMQ,
	}
}

// ProvideDB initializes the database connection.
func ProvideDB(cfg *config.Config) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", cfg.DBDSN)
}

// ProvideRabbitMQ initializes the RabbitMQ connection.
func ProvideRabbitMQ(cfg *config.Config) (*mq.RabbitMQ, error) {
	// Note: Use the correct field names from config.Config.
	return mq.NewRabbitMQ(cfg.RabbitMQDSN, cfg.RabbitMQQueue)
}

// InitializeApp sets up and returns an App with all dependencies injected.
func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRabbitMQ,
		repository.NewMessageRepository,
		repository.NewChatRepository,
		application.NewMessageService,
		api.NewHandler,
		api.NewRouter,
		// Bind *chi.Mux (returned by api.NewRouter) to http.Handler.
		wire.Bind(new(http.Handler), new(*chi.Mux)),
		NewApp,
	)
	return &App{}, nil
}
