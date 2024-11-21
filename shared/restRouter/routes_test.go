package restRouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func InitTestStatusHandlerRoute() Routes {

	return Routes{
		// Route to status GET.
		Route{
			Method:  http.MethodGet,
			Pattern: "/test",
			Handler: StatusHandler,
		}, Route{
			Method:  http.MethodPost,
			Pattern: "/test",
			Handler: StatusHandler,
		}, Route{
			Method:  http.MethodPut,
			Pattern: "/test",
			Handler: func(c *gin.Context) {
				fmt.Println("Called test put method.")
			},
		}, Route{
			Method:  http.MethodPatch,
			Pattern: "/test",
			Handler: func(c *gin.Context) {
				fmt.Println("Called test patch test method.")
			},
		}, Route{
			Method:  http.MethodDelete,
			Pattern: "/test",
			Handler: func(c *gin.Context) {
				fmt.Println("Called test Delete test method.")
			},
		}, Route{
			Method:  http.MethodOptions,
			Pattern: "/test",
			Handler: func(c *gin.Context) {
				fmt.Println("Called test options test method.")
			},
		}, Route{
			Method:  http.MethodHead,
			Pattern: "/test",
			Handler: func(c *gin.Context) {
				fmt.Println("Called test head test method.")
			},
		}}

}

// Creates a new router and checks that contains all expected routes
func TestDeclareRouter(t *testing.T) {
	routes := InitTestStatusHandlerRoute()
	router := NewRouter(routes)

	testRoutes := router.Routes()
	t.Run("TestCheckCorrectRoutesAdded", func(t *testing.T) {
		for idx, route := range routes {
			if testRoutes[idx].Path != route.Pattern ||
				testRoutes[idx].Method != route.Method {
				t.Error("Expected route not found in created router")
			}
		}
	})
	type testCase struct {
		Name           string
		Method         string
		ExpectedStatus int
	}
	testCases := []testCase{
		{
			Name:           "TestCheckStatusHandler",
			Method:         http.MethodGet,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "TestCheckBadRequestStatusHandler",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.Name, func(t *testing.T) {
			req, _ := http.NewRequest(testCase.Method, "/test", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, testCase.ExpectedStatus, resp.Code)
		})
	}
}
