package main

import (
	"fmt"

	"github.com/ACamaraLara/K8sBlockChainDemo/microservices/api-gateway/pkg/inputParams"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
)

func main() {
	// Read input parameters.
	inputParams.SetInputParams()
	helloint := config.GetEnvironIntWithDefault("Hello", 3)
	fmt.Printf("Read variable = %d", helloint)
}
