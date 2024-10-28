package inputParams

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
)

type InputParams struct {
	RESTPort int
	Logger   logger.LoggerConfig
	Rabbit   rabbitmq.RabbitConfig
}

// SetInputParams returns an object that stores service config parameters.
// @return the object with the seted parameters.
func SetInputParams() *InputParams {

	var inputParams InputParams

	flag.IntVar(&inputParams.RESTPort, "restPort", 8080, "REST server port.")
	logger.AddFlagsParams(&inputParams.Logger)
	inputParams.Rabbit.AddFlagsParams()

	flag.Parse()

	return &inputParams
}
