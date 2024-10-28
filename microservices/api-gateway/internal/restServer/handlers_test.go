package restServer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Init function executed before start tests to avoid verbose logs.
func init() {
	// Disables logger for unit testing.
	zerolog.SetGlobalLevel(zerolog.Disabled)
	// Disables gin REST connection verbose logs.
	gin.SetMode(gin.ReleaseMode)
}

// Simulates a GET request to Main Handler. Should not fail.
func TestHandleMainCallNotFail(t *testing.T) {

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	c, _ := gin.CreateTestContext(wr)
	c.Request = req

	Main(c)
	if wr.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", wr.Code)
	}

	if !strings.Contains(wr.Body.String(), "Welcome to Kubernetes Blockchain") {
		t.Errorf(
			`response body "%s" does not contain "Welcome to Kubernetes Blockchain"`,
			wr.Body.String(),
		)
	}
}

// Simulates Bad requests to the different urls of the Service. Should return 404 not found.
func TestBadService(t *testing.T) {
	// For execute different subtests, create a new router.
	router := restRouter.NewRouter()

	// Define different subtests with bad calls
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"TestBadMainCall", http.MethodPost, "/"},
		{"TestBadSgoraCallPut", http.MethodPut, "/K8DEMO"},
		{"TestBadSgoraCallPatch", http.MethodPatch, "/K8DEMO"},
		{"TestBadSgoraCallDelete", http.MethodDelete, "/K8DEMO"},
		{"TestBadIDCallPut", http.MethodPut, "/K8DEMO/123"},
		{"TestBadIDCallPatch", http.MethodPatch, "/"},
		{"TestBadIDCallDelete", http.MethodDelete, "/"},
		{"TestBadIDCallPost", http.MethodPost, "/"},
	}

	// Iterate over the subtests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("Error creating Request: %v", err)
			}

			// Create a ResponseRecorder object to catch service response.
			resp := httptest.NewRecorder()

			// Execute request.
			router.ServeHTTP(resp, req)

			// Check response.
			if resp.Code != http.StatusNotFound {
				t.Errorf("Error in %s: subtest. Expected %d but received %d",
					tt.name, http.StatusBadRequest, resp.Code)
			}
		})
	}
}
