package account

import (
	"api-gateway/internal/api/service"

	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/gin-gonic/gin"
)

// AccountHandler handles user-related API routes
type AccountHandler struct {
	serviceRouter *service.ServiceRouter
}

func NewAccountHandler(serviceRouter *service.ServiceRouter) *AccountHandler {
	return &AccountHandler{
		serviceRouter: serviceRouter,
	}
}

// Signup handles user registration requests.
func (h *AccountHandler) Signup(ctx *gin.Context) {
	h.serviceRouter.SendToService(ctx, "account-service", "/signup")
}

// Login handles user login requests.
func (h *AccountHandler) Login(ctx *gin.Context) {
	h.serviceRouter.SendToService(ctx, "account-service", "/login")
}

// RegisterRoutes registers the account routes.
func (h *AccountHandler) RegisterRoutes() restRouter.Routes {
	return restRouter.Routes{
		{
			Method:  http.MethodPost,
			Pattern: "/signup",
			Handler: h.Signup,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/login",
			Handler: h.Login,
		},
	}
}
