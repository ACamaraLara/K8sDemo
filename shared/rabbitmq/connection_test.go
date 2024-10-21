package rabbitmq

import (
	"testing"

	rabbitMocks "github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq/mocks"

	"github.com/rs/zerolog"
)

// Init function executed before start tests to avoid verbose logs.
func init() {
	// Disables logger for unit testing.
	zerolog.SetGlobalLevel(zerolog.Disabled)

}

// Tests that connection func has the expected behavior.
func TestConnectionNotFail(t *testing.T) {

	rbMQ := &AMQPConn{RbWrapper: &rabbitMocks.RabbitMock{}}

	if err := rbMQ.InitConnection(); err != nil {
		t.Error("Expected none error but one given", err)
	}

	if rbMQ.Conn == nil {
		t.Error("Connection shouldn't be a null pointer.")
	}

	if rbMQ.Channel == nil {
		t.Error("Channel shouldn't be a null pointer.")
	}
}

// Tests that close connection func has the expected behavior.
func TestCloseConnectionNotFail(t *testing.T) {

	rbMQ := &AMQPConn{RbWrapper: &rabbitMocks.RabbitMock{}}

	if err := rbMQ.CloseConnection(); err != nil {
		t.Error("Expected none error but one given", err)
	}

	if rbMQ.Conn != nil {
		t.Error("Connection shouldn't be a null pointer.")
	}

	if rbMQ.Channel != nil {
		t.Error("Channel shouldn't be a null pointer.")
	}
}
