package rabbitmq

import (
	"encoding/json"
	"fmt"
)

// Function to check if a json is well formatted before sending it.
func isJsonWellFormated(jsonBody []byte) bool {
	return (json.Unmarshal(jsonBody, &map[string]interface{}{}) == nil)
}

// PublishObject sends a message in json format to the broker.
// @param jsonBody Is an array of bytes that contains a json message.
// @param exchangeName Where the message will be published. if its empty ("")
// message is sended to default direct.exchange.
// @param routingKey binding key of a queue for the given exchange. If direct.exchange
// is selected, routingKey should be the name of the queue. This is
// because every declared queue gets an implicit route to the default exchange.
// @param mandatory If true, message goes to the specified Queue to the first position.
// @throws error with the specific message error. Null otherwise.
func (rbMQ *AMQPConn) PublishObject(messageObject interface{}, exchangeName, routingKey string, mandatory bool) error {

	jsonBody, err := json.Marshal(messageObject)
	if err != nil {
		return err
	}

	err = rbMQ.RbWrapper.Publish(rbMQ.Channel, jsonBody, exchangeName, routingKey, mandatory)
	if err != nil {

		return err
	}

	return nil
}

// PublishJson sends a pre-formatted JSON message to the broker.
// @param message A byte array containing the JSON-formatted message.
// @param exchangeName Specifies the exchange to publish the message to. If empty (""),
// the message is sent to the default direct exchange.
// @param routingKey Binding key for a queue associated with the given exchange. If using
// the default direct exchange, routingKey should match the queue name to ensure delivery,
// as each queue has an implicit route to the default exchange.
// @param mandatory If set to true, ensures the message is routed to the specified queue
// in a prioritized manner, failing if no route exists.
// @throws error Returns a specific error if the JSON is malformed or if there is an error
// during message publishing; nil otherwise.
func (rbMQ *AMQPConn) PublishJson(message []byte, exchangeName, routingKey string, mandatory bool) error {

	if !isJsonWellFormated(message) {
		return fmt.Errorf("json bad formatted, cannot send it to the broker")
	}

	err := rbMQ.RbWrapper.Publish(rbMQ.Channel, message, exchangeName, routingKey, mandatory)
	if err != nil {

		return err
	}

	return nil
}
