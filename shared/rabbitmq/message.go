//go:build !skip

package rabbitmq

import "github.com/rabbitmq/amqp091-go"

// RabbitMessage is a concrete implementation of the Message interface.
type RabbitMessage struct {
	Delivery *amqp091.Delivery
}

// NewRabbitMessage creates a new RabbitMessage.
func NewRabbitMessage(delivery *amqp091.Delivery) *RabbitMessage {
	return &RabbitMessage{Delivery: delivery}
}

// RejectMsg implements the RejectMsg method from the Message interface.
func (msg *RabbitMessage) Reject(requeue bool) error {
	return msg.Delivery.Reject(requeue)
}

// Ack implements the Ack method from the Message interface.
func (msg *RabbitMessage) Ack(multiple bool) error {
	return msg.Delivery.Ack(multiple)
}
