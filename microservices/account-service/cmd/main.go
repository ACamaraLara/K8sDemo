package main

import (
	"account-service/internal/inputParams"
	"account-service/internal/restServer"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/rs/zerolog/log"
)

func main() {
	// Read input parameters.
	inputParams := inputParams.SetInputParams()

	// Create the queue where ZeroLog will enqueue logs.
	loggerOutput := &logger.LoggerOutput{LogQueue: make(chan []byte, 1000)}

	// Init Logger with selected level.
	if err := logger.InitServiceLogger(inputParams.Logger, loggerOutput); err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	if err := logger.StartLokiLogPublishRoutine(loggerOutput); err != nil {
		fmt.Println("Error initializing Loki log routine:", err)
		return
	}

	// Connect to MongoDB data base with the input parameters.
	log.Info().Msg("Connecting to mongodb..." + inputParams.Mongo.GetURL())
	mongoConn := mongodb.NewMongoDBClient(inputParams.Mongo)

	err := mongoConn.ConnectMongoClient()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	defer mongoConn.DisconnectMongoClient()

	rbMQ := rabbitmq.NewAMQPConn(inputParams.Rabbit)

	if err := rbMQ.InitConnection(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	// FIXME: move queue declaration to msgbrokerlib
	if err := rbMQ.DeclareQueue("USERS", false, false, false, false); err != nil {
		log.Fatal().Msg(err.Error())
	}

	// Defer means that shuld be executed at the end of current scope.
	defer rbMQ.CloseConnection()

	restServer.InitRestRoutes(mongoConn)

	// Creates a muxer/router and adds routes to it (POSTS, GETS...).
	router := restRouter.NewRouter()

	listenPort := ":" + strconv.Itoa(inputParams.RESTPort)
	log.Info().Msg("Listening for HTTP requests on port " + listenPort)

	// Starts listening for HTTP requests.
	log.Fatal().Msg(http.ListenAndServe(listenPort, router).Error())
	log.Info().Msg("Exiting...")

}
