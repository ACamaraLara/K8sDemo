package rabbitmq

import (
	"fmt"
	"sync"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPublishObject(t *testing.T) {
	mockChannel := new(MockChannel)

	rbMQClient := &RabbitMQClient{
		Ch: mockChannel,
	}

	type testPublish struct {
		Name          string
		MessageObject interface{}
		ExpectedError error
	}

	testCases := []testPublish{
		{
			Name: "TestValidObject",
			MessageObject: map[string]string{
				"name": "testValidObject",
			},
			ExpectedError: nil,
		},
		{
			Name:          "TestInvalidObject",
			MessageObject: func() {}, // Invalid object that can't be marshalled to JSON
			ExpectedError: fmt.Errorf("json: unsupported type: func()"),
		},
	}

	mockChannel.On("Publish", "testQueue", "", true, false, mock.Anything).
		Return(nil).Once()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := rbMQClient.PublishObject(tc.MessageObject, "testQueue", true)

			// Assert the expected error
			if err != nil && err.Error() != tc.ExpectedError.Error() {
				t.Fatalf("Expected error %v, got %v", tc.ExpectedError, err)
			}

			if tc.ExpectedError == nil {
				mockChannel.AssertExpectations(t)
			}
		})
	}
}

func TestPublishJson(t *testing.T) {
	mockChannel := new(MockChannel)

	rbMQClient := &RabbitMQClient{
		Ch: mockChannel,
	}

	type testPublish struct {
		Name          string
		JsonBody      []byte
		ExpectedError error
	}

	testCases := []testPublish{
		{
			Name:          "TestValidJson",
			JsonBody:      []byte("{\"name\":\"testValidObject\"}"),
			ExpectedError: nil,
		},
		{
			Name:          "TestInvalidJson",
			JsonBody:      []byte("InvalidJson"),
			ExpectedError: fmt.Errorf("json form broker wrong formatted"),
		},
	}

	mockChannel.On("Publish", "testQueue", "", true, false, mock.Anything).
		Return(nil).Once()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := rbMQClient.PublishJson(tc.JsonBody, "testQueue", true)

			if err != nil && err.Error() != tc.ExpectedError.Error() {
				t.Fatalf("Expected error %v, got %v", tc.ExpectedError, err)
			}

			if tc.ExpectedError == nil {
				mockChannel.AssertExpectations(t)
			}
		})
	}
}

func TestDeclareQueue(t *testing.T) {
	mockChannel := new(MockChannel)

	rbMQClient := &RabbitMQClient{
		Ch: mockChannel,
	}

	testCases := []struct {
		Name          string
		ExpectedError error
	}{
		{
			Name:          "TestSuccessQueueDeclare",
			ExpectedError: nil,
		},
		{
			Name:          "TestErrorInQueueDeclare",
			ExpectedError: fmt.Errorf("Queue declare failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			mockChannel.On("QueueDeclare", "testQueue", true, false, false, false, mock.Anything).
				Return(amqp.Queue{}, tc.ExpectedError).Once()

			err := rbMQClient.DeclareQueue("testQueue", true, false, false, false)

			if tc.ExpectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.ExpectedError.Error())
			}

			mockChannel.AssertExpectations(t)
		})
	}
}

func TestConsumeOnQueueFailDeliveryCreation(t *testing.T) {
	mockChannel := new(MockChannel)

	rbMQClient := &RabbitMQClient{
		Ch: mockChannel,
	}

	testError := fmt.Errorf("mock error")

	mockChannel.On("Consume", "testQueue", "", false, false, false, false, mock.Anything).
		Return(make(<-chan Message), testError).Once()

	err := rbMQClient.ConsumeOnQueue("testQueue", &sync.WaitGroup{}, false, false, nil)

	assert.Contains(t, err.Error(), testError.Error(), "Error is different than expected")

	mockChannel.AssertExpectations(t)
}

type testConsume struct {
	RbtClient   *RabbitMQClient
	MockChannel *MockChannel
	Wg          *sync.WaitGroup
	MsgCh       chan Message
	MockMsg     *MockMessage
}

func newTestConsumeMessage() *testConsume {

	mockChannel := new(MockChannel)
	return &testConsume{
		RbtClient:   &RabbitMQClient{Ch: mockChannel},
		MockChannel: mockChannel,
		Wg:          &sync.WaitGroup{},
		MsgCh:       make(chan Message, 3),
		MockMsg:     &MockMessage{Body: []byte("TestMessage")},
	}
}

func (tstC *testConsume) addMessageToChannel(msg Message) <-chan Message {
	tstC.MsgCh <- msg
	return (<-chan Message)(tstC.MsgCh)
}

func TestConsumeOnQueueProcessMessage(t *testing.T) {
	testConsume := newTestConsumeMessage()

	// Insert test message inside channel
	processTestFunc := func(m Message) error {
		assert.Equal(t, m, testConsume.MockMsg)
		testConsume.Wg.Done()
		return nil
	}

	testConsume.MockChannel.On("Consume", "testQueue", "", false, false, false, false, mock.Anything).
		Return(testConsume.addMessageToChannel(testConsume.MockMsg), nil).Once()

	t.Run("TestConsumeMessageWithACK", func(t *testing.T) {
		testConsume.MockMsg.On("Ack", false).Return(nil).Once()

		err := testConsume.RbtClient.ConsumeOnQueue("testQueue", testConsume.Wg, true, false, processTestFunc)

		assert.NoError(t, err)

		testConsume.MockChannel.AssertExpectations(t)

		testConsume.Wg.Wait()

	})

	t.Run("TestConsumeMessageACKErr", func(t *testing.T) {
		expectedErr := fmt.Errorf("testErrorSendigAck")
		testConsume.MockMsg.On("Ack", false).Return(expectedErr).Once()

		testConsume.Wg.Add(1)
		testConsume.addMessageToChannel(testConsume.MockMsg)

		testConsume.MockChannel.AssertExpectations(t)

		testConsume.Wg.Wait()

	})

}

func TestConsumeOnQueueRejectMessage(t *testing.T) {
	testConsume := newTestConsumeMessage()

	// Insert test message inside channel
	processTestFunc := func(m Message) error {
		assert.Equal(t, m, testConsume.MockMsg)
		return fmt.Errorf("test error processing message")
	}

	testConsume.MockChannel.On("Consume", "testQueue", "", false, false, false, false, mock.Anything).
		Return(testConsume.addMessageToChannel(testConsume.MockMsg), nil).Once()

	t.Run("TestRejectMessageNoError", func(t *testing.T) {
		testConsume.MockMsg.On("Reject", false).Return(nil).Once()

		err := testConsume.RbtClient.ConsumeOnQueue("testQueue", testConsume.Wg, false, false, processTestFunc)

		assert.NoError(t, err)
		assert.True(t, testConsume.MockChannel.AssertExpectations(t))

	})

	t.Run("TestConsumeMessageACKErr", func(t *testing.T) {
		expectedErr := fmt.Errorf("testErrorRejectingMessage")
		testConsume.MockMsg.On("Reject", false).Return(expectedErr).Once()

		testConsume.addMessageToChannel(testConsume.MockMsg)

		testConsume.MockChannel.AssertExpectations(t)
	})

}
