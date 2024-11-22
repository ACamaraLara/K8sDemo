package users

import (
	"io"
	"net/http"
	"time"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const accountServiceUri string = "http://account-service-svc.microservices.svc.cluster.local/"

var usersClient = &http.Client{
	Timeout: 5 * time.Second,
}

func AuthHandler(c *gin.Context, rbMQ *rabbitmq.RabbitMQClient, action string) {
	log.Info().Msg("Registering new user")
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expected POST method to register a user!"})
		return
	}
	// Proxie request to account-service.
	req, err := http.NewRequest(http.MethodPost, accountServiceUri+action, c.Request.Body)
	if err != nil {
		log.Error().Msg("failed to create request to account service:" + err.Error())
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "failed to create request to account service: " + err.Error()})
		return
	}

	// Copy headers from original request
	for key, value := range c.Request.Header {
		req.Header[key] = value
	}

	resp, err := usersClient.Do(req)
	if err != nil {
		log.Error().Msg("Error sending request to account service:" + err.Error())
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error sending request to account service: " + err.Error()})
		return
	}

	defer resp.Body.Close()

	// Copy the response from account service back to the client
	c.Writer.WriteHeader(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
