package restServer

import (
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
	"github.com/gin-gonic/gin"
)

func InitRestRoutes(mongoClient *database.DBManager) restRouter.Routes {

	// In case of new routes, declare them here.
	return restRouter.Routes{
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
