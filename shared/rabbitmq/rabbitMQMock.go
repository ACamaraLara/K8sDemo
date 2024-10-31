package rabbitmq

import (
	rabbitmq "github.com/rabbitmq/amqp091-go"

	"github.com/stretchr/testify/mock"
)

// This struct implements methods of IRabbitWrapper. Simulates calls to rabbitMQ go framework.
type RabbitMock struct {
	mock.Mock
	Msgs chan rabbitmq.Delivery
}

func NewRabbitMock() *RabbitMock { return &RabbitMock{Msgs: make(chan rabbitmq.Delivery, 100)} }

func (rbM *RabbitMock) Dial(url string) (*rabbitmq.Connection, error) {

	return &rabbitmq.Connection{}, nil
}

func (rbM *RabbitMock) Channel(conn *rabbitmq.Connection) (*rabbitmq.Channel, error) {
	return &rabbitmq.Channel{}, nil
}

func (rbM *RabbitMock) CloseChannel(ch *rabbitmq.Channel) error {
	return nil
}

func (rbM *RabbitMock) CloseConnection(conn *rabbitmq.Connection) error {
	return nil
}

func (rbM *RabbitMock) Consume(ch *rabbitmq.Channel, queueName string) (<-chan rabbitmq.Delivery, error) {
	return rbM.Msgs, nil
}

func (rbM *RabbitMock) Publish(ch *rabbitmq.Channel, jsonBody []byte, exchangeName, routingKey string, mandatory bool) error {

	fakeMessage := rabbitmq.Delivery{
		Body:        jsonBody,
		ContentType: exchangeName,
	}
	rbM.Msgs <- fakeMessage

	return nil
}

func (rbM *RabbitMock) QueueDeclare(ch *rabbitmq.Channel, queueName string, durable, autoDelete, exclusive, noWait bool) (rabbitmq.Queue, error) {
	return rabbitmq.Queue{Name: "testQueue"}, nil
}

func (rbM *RabbitMock) RejectMsg(msg *rabbitmq.Delivery, requeue bool) error {
	args := rbM.Called(msg, requeue)
	return args.Error(0)
}

func (rbM *RabbitMock) AckMsg(msg *rabbitmq.Delivery, multiple bool) error {
	args := rbM.Called(msg, multiple)
	return args.Error(0)
}
