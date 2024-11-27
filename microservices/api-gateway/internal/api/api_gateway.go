package api

import (
	"api-gateway/internal/api/account"
	"api-gateway/internal/api/service"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
)

type serviceHandler interface {
	RegisterRoutes() restRouter.Routes
}

func InitGatewayRoutes() restRouter.Routes {
	serviceRouter := service.NewServiceRouter()
	serviceHandlers := createHandlers(serviceRouter)
	return registerRoutes(serviceHandlers)
}

func createHandlers(sr *service.ServiceRouter) []serviceHandler {
	// More microservices handlers should be added here,
	return []serviceHandler{account.NewAccountHandler(sr)}
}

func registerRoutes(serviceHandlers []serviceHandler) restRouter.Routes {
	routes := restRouter.Routes{}
	for _, sh := range serviceHandlers {
		// Process routes of all microservices.
		routes = append(routes, sh.RegisterRoutes()...)
	}
	return routes
}
