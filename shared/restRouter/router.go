package restRouter

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Struct that stores information of a single
// route for current service.
type Route struct {
	Method  string
	Pattern string
	Handler gin.HandlerFunc
}

type Routes []Route

func NewRouter(routes Routes) *gin.Engine {
	// Avoid GIN verbose messages.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	// Status handler should be available in all services inside the cluster.
	appendStatusHandler(routes)
	registerRoutes(router, routes)
	return router
}

func appendStatusHandler(routes Routes) Routes {
	return append(routes, Route{Method: http.MethodGet, Pattern: "/status", Handler: statusHandler})
}

func registerRoutes(router *gin.Engine, routes Routes) {
	for _, route := range routes {

		//Add route to declared router.
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.Handler)
		case http.MethodPost:
			router.POST(route.Pattern, route.Handler)
		case http.MethodPut:
			router.PUT(route.Pattern, route.Handler)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.Handler)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.Handler)
		case http.MethodHead:
			router.HEAD(route.Pattern, route.Handler)
		case http.MethodOptions:
			router.OPTIONS(route.Pattern, route.Handler)
		default:
			log.Warn().Msg("Invalid HTTP method specified: " + route.Method)
		}
	}
}

// This is the default status handler that will be used in every service.
func statusHandler(c *gin.Context) {
	log.Info().Msg("Status endpoint hit")

	if c.Request.Method != http.MethodGet {
		log.Error().Msg("Invalid method for status endpoint")
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "invalid method, only GET is allowed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
