package inputParams

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database/config"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/go-playground/assert/v2"
)

func TestSetInputParams_DefaultValues(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)

	inputParams := SetInputParams()
	fmt.Printf("input params = %+v", inputParams)

	assert.Equal(t, "80", inputParams.RESTPort)
	assert.Equal(t, "mongo", inputParams.DBType)

	assert.Equal(t, rabbitmq.DefaultRabbitHost, inputParams.Rabbit.Host)
	assert.Equal(t, rabbitmq.DefaultRabbitPort, inputParams.Rabbit.Port)
	assert.Equal(t, rabbitmq.DefaultRabbitUser, inputParams.Rabbit.User)
	assert.Equal(t, rabbitmq.DefaultRabbitPass, inputParams.Rabbit.Passwd)

	assert.Equal(t, config.DefaultMongoDBHost, inputParams.DBConf.Host)
	assert.Equal(t, config.DefaultMongoDBPort, inputParams.DBConf.Port)
	assert.Equal(t, config.DefaultMongoDBName, inputParams.DBConf.DbName)
	assert.Equal(t, config.DefaultMongoDBCollection, inputParams.DBConf.Collections[0])
	assert.Equal(t, config.DefaultMongoDBUserName, inputParams.DBConf.User)
	assert.Equal(t, config.DefaultMongoDBUserPass, inputParams.DBConf.Passwd)
}
