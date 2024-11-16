package inputParams

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/go-playground/assert/v2"
)

func TestSetInputParams_DefaultValues(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)

	inputParams := SetInputParams()
	fmt.Printf("input params = %+v", inputParams)

	assert.Equal(t, 80, inputParams.RESTPort)

	assert.Equal(t, rabbitmq.DefaultRabbitHost, inputParams.Rabbit.Host)
	assert.Equal(t, rabbitmq.DefaultRabbitPort, inputParams.Rabbit.Port)
	assert.Equal(t, rabbitmq.DefaultRabbitUser, inputParams.Rabbit.User)
	assert.Equal(t, rabbitmq.DefaultRabbitPass, inputParams.Rabbit.Passwd)

	assert.Equal(t, mongodb.DefaultMongoDBHost, inputParams.Mongo.Host)
	assert.Equal(t, mongodb.DefaultMongoDBPort, inputParams.Mongo.Port)
	assert.Equal(t, mongodb.DefaultMongoDBName, inputParams.Mongo.DbName)
	assert.Equal(t, mongodb.DefaultMongoDBCollection, inputParams.Mongo.Collection)
	assert.Equal(t, mongodb.DefaultMongoDBUserName, inputParams.Mongo.User)
	assert.Equal(t, mongodb.DefaultMongoDBUserPass, inputParams.Mongo.Passwd)
}
