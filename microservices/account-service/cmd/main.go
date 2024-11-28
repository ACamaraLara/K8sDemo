package main

import (
	"account-service/internal/account"
	"account-service/internal/api"
	"account-service/internal/inputParams"
	"context"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/jwtManager"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting account-service application")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // This will allow us to cancel the context on app shutdown

	inputParams := inputParams.SetInputParams()

	logger := logger.InitServiceLogger(inputParams.Logger)

	if err := logger.StartLokiLogPublishRoutine(); err != nil {
		log.Fatal().Err(err).Msg("Error initializing Loki log routine.")
	}

	mongoDB, err := mongodb.NewMongoDBClient(ctx, &inputParams.Mongo)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing MongoDB client.")
	}

	defer mongoDB.Client.Disconnect(ctx)

	accCtrl := account.NewAccountController(mongoDB)
	jwtMgr, err := jwtManager.NewManager(&inputParams.JWT)
	if err != nil {
		log.Fatal().Msgf("Error creating JWT manager object. %+v", err)
	}

	router := restRouter.NewRouter(api.SetAccountRoutes(accCtrl, jwtMgr))

	log.Info().Msg("Listening for HTTP requests on port " + inputParams.RESTPort)

	log.Fatal().Msg(http.ListenAndServe(inputParams.RESTPort, router).Error())
	log.Info().Msg("Exiting...")

}
