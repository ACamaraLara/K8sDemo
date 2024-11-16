package restRouter

import (
	"fmt"
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

// Vector to store declared routes.
type Routes []Route

// Declare RoutesRepo and the handlers (actions) to do when
// one of them are called.
var RoutesRepo Routes

func NewRouter() *gin.Engine {
	// Avoid GIN verbose messages.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	// Create the muxer of our Rest server that will
	// route each request and response to correspondent
	// declared route. More info at:
	// https://go.dev/doc/tutorial/web-service-gin
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1 MB limit
		c.Next()
	})

	for _, route := range RoutesRepo {

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

	// Return muxer with all its added routes.
	return router
}

// This is the default status handler that will be used in every service.
func StatusHandler(c *gin.Context) {
	log.Info().Msg("Called GET status method.")
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "Expected GET method!.")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
