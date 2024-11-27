package service

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ServiceRouter manages communication between the api-gateway and other microservices inside the cluster.
type ServiceRouter struct {
	serviceURLs map[string]string
}

func NewServiceRouter() *ServiceRouter {
	// Initialize with the base URLs for all services
	return &ServiceRouter{
		serviceURLs: map[string]string{
			"account-service": "http://account-service-svc.microservices.svc.cluster.local",
			// Add future services URLs here.
		},
	}
}

// SendToService forwards the request to the target service and copies the response back to the client.
func (sr *ServiceRouter) SendToService(ctx *gin.Context, serviceName, path string) {
	log.Info().Str("Method", ctx.Request.Method).Str("Path", path).Msg("Received new request.")
	baseURL, ok := sr.serviceURLs[serviceName]
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Service not found"})
		return
	}

	targetURL := sr.buildTargetURL(baseURL, path)

	ProxyRequest(ctx, targetURL)
}

func (sr *ServiceRouter) buildTargetURL(baseURL, path string) string {
	targetURL, err := url.Parse(baseURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse service URL")
		return ""
	}

	targetURL.Path = path
	return targetURL.String()
}
