package restServer

import (
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
	"github.com/gin-gonic/gin"
)

func InitRestRoutes(mongoClient *mongodb.MongoDBClient) {

	// In case of new routes, declare them here.
	restRouter.RoutesRepo = restRouter.Routes{
		// Route to status GET.
		restRouter.Route{
			Method:  http.MethodGet,
			Pattern: "/status",
			Handler: restRouter.StatusHandler,
		},
		// Route to register new user.
		restRouter.Route{
			Method:  http.MethodPost,
			Pattern: "/signup",
			Handler: func(c *gin.Context) {
				SignupHandler(c, mongoClient)
			},
		},
		// Route to login user.
		restRouter.Route{
			Method:  http.MethodPost,
			Pattern: "/login",
			Handler: func(c *gin.Context) {
				LoginHandler(c, mongoClient)
			},
		},
	}
}
