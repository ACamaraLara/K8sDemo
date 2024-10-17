package restServer

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// This is the main page of our web server. The information will be shown
// while doing a REST GET to the URL of this service
// (http://www.localhost:8080/ in this case).
func Main(c *gin.Context) {
	log.Info().Msg("Called GET main method.")
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "Expected GET method!.")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(c.Writer, "Welcome to Kubernetes Blockchain Demo main Page!")
	fmt.Fprintln(c.Writer, "Usage:")
	fmt.Fprintln(c.Writer, "POST -> URL/k8sdemo: stores given Json info inside dataBase as new entry. In case of bad data, returns an error (422)")
}

// This is the status page of our web server. It just returns if the service
// is available at anytime.
func StatusHandler(c *gin.Context) {
	log.Info().Msg("Called GET status method.")
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "Expected GET method!.")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
