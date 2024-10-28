package restRouter

import (
	"encoding/json"
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
		case "GET":
			router.GET(route.Pattern, route.Handler)
		case "POST":
			router.POST(route.Pattern, route.Handler)
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

func ReadRequestBody(c *gin.Context) ([]byte, error) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body."})
		return nil, err
	}

	// Close the request body
	if err := c.Request.Body.Close(); err != nil {
		log.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "failed to close request body:" + err.Error()})
		return nil, err
	}

	return body, nil
}

func DecodeRequestBody(c *gin.Context, target interface{}) error {
	// Read the request body
	body, err := ReadRequestBody(c)
	if err != nil {
		return err // The error is already handled in ReadRequestBody
	}

	// Unmarshal the body into the target object
	if err := json.Unmarshal(body, target); err != nil {
		// If cannot decode received JSON, return the error to the client.
		c.Header("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(c.Writer).Encode(gin.H{"error": "unprocessable entity"}); err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		return err
	}

	return nil
}
