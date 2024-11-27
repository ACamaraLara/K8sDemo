package inputParams

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/utils"
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
		utils.GetEnvironWithDefault("RABBITMQ__PORT", "8080"), "RabbitMQ broker port (RABBITMQ__PORT).")
	logger.AddFlagsParams(&inputParams.Logger)
	inputParams.Rabbit.AddFlagsParams()

	return &inputParams
}
