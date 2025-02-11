package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize the application using Wire.
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.RabbitMQ.Close()

	srv := &http.Server{
		Addr:         ":" + app.Config.HTTPPort,
		Handler:      app.Router,
		ReadTimeout:  time.Duration(app.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(app.Config.IdleTimeout) * time.Second,
	}

	// Use the HTTP port from the configuration.
	log.Printf("Server starting on port %s", app.Config.HTTPPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
