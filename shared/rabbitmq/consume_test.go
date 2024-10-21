package rabbitmq

import (
	"fmt"
	"sync"
	"testing"
	"time"

	rabbitMocks "github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq/mocks"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Tests that Consume function cannot be called before connection hasn't been established.
func TestConsumeMessageWithoutConnection(t *testing.T) {
	rbMQ := &AMQPConn{}

	expectedError := fmt.Errorf("service not connected to broker")
	testProcessMsg := func([]byte) error {
		return nil
	}

	err := rbMQ.ConsumeOnQueue("testQueue", &sync.WaitGroup{}, testProcessMsg)

	if err.Error() != expectedError.Error() {
		t.Fatalf("Got a different error than expected. Expected: %s",
			fmt.Errorf("service not connected to broker").Error())
	}

}

type testConsumeMsg struct {
	RbMQ        *AMQPConn
	RbMock      *rabbitMocks.RabbitMock
	ProcessFunc func([]byte) error
	Wg          sync.WaitGroup
	Msg         *amqp.Delivery
}

func (tstC *testConsumeMsg) setSimpleTestConsumeParams() {

	// Create Mock object and allocate it in  an AMQPConn object.
	tstC.RbMQ = &AMQPConn{Conn: &amqp.Connection{}}
	tstC.RbMock = rabbitMocks.NewRabbitMock()
	tstC.RbMQ.RbWrapper = tstC.RbMock

	// WaitGroupVariable
	tstC.Wg = sync.WaitGroup{}

	// Test message.
	tstC.Msg = &amqp.Delivery{
		Body:        []byte(`{"name": "TestMessage"}`),
		ContentType: MsgTypeJson,
	}

}

func (tstC *testConsumeMsg) setprocessMessageOK() {

	// Function to execute while receiving the message.
	tstC.ProcessFunc = func([]byte) error {
		return nil
	}

}

// Tests the normal execution of a message consumption when the received message
// is well processed and the ACK is sended to the broker.
func TestConsumeMessageOk(t *testing.T) {

	tstC := testConsumeMsg{}

	tstC.setSimpleTestConsumeParams()
	tstC.setprocessMessageOK()

	err := tstC.RbMQ.ConsumeOnQueue("testQueue", &tstC.Wg, tstC.ProcessFunc)
	if err != nil {
		t.Error("Expected none error but got one ", err)
	}

	tstC.RbMock.On("AckMsg", tstC.Msg, false).Return(nil)

	tstC.RbMQ.RbWrapper.Publish(&amqp.Channel{}, tstC.Msg.Body, tstC.Msg.ContentType, "", false)

	// Wait minimum time for consume message after publish.
	time.Sleep(time.Millisecond)

	tstC.RbMock.AssertNumberOfCalls(t, "AckMsg", 1)
	tstC.RbMock.AssertCalled(t, "AckMsg", tstC.Msg, false)

	tstC.Wg.Done()

}

func (tstC *testConsumeMsg) setBadMessageType() {
	// Test message.
	tstC.Msg = &amqp.Delivery{
		Body:        []byte(`Plain text`),
		ContentType: "text/plain",
	}
}

// Tests if a message that is not a json is rejected.
func TestRejectBadMsgType(t *testing.T) {

	tstC := testConsumeMsg{}

	tstC.setSimpleTestConsumeParams()
	tstC.setprocessMessageOK()
	tstC.setBadMessageType()

	err := tstC.RbMQ.ConsumeOnQueue("testQueue", &tstC.Wg, tstC.ProcessFunc)
	if err != nil {
		t.Error("Expected none error but got one ", err)
	}

	tstC.RbMock.On("RejectMsg", tstC.Msg, false).Return(nil)

	tstC.RbMQ.RbWrapper.Publish(&amqp.Channel{}, tstC.Msg.Body, tstC.Msg.ContentType, "", false)

	// Wait minimum time for consume message after publish.
	time.Sleep(time.Millisecond)

	// Retry messagges 3 times
	tstC.RbMock.AssertCalled(t, "RejectMsg", tstC.Msg, false)

	tstC.Wg.Done()

}

func (tstC *testConsumeMsg) setBadFormatedJsonMessage() {
	// Test message.
	tstC.Msg = &amqp.Delivery{
		Body:        []byte(`Fake Bad Json Message`),
		ContentType: MsgTypeJson,
	}
}

func TestRejectBadFormattedMsg(t *testing.T) {

	tstC := testConsumeMsg{}

	tstC.setSimpleTestConsumeParams()
	tstC.setprocessMessageOK()
	tstC.setBadFormatedJsonMessage()

	err := tstC.RbMQ.ConsumeOnQueue("testQueue", &tstC.Wg, tstC.ProcessFunc)
	if err != nil {
		t.Error("Expected none error but got one ", err)
	}

	tstC.RbMock.On("RejectMsg", tstC.Msg, false).Return(nil)

	tstC.RbMQ.RbWrapper.Publish(&amqp.Channel{}, tstC.Msg.Body, tstC.Msg.ContentType, "", false)

	// Wait minimum time for consume message after publish.
	time.Sleep(time.Millisecond)

	// Retry messagges 3 times
	tstC.RbMock.AssertCalled(t, "RejectMsg", tstC.Msg, false)

	tstC.Wg.Done()

}

func (tstC *testConsumeMsg) setprocessMessageKO() {

	// Function to execute while receiving the message.
	tstC.ProcessFunc = func([]byte) error {
		return fmt.Errorf("fake process message error.")
	}

}

// Tests rejection of a message when it hasn't been well processed.
func TestConsumeRejectMessageBadProcess(t *testing.T) {

	tstC := testConsumeMsg{}

	tstC.setSimpleTestConsumeParams()
	tstC.setprocessMessageKO()

	err := tstC.RbMQ.ConsumeOnQueue("testQueue", &tstC.Wg, tstC.ProcessFunc)
	if err != nil {
		t.Error("Expected none error but got one ", err)
	}

	tstC.RbMock.On("RejectMsg", tstC.Msg, false).Return(nil)

	tstC.RbMQ.RbWrapper.Publish(&amqp.Channel{}, tstC.Msg.Body, tstC.Msg.ContentType, "", false)

	// Wait minimum time for consume message after publish.
	time.Sleep(time.Millisecond)

	tstC.RbMock.AssertCalled(t, "RejectMsg", tstC.Msg, false)

	tstC.Wg.Done()

}
