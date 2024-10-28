package rabbitmq

import (
	"strconv"

	"github.com/rs/zerolog/log"
)

// Returns a url with the necessary format to connect to RabbitMQ broker.
func (rbMQ *AMQPConn) getURL() string {
	return "amqp://" + rbMQ.User + ":" + rbMQ.Passwd + "@" + rbMQ.Address + ":" + strconv.Itoa(rbMQ.Port) + "/"
}

// InitConnection starts a connection to RabbitMQ broker.
// @param address of the broker to connect.
// @param port port where broker is listenning to new connections.
// @throws error if connection couldn't be established. Nil if connection
// succesded.
func (rbMQ *AMQPConn) InitConnection() error {

	url := rbMQ.getURL()
	log.Info().Msg("Rabbit URL " + url)

	conn, err := rbMQ.RbWrapper.Dial(url)

	if err != nil {
		log.Error().Msg("Connection cannot be established: " + err.Error())
		return err
	}

	rbMQ.Conn = conn

	ch, err := rbMQ.RbWrapper.Channel(rbMQ.Conn)

	if err != nil {
		log.Error().Msg("Error obtaining channel: " + err.Error())
		return err
	}

	// Enable Publisher Confirms on the channel
	err = ch.Confirm(false)
	if err != nil {
		log.Error().Msg("Failed to enable publisher confirms: " + err.Error())
		return err
	}

	rbMQ.Channel = ch

	log.Info().Msg("Connection established with address " + url)

	return nil
}

// CloseConnection closes actual connection to the broker.
// @throws error if there was a problem closing connection.
// Nil otherwise.
func (rbMQ *AMQPConn) CloseConnection() error {

	if err := rbMQ.RbWrapper.CloseChannel(rbMQ.Channel); err != nil {
		log.Error().Msg("Error closing channel: " + err.Error())
		return err
	}

	rbMQ.Channel = nil

	if err := rbMQ.RbWrapper.CloseConnection(rbMQ.Conn); err != nil {
		log.Error().Msg("Error closing connection: " + err.Error())
		return err
	}

	rbMQ.Conn = nil
	log.Debug().Msg("Connection closed.")

	return nil
}

// CheckIsConnected checks if the service connection to the broker is established.
// @return True if the service is connected. False otherwise.
func (rbMQ *AMQPConn) checkIsConnected() bool {
	return rbMQ.Conn != nil
}
