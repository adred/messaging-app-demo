// dummy_rabbitmq.go
package main

// dummyRabbitMQ is a dummy implementation of the RabbitMQ interface for testing.
type dummyRabbitMQ struct{}

func (d *dummyRabbitMQ) PublishMessage(body []byte) error {
	// Immediately succeed.
	return nil
}

func (d *dummyRabbitMQ) Close() {
	// Do nothing.
}
