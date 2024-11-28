package inputParams

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
)

type InputParams struct {
	RESTPort string
	Logger   logger.LoggerConfig
	Rabbit   rabbitmq.RabbitConfig
}

// SetInputParams returns an object that stores service config parameters.
// @return the object with the seted parameters.
func SetInputParams() *InputParams {

	var inputParams InputParams

	flag.StringVar(&inputParams.RESTPort, "restPort",
		config.GetEnvironWithDefault("REST_PORT", "8080"), "Port to listen http requests (REST_PORT).")
	inputParams.RESTPort = ":" + inputParams.RESTPort
	logger.AddFlagsParams(&inputParams.Logger)
	inputParams.Rabbit.AddFlagsParams()

	return &inputParams
}
