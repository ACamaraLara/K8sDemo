package restServer

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
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

// Vector to store declared routes.
type Routes []Route

// Declare RoutesRepo and the handlers (actions) to do when
// one of them are called.
var RoutesRepo Routes

func InitRestRoutes( /*rbMQ *amqp.AMQPConn*/ ) {

	// In case of new routes, declare them here.
	RoutesRepo = Routes{
		// This is the service entry point. Main page of IP::port service.
		// Is just an example.
		Route{
			"GET", // Method
			"/",   // Pattern
			Main,  // Handler (action to do)
		},
		// Route to status GET.
		Route{
			"GET",
			"/status",
			StatusHandler,
		},
		// Route to POST method in a declared
		// section (SGORA in this case)
		Route{
			"POST",
			"/k8sdemo",
			func(c *gin.Context) {
				log.Info().Msg("Called POST method.")
				if c.Request.Method != http.MethodPost {
					c.Writer.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(c.Writer, "Expected GET method!.")
					return
				}

				c.Writer.WriteHeader(http.StatusOK)
			},
		},
	}
}

func NewRouter() *gin.Engine {
	// Avoid GIN verbose messages.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	// Create the muxer of our Rest server that will
	// route each request and response to correspondent
	// declared route. More info at:
	// https://go.dev/doc/tutorial/web-service-gin
	router := gin.Default()

	for _, route := range RoutesRepo {

		//Add route to declared router.
		switch route.Method {
		case "GET":
			router.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			router.POST(route.Pattern, route.HandlerFunc)
		default:
			log.Warn().Msg("Invalid HTTP method specified: " + route.Method)
		}
	}

	// Return muxer with all its added routes.
	return router
}
