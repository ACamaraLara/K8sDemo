package service

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProxyRequest forwards the request to the target service and copies the response back to the client.
func ProxyRequest(ctx *gin.Context, targetURL string) {
	proxyRequest, err := http.NewRequest(ctx.Request.Method, targetURL, ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	proxyRequest.Header = ctx.Request.Header

	client := &http.Client{}
	response, err := client.Do(proxyRequest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to communicate with service:" + err.Error()})
		return
	}
	defer response.Body.Close()

	copyResponse(ctx, response)
}

func copyResponse(ctx *gin.Context, response *http.Response) {
	ctx.Writer.WriteHeader(response.StatusCode)

	for key, values := range response.Header {
		for _, value := range values {
			ctx.Header(key, value)
		}
	}

	io.Copy(ctx.Writer, response.Body)
}
