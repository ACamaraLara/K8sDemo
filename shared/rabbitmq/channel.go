//go:build !skip

package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type RabbitMQChannel struct {
	channel *amqp091.Channel
}

func (c *RabbitMQChannel) Close() error {
	return c.channel.Close()
}

func (c *RabbitMQChannel) Consume(queueName, consumer string, autoAck, exclusive, noLocal, noWait bool,
	args amqp091.Table) (<-chan Message, error) {
	deliveries, err := c.channel.Consume(queueName, consumer, autoAck, exclusive, noLocal, noWait, args)
	if err != nil {
		return nil, err
	}
	msgChannel := make(chan Message)

	// Process deliveries directly and send them to msgChannel.
	for delivery := range deliveries {
		// Wrap the delivery and send it to the msgChannel.
		msgChannel <- NewRabbitMessage(&delivery)
	}

	close(msgChannel)

	return msgChannel, nil
}

func (c *RabbitMQChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	return c.channel.Publish(exchange, key, mandatory, immediate, msg)
}

func (c *RabbitMQChannel) QueueDeclare(queueName string, durable, autoDelete, exclusive, noWait bool,
	args amqp091.Table) (amqp091.Queue, error) {
	return c.channel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, args)
}
