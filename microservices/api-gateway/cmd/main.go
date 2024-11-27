package main

import (
	"net/http"
	"strconv"

	"api-gateway/internal/api"
	"api-gateway/internal/inputParams"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/rs/zerolog/log"
)

func main() {
	// Read input parameters.
	inputParams := inputParams.SetInputParams()
	log.Info().Msg("Starting api-gateway application")

	// Init Logger with selected level.
	logger := logger.InitServiceLogger(inputParams.Logger)

	if err := logger.StartLokiLogPublishRoutine(); err != nil {
		log.Fatal().Err(err).Msg("Starting api-gateway application")
	}

	log.Info().Msg("Initializing api-gateway routes.")
	router := restRouter.NewRouter(api.InitGatewayRoutes())
	log.Info().Msg("Routes initialized.")

	listenPort := ":" + strconv.Itoa(inputParams.RESTPort)
	log.Info().Msg("Listening for HTTP requests on port " + listenPort)

	// Starts listening for HTTP requests.
	log.Fatal().Msg(http.ListenAndServe(listenPort, router).Error())
}
