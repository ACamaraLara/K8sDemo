package main

import (
	"account-service/internal/inputParams"
	"account-service/internal/restServer"
	"context"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // This will allow us to cancel the context on app shutdown

	inputParams := inputParams.SetInputParams()

	logger := logger.InitServiceLogger(inputParams.Logger)

	if err := logger.StartLokiLogPublishRoutine(); err != nil {
		log.Error().Err(err).Msg("Error initializing Loki log routine:")
		return
	}
	database, err := database.NewDatabase(ctx, inputParams.DBType, &inputParams.DBConf)
	if err != nil {

	}
	/*:ToDo defer mongoDB.Client.Disconnect(ctx)*/
	router := restRouter.NewRouter(restServer.InitRestRoutes(database))

	log.Info().Msg("Listening for HTTP requests on port " + inputParams.RESTPort)
	log.Fatal().Msg(http.ListenAndServe(inputParams.RESTPort, router).Error())
	log.Info().Msg("Exiting...")

}
