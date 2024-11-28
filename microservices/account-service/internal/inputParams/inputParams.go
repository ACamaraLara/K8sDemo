package inputParams

import (
	"flag"
	"time"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/jwtManager"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
)

type InputParams struct {
	RESTPort          string
	JWTExpirationTime time.Duration
	Logger            logger.LoggerConfig
	Rabbit            rabbitmq.RabbitConfig
	Mongo             mongodb.MongoConfig
	JWT               jwtManager.Config
}

// SetInputParams returns an object that stores service config parameters.
// @return the object with the seted parameters.
func SetInputParams() *InputParams {
	var inputParams InputParams

	flag.StringVar(&inputParams.RESTPort, "restPort",
		config.GetEnvironWithDefault("REST_PORT", "80"), "Port to listen REST requests (RABBITMQ__PORT).")
	inputParams.RESTPort = ":" + inputParams.RESTPort
	logger.AddFlagsParams(&inputParams.Logger)
	inputParams.Rabbit.AddFlagsParams()
	inputParams.Mongo.AddFlagsParams()
	inputParams.JWT.AddFlagsParams()

	return &inputParams
}
