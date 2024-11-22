package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

type MockMessage struct {
	mock.Mock
	Body []byte
}

func (m *MockMessage) Reject(requeue bool) error {
	args := m.Called(requeue)
	return args.Error(0)
}

func (m *MockMessage) Ack(multiple bool) error {
	args := m.Called(multiple)
	return args.Error(0)
}

type MockConnection struct {
	mock.Mock
}

func (conn *MockConnection) Channel() (Channel, error) {
	args := conn.Called()
	return args.Get(0).(Channel), args.Error(1)
}

func (conn *MockConnection) Close() error {
	args := conn.Called()
	return args.Error(0)
}

type MockChannel struct {
	mock.Mock
	Msgs chan MockMessage // Create channel to simulate message consumtion.
}

func (c *MockChannel) Close() error {
	args := c.Called()
	return args.Error(0)
}

func (c *MockChannel) Consume(queueName, consumer string, autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table) (<-chan Message, error) {
	call := c.Called(queueName, consumer, autoAck, exclusive, noLocal, noWait, args)
	return call.Get(0).(<-chan Message), call.Error(1)
}

func (c *MockChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := c.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func (c *MockChannel) QueueDeclare(queueName string, durable, autoDelete, exclusive, noWait bool,
	args amqp.Table) (amqp.Queue, error) {
	retArgs := c.Called(queueName, durable, autoDelete, exclusive, noWait, args)
	return retArgs.Get(0).(amqp.Queue), retArgs.Error(1)
}
