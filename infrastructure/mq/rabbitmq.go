package mq

import (
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQInterface defines the methods required by the messaging service.
type RabbitMQInterface interface {
	PublishMessage(body []byte) error
	Close()
}

// RabbitMQ wraps the connection, channel, and queue declaration.
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

// NewRabbitMQ creates and returns a new RabbitMQ instance.
func NewRabbitMQ(connStr, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Queue:      q,
	}, nil
}

// PublishMessage sends a message to the RabbitMQ queue.
func (r *RabbitMQ) PublishMessage(body []byte) error {
	err := r.Channel.Publish(
		"",           // exchange (using default exchange)
		r.Queue.Name, // routing key (queue name)
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}
	return nil
}

// Close cleans up the RabbitMQ connection and channel.
func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Connection.Close()
}
