package main

import (
	"os"
	"testing"
)

func TestInitializeApp(t *testing.T) {
	// Enable test mode.
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	app, err := InitializeApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	if app.Config == nil {
		t.Fatal("app.Config is nil")
	}
	if app.Router == nil {
		t.Fatal("app.Router is nil")
	}
	if app.RabbitMQ == nil {
		t.Fatal("app.RabbitMQ is nil")
	}
	app.RabbitMQ.Close()
}
