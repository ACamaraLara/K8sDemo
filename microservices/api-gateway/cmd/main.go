package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"api-gateway/internal/restServer"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/inputParams"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Read input parameters.
	inputParams := inputParams.SetInputParams()
	fmt.Println("Starting application")

	// Create channel where zerolog will enqueue all logs
	loggerOutput := &logger.LoggerOutput{LogQueue: make(chan []byte, 10000)}

	// Init Logger with selected level.
	if err := logger.InitServiceLogger(logger.LoggerConfig{LogLevel: zerolog.LevelDebugValue},
		loggerOutput); err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	if err := logger.StartLokiLogPublishRoutine(loggerOutput); err != nil {
		fmt.Println("Error initializing Loki log routine:", err)
		return
	}

	restServer.InitRestRoutes() //&rbMQ)

	// Creates a muxer/router and adds routes to it (POSTS, GETS...).
	router := restServer.NewRouter()

	go func() {
		count := 0
		for {
			log.Debug().Int("count", count).Msg("Incrementing counter loki + stdout")
			count++
			time.Sleep(10 * time.Second)
		}
	}()

	listenPort := ":" + strconv.Itoa(inputParams.RESTPort)
	log.Info().Msg("Listening for HTTP requests on port " + listenPort)

	// Starts listening for HTTP requests.
	log.Fatal().Msg(http.ListenAndServe(listenPort, router).Error())
	select {}
}
