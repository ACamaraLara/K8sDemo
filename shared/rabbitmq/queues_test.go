package rabbitmq

import (
	"testing"
)

// Tests if the Queue is stored correctly after create it.
func TestDeclareQueueNotFail(t *testing.T) {

	rbMQ := &AMQPConn{RbWrapper: &RabbitMock{}}

	if err := rbMQ.DeclareQueue("testQueue", false, false, false, false); err != nil {
		t.Fatalf("Error not expected, but one given.: %s", err.Error())
	}

	if len(rbMQ.Queues) != 1 {
		t.Fatal("No queue stored inside the vector")
	}

}
