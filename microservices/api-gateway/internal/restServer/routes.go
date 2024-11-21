package restServer

import (
	"api-gateway/internal/restServer/users"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
)

func InitRestRoutes(rbMQ *rabbitmq.RabbitMQClient) restRouter.Routes {

	// In case of new routes, declare them here.
	routes := restRouter.Routes{
		// This is the service entry point. Main page of IP::port service.
		// Is just an example.
		restRouter.Route{
			Method:  http.MethodGet, // Method
			Pattern: "/",            // Pattern
			Handler: Main,           // Handler (action to do)
		},
		// Route to status GET.
		restRouter.Route{
			Method:  http.MethodGet,
			Pattern: "/status",
			Handler: restRouter.StatusHandler,
		},
	}

	return append(routes, users.GetUsersRoutes(rbMQ)...)
}
