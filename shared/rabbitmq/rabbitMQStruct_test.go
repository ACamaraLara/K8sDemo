package rabbitmq

import (
	"testing"

	"github.com/rs/zerolog"
)

// Init function executed before start tests to avoid verbose logs.
func init() {
	// Disables logger for unit testing.
	zerolog.SetGlobalLevel(zerolog.Disabled)

}

// This tests checks different combinations of AMQP init parameters.
func TestSetAMQPParams(t *testing.T) {
	// Subtest to test different parameter combinations
	t.Run("ValidParams", func(t *testing.T) {
		// Test case 1: Valid values
		t.Run("ValidValues", func(t *testing.T) {
			address := "localhost"
			port := 5672
			user := "user"
			password := "pass"

			rbMQ := NewAMQPConn(RabbitConfig{address, port, user, password})

			// Verify that values are assigned correctly
			if rbMQ.Address != address {
				t.Errorf("Expected Address='%s', but got: %s", address, rbMQ.Address)
			}
			if rbMQ.Port != port {
				t.Errorf("Expected Port='%d', but got: %d", port, rbMQ.Port)
			}
			if rbMQ.User != user {
				t.Errorf("Expected User='%s', but got: %s", user, rbMQ.User)
			}
			if rbMQ.Passwd != password {
				t.Errorf("Expected Pass='%s', but got: %s", password, rbMQ.Passwd)
			}
		})
	})
}
