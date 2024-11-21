//go:build !skip

package rabbitmq

import "github.com/rabbitmq/amqp091-go"

// Wrapper for amqp.Connection
type RabbitMQConnection struct {
	conn *amqp091.Connection
}

// returns a wrapper for the channel object.
func (c *RabbitMQConnection) Channel() (Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQChannel{channel: ch}, nil
}

// close channel connection.
func (c *RabbitMQConnection) Close() error {
	return c.conn.Close()
}
