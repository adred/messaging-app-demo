package mq

import (
	"os"
	"testing"
)

func TestNewRabbitMQ(t *testing.T) {
	// Use environment variables for DSN and queue name if provided.
	dsn := os.Getenv("RABBITMQ_DSN")
	queueName := os.Getenv("RABBITMQ_QUEUE")
	if dsn == "" || queueName == "" {
		t.Skip("RABBITMQ_DSN or RABBITMQ_QUEUE not set; skipping RabbitMQ integration test")
	}
	rmq, err := NewRabbitMQ(dsn, queueName)
	if err != nil {
		t.Fatalf("NewRabbitMQ failed: %v", err)
	}
	// Test publish
	err = rmq.PublishMessage([]byte("test message"))
	if err != nil {
		t.Errorf("PublishMessage failed: %v", err)
	}
	rmq.Close()
}
