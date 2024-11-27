package inputParams

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database/config"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/utils"
)

type InputParams struct {
	RESTPort string
	DBType   string
	Logger   logger.LoggerConfig
	Rabbit   rabbitmq.RabbitConfig
	DBConf   config.DBConfig
}

// SetInputParams returns an object that stores service config parameters.
// @return the object with the seted parameters.
func SetInputParams() *InputParams {
	var inputParams InputParams

	flag.StringVar(&inputParams.RESTPort, "restPort",
		utils.GetEnvironWithDefault("REST_PORT", "80"), "Rest listening port (REST_PORT).")
	flag.StringVar(&inputParams.DBType, "db-type",
		utils.GetEnvironWithDefault("dbType", "mongo"), "Type of database that will be used (DB_TYPE).")
	logger.AddFlagsParams(&inputParams.Logger)
	inputParams.Rabbit.AddFlagsParams()
	inputParams.DBConf.AddFlagsParams()

	return &inputParams
}
