package inputParams

import (
	"flag"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/go-playground/assert/v2"
)

func TestSetInputParams_DefaultValues(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)

	inputParams := SetInputParams()

	assert.Equal(t, ":8080", inputParams.RESTPort)

	assert.Equal(t, rabbitmq.DefaultRabbitHost, inputParams.Rabbit.Host)
	assert.Equal(t, rabbitmq.DefaultRabbitPort, inputParams.Rabbit.Port)
	assert.Equal(t, rabbitmq.DefaultRabbitUser, inputParams.Rabbit.User)
	assert.Equal(t, rabbitmq.DefaultRabbitPass, inputParams.Rabbit.Passwd)
}
