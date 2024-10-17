package restServer

import (
	"testing"
)

// Creates a new router and checks that contains all expected routes
func TestDeclareRouter(t *testing.T) {
	router := NewRouter()

	testRoutes := router.Routes()

	for idx, route := range RoutesRepo {
		if testRoutes[idx].Path != route.Pattern ||
			testRoutes[idx].Method != route.Method {
			t.Fatal("Expected route not found in created router")
		}
	}
}
