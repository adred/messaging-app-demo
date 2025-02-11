package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables.
	os.Setenv("RABBITMQ_DSN", "amqp://guest:guest@localhost:5672/")
	os.Setenv("RABBITMQ_QUEUE", "test_queue")
	os.Setenv("HTTP_PORT", "3000")
	defer os.Unsetenv("RABBITMQ_DSN")
	defer os.Unsetenv("RABBITMQ_QUEUE")
	defer os.Unsetenv("HTTP_PORT")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.RabbitMQDSN != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("expected RabbitMQDSN 'amqp://guest:guest@localhost:5672/', got '%s'", cfg.RabbitMQDSN)
	}
	if cfg.RabbitMQQueue != "test_queue" {
		t.Errorf("expected RabbitMQ_QUEUE 'test_queue', got '%s'", cfg.RabbitMQQueue)
	}
	if cfg.HTTPPort != "3000" {
		t.Errorf("expected HTTP_PORT '3000', got '%s'", cfg.HTTPPort)
	}
}
