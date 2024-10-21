package rabbitmq

import (
	"github.com/rs/zerolog/log"
)

// QueueExists Checks if a given queue name exists at the broker.
// @param queueName represents the name of the queue that will be checked.
// @return true if exists, false if not.
func (rbMQ *AMQPConn) QueueExists(queueName string) bool {

	// Declaring queue in passive mode checks for existance of a queue. If an error ocurrs
	// it means that the queue doesn't exist.
	_, err := rbMQ.Channel.QueueDeclarePassive(queueName, false, false, false, false, nil)

	if err != nil {
		// Warn about the generated error to check that its not a critical error.
		// (Normaly, the given error will be that the queue doesn't exist)
		log.Warn().Msg(err.Error())
		return false
	}

	return true

}

// Declares a queue. Check QueueDeclare function to know each parameter functionality.
// To take control of what is declared, best practice is to declare Queue Non-Durable and
// Non-Auto-Deleted. This means that the queue will exist while the server that has declared
// it continues running.
// @param queueName Name of the queue that will be created.
// @return error in case queue couldn't be created, nil if creation is success.
func (rbMQ *AMQPConn) DeclareQueue(queueName string, durable, autoDelete, exclusive, noWait bool) error {

	queue, err := rbMQ.RbWrapper.QueueDeclare(rbMQ.Channel, queueName, durable, autoDelete, exclusive, noWait)

	if err != nil {
		log.Error().Msg("Error declaring queue")
		return err
	}

	log.Info().Msg("Declared queue with name " + queue.Name)

	// Add queue to this AMQP Connection.
	rbMQ.Queues = append(rbMQ.Queues, &queue)

	return nil
}
