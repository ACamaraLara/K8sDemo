package rabbitmq

import (
	"encoding/json"
	"fmt"
)

// Function to check if a json is well formatted before sending it.
func isJsonWellFormated(jsonBody []byte) bool {
	return (json.Unmarshal([]byte(jsonBody), &map[string]interface{}{}) == nil)
}

// PublishJsonMessage sends a message in json format to the broker.
// @param jsonBody Is an array of bytes that contains a json message.
// @param exchangeName Where the message will be published. if its empty ("")
// message is sended to default direct.exchange.
// @param routingKey binding key of a queue for the given exchange. If direct.exchange
// is selected, routingKey should be the name of the queue. This is
// because every declared queue gets an implicit route to the default exchange.
// @param mandatory If true, message goes to the specified Queue to the first position.
// @throws error with the specific message error. Null otherwise.
func (rbMQ *AMQPConn) PublishJsonMessage(jsonBody []byte, exchangeName, routingKey string, mandatory bool) error {

	if !isJsonWellFormated(jsonBody) {
		return fmt.Errorf("json bad formatted, cannot send it to the broker")
	}
	err := rbMQ.RbWrapper.Publish(rbMQ.Channel, jsonBody, exchangeName, routingKey, mandatory)
	if err != nil {

		return err
	}

	return nil
}
