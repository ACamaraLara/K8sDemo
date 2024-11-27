package restServer

import (
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
)

// Creates a new router and checks that contains all expected routes
func TestDeclareRouter(t *testing.T) {
	routes := InitRestRoutes(&database.DBManager{})

	router := restRouter.NewRouter(routes)

	testRoutes := router.Routes()

	for idx, route := range routes {
		if testRoutes[idx].Path != route.Pattern ||
			testRoutes[idx].Method != route.Method {
			t.Fatal("Expected route not found in created router")
		}
	}
}
