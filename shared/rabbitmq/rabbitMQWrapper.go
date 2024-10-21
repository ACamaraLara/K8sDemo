package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitWrapper interface {
	Dial(url string) (*amqp.Connection, error)
	Channel(conn *amqp.Connection) (*amqp.Channel, error)
	CloseChannel(ch *amqp.Channel) error
	CloseConnection(conn *amqp.Connection) error
	Consume(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error)
	Publish(ch *amqp.Channel, jsonBody []byte, exchangeName, routingKey string, mandatory bool) error
	QueueDeclare(ch *amqp.Channel, queueName string, durable, autoDelete, exclusive, noWait bool) (amqp.Queue, error)
	RejectMsg(msg *amqp.Delivery, requeue bool) error
	AckMsg(msg *amqp.Delivery, multiple bool) error
}

type RabbitWrapper struct{}

func (rbW *RabbitWrapper) Dial(url string) (*amqp.Connection, error) {

	return amqp.Dial(url)
}

func (rbW *RabbitWrapper) Channel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

func (rbW *RabbitWrapper) CloseChannel(ch *amqp.Channel) error {
	return ch.Close()
}

func (rbW *RabbitWrapper) CloseConnection(conn *amqp.Connection) error {
	return conn.Close()
}

func (rbW *RabbitWrapper) Consume(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (rbW *RabbitWrapper) Publish(ch *amqp.Channel, jsonBody []byte, exchangeName, routingKey string, mandatory bool) error {
	// Empty context by the moment. A context is used to create signals for the message as a
	// timeout or a reception confirmation signal.

	ctx := context.Background()

	return ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		mandatory,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)
}

func (rbW *RabbitWrapper) QueueDeclare(ch *amqp.Channel, queueName string, durable, autoDelete, exclusive, noWait bool) (amqp.Queue, error) {
	return ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
}

func (rbW *RabbitWrapper) RejectMsg(msg *amqp.Delivery, requeue bool) error {
	return msg.Reject(requeue)
}

func (rbW *RabbitWrapper) AckMsg(msg *amqp.Delivery, multiple bool) error {
	return msg.Ack(multiple)
}
