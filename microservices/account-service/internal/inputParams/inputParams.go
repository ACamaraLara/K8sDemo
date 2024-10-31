package inputParams

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
)

type InputParams struct {
	RESTPort int
	Logger   logger.LoggerConfig
	Rabbit   rabbitmq.RabbitConfig
	Mongo    mongodb.MongoConfig
}

// SetInputParams returns an object that stores service config parameters.
// @return the object with the seted parameters.
func SetInputParams() *InputParams {
	var inputParams InputParams

	flag.IntVar(&inputParams.RESTPort, "restPort",
		config.GetEnvironIntWithDefault("REST_PORT", 80), "RabbitMQ broker port (RABBITMQ__PORT).")
	logger.AddFlagsParams(&inputParams.Logger)
	inputParams.Rabbit.AddFlagsParams()
	inputParams.Mongo.AddFlagsParams()

	return &inputParams
}
