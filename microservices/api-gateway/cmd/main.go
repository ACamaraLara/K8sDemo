package main

import (
	"fmt"
	"time"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	serviceName string = "API-gateway"
)

func main() {
	// Read input parameters.
	fmt.Println("Starting application")

	// Create channel where zerolog will enqueue all logs
	loggerOutput := &logger.LoggerOutput{LogQueue: make(chan []byte, 1000)}

	// Init Logger with selected level.
	if err := logger.InitServiceLogger(logger.LoggerConfig{LogLevel: zerolog.LevelInfoValue},
		serviceName, loggerOutput); err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	if err := logger.StartLokiLogPublishRoutine(loggerOutput); err != nil {
		fmt.Println("Error initializing Loki log routine:", err)
		return
	}

	go func() {
		count := 0
		for {
			log.Info().Int("count", count).Msg("Incrementing counter loki + stdout")
			// log.Info().Msgf("Incrementing counter loki: %d", count)
			count++
			time.Sleep(10 * time.Second)
		}
	}()

	select {}
}
