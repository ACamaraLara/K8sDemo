package restServer

import (

	// amqp_mocks "msgbrokerlib/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	// rabbitmq "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// Init function executed before start tests to avoid verbose logs.
func init() {
	// Disables logger for unit testing.
	zerolog.SetGlobalLevel(zerolog.Disabled)
	// Disables gin REST connection verbose logs.
	gin.SetMode(gin.ReleaseMode)
}

// Simulates a POST request to the service. Should not fail.
func TestHandleAddParkingInfoNotFail(t *testing.T) {

	// testJsonInfo, err := os.ReadFile("testdata/testPOST.json")

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// var testInfoBuff bytes.Buffer

	// testInfoBuff.Write(testJsonInfo)

	// wr := httptest.NewRecorder()
	// req := httptest.NewRequest(http.MethodPost, "/SGORA", &testInfoBuff)

	// c, _ := gin.CreateTestContext(wr)
	// c.Request = req

	// conn := &amqp.AMQPConn{RbWrapper: &amqp_mocks.RabbitMock{Msgs: make(chan rabbitmq.Delivery, 1)}}
	// ParkingMeterInfoPublish(c, conn)
	// if wr.Code != http.StatusCreated {
	// 	t.Errorf("got HTTP status code %d, expected 201", wr.Code)
	// }
}

// Simulates a bad POST request. Should fail.
func TestHandleOublishParkingInfoFail(t *testing.T) {

	// testFailInfo := []byte("This is fail info")

	// var testInfoBuff bytes.Buffer

	// testInfoBuff.Write(testFailInfo)

	// wr := httptest.NewRecorder()
	// req := httptest.NewRequest(http.MethodPost, "/SGORA", &testInfoBuff)

	// c, _ := gin.CreateTestContext(wr)
	// c.Request = req

	// conn := &amqp.AMQPConn{RbWrapper: &amqp_mocks.RabbitMock{}}
	// ParkingMeterInfoPublish(c, conn)

	// if wr.Code != http.StatusUnprocessableEntity {
	// 	t.Fatalf("got HTTP status code %d, expected 422", wr.Code)
	// }
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
			`response body "%s" does not contain "Welcome to SGORA"`,
			wr.Body.String(),
		)
	}
}

// Simulates Bad requests to the different urls of the Service. Should return 404 not found.
func TestBadService(t *testing.T) {
	// For execute different subtests, create a new router.
	router := NewRouter()

	// Define different subtests with bad calls
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"TestBadMainCall", http.MethodPost, "/"},
		{"TestBadSgoraCallPut", http.MethodPut, "/SGORA"},
		{"TestBadSgoraCallPatch", http.MethodPatch, "/SGORA"},
		{"TestBadSgoraCallDelete", http.MethodDelete, "/SGORA"},
		{"TestBadIDCallPut", http.MethodPut, "/SGORA/123"},
		{"TestBadIDCallPatch", http.MethodPatch, "/SGORA/123"},
		{"TestBadIDCallDelete", http.MethodDelete, "/SGORA/123"},
		{"TestBadIDCallPost", http.MethodPost, "/SGORA/123"},
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
