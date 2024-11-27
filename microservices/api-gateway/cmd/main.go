package main

import (
	"fmt"
	"net/http"

	"api-gateway/internal/inputParams"
	"api-gateway/internal/restServer"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/rs/zerolog/log"
)

func main() {
	// Read input parameters.
	inputParams := inputParams.SetInputParams()
	fmt.Println("Starting application")

	// Init Logger with selected level.
	logger := logger.InitServiceLogger(inputParams.Logger)

	if err := logger.StartLokiLogPublishRoutine(); err != nil {
		fmt.Println("Error initializing Loki log routine:", err)
		return
	}

	rbMQ, err := rabbitmq.NewRabbitMQClient(&inputParams.Rabbit)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// Close connection before app ends.
	defer rbMQ.Conn.Close()

	if err := rbMQ.DeclareQueue("USERS", false, false, false, false); err != nil {
		log.Fatal().Msg(err.Error())
	}
	// Creates a muxer/router and adds routes to it (POSTS, GETS...).
	router := restRouter.NewRouter(restServer.InitRestRoutes(rbMQ))

	log.Info().Msg("Listening for HTTP requests on port " + inputParams.RESTPort)
	// Starts listening for HTTP requests.
	log.Fatal().Msg(http.ListenAndServe(inputParams.RESTPort, router).Error())
}
