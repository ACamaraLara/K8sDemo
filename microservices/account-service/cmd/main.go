package main

import (
	"account-service/internal/inputParams"
	"account-service/internal/restServer"
	"context"
	"net/http"
	"strconv"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
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

	log.Info().Msg("Connecting to mongodb..." + inputParams.Mongo.GetURL())
	mongoDB, err := mongodb.NewMongoDBClient(ctx, &inputParams.Mongo)
	if err != nil {

	}

	defer mongoDB.Client.Disconnect(ctx)

	router := restRouter.NewRouter(restServer.InitRestRoutes(mongoDB))

	listenPort := ":" + strconv.Itoa(inputParams.RESTPort)
	log.Info().Msg("Listening for HTTP requests on port " + listenPort)

	log.Fatal().Msg(http.ListenAndServe(listenPort, router).Error())
	log.Info().Msg("Exiting...")

}
