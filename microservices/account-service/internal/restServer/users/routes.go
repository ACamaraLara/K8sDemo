package users

import (
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/rabbitmq"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
	"github.com/gin-gonic/gin"
)

func GetUsersRoutes(rbMQ *rabbitmq.AMQPConn) restRouter.Routes {
	userRoutes := restRouter.Routes{
		restRouter.Route{
			Method:  http.MethodPost,
			Pattern: "/api/users/signup",
			Handler: func(c *gin.Context) {
				AuthHandler(c, rbMQ, "signup")
			},
		},
		restRouter.Route{
			Method:  http.MethodPost,
			Pattern: "//api/users/login",
			Handler: func(c *gin.Context) {
				AuthHandler(c, rbMQ, "login")
			},
		},
	}

	return userRoutes
}
