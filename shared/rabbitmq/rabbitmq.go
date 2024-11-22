package rabbitmq

import (
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type RabbitMQClient struct {
	Conf   *RabbitConfig
	Conn   Connection
	Ch     Channel
	Queues []*amqp.Queue
}

func NewRabbitMQClient(cfg *RabbitConfig) (*RabbitMQClient, error) {
	url := cfg.GetURL()
	log.Info().Msg("Rabbit URL " + url)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error().Msg("Connection cannot be established: " + err.Error())
		return nil, err
	}

	wrappedConn := &RabbitMQConnection{conn: conn}

	ch, err := wrappedConn.Channel()
	if err != nil {
		log.Error().Msg("Error obtaining channel: " + err.Error())
		return nil, err
	}

	log.Info().Msg("Connection to RabbitMQ broker established with address " + url)

	return &RabbitMQClient{
		Conn: wrappedConn,
		Ch:   ch,
		Conf: cfg,
	}, nil

}

// PublishObject sends a message in json format to the broker.
// @throws error with the specific message error. Null otherwise.
func (rbMQ *RabbitMQClient) PublishObject(messageObject interface{}, queueName string, mandatory bool) error {

	jsonBody, err := json.Marshal(messageObject)
	if err != nil {
		return err
	}

	return rbMQ.Ch.Publish(queueName, "", mandatory, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		})

}

// PublishJson sends a pre-formatted JSON message to the broker.
// @throws error Returns a specific error if the JSON is malformed or if there is an error
// during message publishing; nil otherwise.
func (rbMQ *RabbitMQClient) PublishJson(jsonBody []byte, exchangeName string, mandatory bool) error {

	if !json.Valid(jsonBody) {
		return fmt.Errorf("json form broker wrong formatted")
	}

	return rbMQ.Ch.Publish(exchangeName, "", mandatory, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		})
}

// Declares a queue. Check QueueDeclare function to know each parameter functionality.
// To take control of what is declared, generaly practice is to declare Queue Non-Durable and
// Non-Auto-Deleted. This means that the queue will exist while the server that has declared
// it continues running. Better to make it durable if more than one service is publishing at the same queue
// @return error in case queue couldn't be created, nil if creation is success.
func (rbMQ *RabbitMQClient) DeclareQueue(queueName string, durable, autoDelete, exclusive, noWait bool) error {

	queue, err := rbMQ.Ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)

	if err != nil {
		log.Error().Msg("Error declaring queue")
		return err
	}

	log.Info().Msg("Declared queue with name " + queue.Name)

	// Add queue to this AMQP Connection.
	rbMQ.Queues = append(rbMQ.Queues, &queue)

	return nil
}

// ConsumeOnQueue declares a consumer in the broker and listen for messages in the given queue Name.
// @param processMessage callback to the function that the program should execute for the consumed message.
func (rbMQ *RabbitMQClient) ConsumeOnQueue(queueName string, wg *sync.WaitGroup, ack, requeue bool, processMessage func(Message) error) error {
	// Declare Consumer in the broker, only queue name is necessary for a basic consumption of messages.
	msgs, err := rbMQ.Ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error configuring consumer: %s", err)
	}

	// Add a waiting signal to wait group variable. Main program will wait until wg.Done is called.
	wg.Add(1)

	go func() {
		// Free waitGroup variable after exit from the go routine.
		defer wg.Done()

		for msg := range msgs {
			log.Debug().Msg("Mensaje received")
			// Process message with provided function.
			if err := processMessage(msg); err != nil {

				log.Error().Msg("Error processing consumed message, message will be rejected: " + err.Error())

				// Reject message and send it back to the queue if requeue == true.
				if err := msg.Reject(requeue); err != nil {
					log.Error().Msg("Failed to reject message:  " + err.Error())
				}
				continue
			}
			if ack {
				// ACK to the broker. This ACK lets the broker delete the message from the queue.
				// In case there are more consumers of the queue, broker will wait for all of them.
				if err := msg.Ack(false); err != nil {
					log.Error().Msg("Failed to ack message:  " + err.Error())
				}
			}
		}
	}()

	log.Info().Msg("Listening messages on queue " + queueName)

	return nil
}
