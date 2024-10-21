package rabbitmq

import (
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

const (
	MsgTypeJson = "application/json"
)

// ConsumeOnQueue declares a consumer in the broker and listen for messages in the given queue Name.
// @param queueName Name of the queue where the consumer will listen.
// @param msgType Refers to the type message that will be consumed. If a consumed message has a different type
// the consumer will discard it.
// @param numRetries Number of retries that a message is processed before discard it.
// @param msToRetry Time in ms that the consumer will wait until retry the message.
// @param wg Wait group variable to wait for the program finish until the go routine ends.
// @param processMessage callback to the function that the program should execute for the consumed message.
func (rbMQ *AMQPConn) ConsumeOnQueue(queueName string, wg *sync.WaitGroup, processMessage func([]byte) error) error {

	if !rbMQ.checkIsConnected() {
		return fmt.Errorf("service not connected to broker")
	}
	// Declare Consumer in the broker, only queue name is necessary for a basic consumption of messages.
	msgs, err := rbMQ.RbWrapper.Consume(rbMQ.Channel, queueName)
	if err != nil {
		return fmt.Errorf("error configuring consumer: %s", err)
	}

	// Add a waiting signal to wait group variable. Main program will wait until wg.Done is called.
	wg.Add(1)

	go func() {
		// Free waitGroup variable after exit from the go routine.
		defer wg.Done()

		for msg := range msgs {

			if err := rbMQ.checkValidJsonMessage(&msg); err != nil {
				log.Warn().Msg("Format error with received message: " + err.Error())
				// If message hasn't got correct format, reject it an go to the next one.
				if err := rbMQ.RbWrapper.RejectMsg(&msg, false); err != nil {
					log.Error().Msg("Failed to reject message:  " + err.Error())
				}
				continue
			}

			log.Debug().Msg("Mensaje received")
			// Send message to database.
			if err := processMessage(msg.Body); err != nil {

				log.Error().Msg("Error processing consumed message, message will be rejected: " + err.Error())

				if err := rbMQ.RbWrapper.RejectMsg(&msg, false); err != nil {
					log.Error().Msg("Failed to reject message:  " + err.Error())
				}
				continue
			}

			// ACK to the broker. This ACK lets the broker delete the message from the queue.
			// In case there are more consumers of the queue, broker will wait for all of them.
			if err := rbMQ.RbWrapper.AckMsg(&msg, false); err != nil {
				log.Error().Msg("Failed to ack message:  " + err.Error())
			}
		}
	}()

	log.Info().Msg("Listening messages on queue " + queueName)

	return nil
}

func (rbMQ *AMQPConn) checkValidJsonMessage(msg *amqp.Delivery) error {

	if msg.ContentType != MsgTypeJson {
		return fmt.Errorf("this service only supports " + MsgTypeJson + " messages")
	}

	if !json.Valid(msg.Body) {
		return fmt.Errorf("jSon form broker wrong formatted")
	}

	return nil
}
