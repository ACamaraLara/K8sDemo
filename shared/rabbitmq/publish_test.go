package rabbitmq

import (
	"fmt"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Tests valid and invalid publication of json object to the broker.
func TestPublishObject(t *testing.T) {

	rbMQ := &AMQPConn{RbWrapper: &RabbitMock{Msgs: make(chan amqp.Delivery, 1)}}

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

	for _, testCase := range testCases {

		t.Run(testCase.Name, func(t *testing.T) {
			err := rbMQ.PublishObject(testCase.MessageObject, "testExchange", "testRoutingKey", false)

			// Check if the expected error is obtained
			if err == nil && testCase.ExpectedError != nil {
				t.Fatalf("Expected an error, but got none.")
			} else if err != nil && testCase.ExpectedError == nil {
				t.Fatalf("Error not expected, but one given: %s", err.Error())
			} else if err != nil && testCase.ExpectedError != nil && err.Error() != testCase.ExpectedError.Error() {
				t.Fatalf("Got a different error than expected. Expected: %s, Got: %s", testCase.ExpectedError.Error(), err.Error())
			}
		})

	}
}

// Tests valid and invalid publication of json message to the broker.
func TestPublishJsonMessage(t *testing.T) {

	rbMQ := &AMQPConn{RbWrapper: &RabbitMock{Msgs: make(chan amqp.Delivery, 1)}}

	type testPublish struct {
		Name          string
		JsonBody      []byte
		ExpectedError error
	}

	testCases := []testPublish{
		{
			Name:          "TestValidJson",
			JsonBody:      []byte(`{"name": "testValidJson"}`),
			ExpectedError: nil,
		},
		{
			Name:          "TestInvalidJson",
			JsonBody:      []byte(`{"name": "testInvalidJson"`),
			ExpectedError: fmt.Errorf("json bad formatted, cannot send it to the broker"),
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.Name, func(t *testing.T) {
			err := rbMQ.PublishJson(testCase.JsonBody, "testExchange", "testRoutingKey", false)

			// Check if the expected error is obtained
			if err == nil && testCase.ExpectedError != nil {
				t.Fatalf("Expected an error, but got none.")
			} else if err != nil && testCase.ExpectedError == nil {
				t.Fatalf("Error not expected, but one given.: %s", err.Error())
			} else if err != nil && testCase.ExpectedError != nil && err.Error() != testCase.ExpectedError.Error() {
				t.Fatalf("Got a different error than expected. Expected: %s, Got: %s", testCase.ExpectedError.Error(), err.Error())
			}
		})

	}

}
