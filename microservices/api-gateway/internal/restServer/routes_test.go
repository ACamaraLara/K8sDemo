package restServer

import (
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
)

// Creates a new router and checks that contains all expected routes
func TestDeclareRouter(t *testing.T) {
	InitRestRoutes(&rabbitmq.AMQPConn{})

	router := restRouter.NewRouter()

	testRoutes := router.Routes()

	for idx, route := range restRouter.RoutesRepo {
		if testRoutes[idx].Path != route.Pattern ||
			testRoutes[idx].Method != route.Method {
			t.Fatal("Expected route not found in created router")
		}
	}
}
