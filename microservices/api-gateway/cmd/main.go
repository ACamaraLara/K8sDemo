package main

import (
	"fmt"

	"github.com/ACamaraLara/K8sBlockChainDemo/microservices/api-gateway/pkg/inputParams"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/logger"
)

func main() {
	// Read input parameters.
	inputParams := inputParams.SetInputParams()

	// Create the queue where ZeroLog will add all published Logs.
	// :TODO: Think about what would be an appropiate queue size.
	loggerOutput := &logger.LoggerOutput{LogQueue: make(chan []byte, 1000)}

	// Init Logger with selected level.
	if err := logger.InitServiceLogger(inputParams.Logger, "SGORA-publisher", loggerOutput); err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	if err := logger.StartLokiLogPublishRoutine(loggerOutput); err != nil {
		fmt.Println("Error initializing Loki log routine:", err)
		return
	}
}
