package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Message interface {
	Reject(requeue bool) error
	Ack(multiple bool) error
}

type Channel interface {
	Close() error
	Consume(queueName string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool,
		args amqp.Table) (<-chan Message, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	QueueDeclare(queueName string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
}

type Connection interface {
	Channel() (Channel, error)
	Close() error
}
