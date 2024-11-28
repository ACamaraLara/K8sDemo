package api

import (
	"account-service/internal/account"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/jwtManager"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
	"github.com/gin-gonic/gin"
)

func SetAccountRoutes(accController *account.AccountController, jwtMgr *jwtManager.Manager) restRouter.Routes {
	return restRouter.Routes{
		restRouter.Route{
			Method:  http.MethodPost,
			Pattern: "/signup",
			Handler: func(c *gin.Context) {
				SignupHandler(c, accController)
			},
		},
		restRouter.Route{
			Method:  http.MethodPost,
			Pattern: "/login",
			Handler: func(c *gin.Context) {
				LoginHandler(c, accController, jwtMgr)
			},
		},
	}
}
