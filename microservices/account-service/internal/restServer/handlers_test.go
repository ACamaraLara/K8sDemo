package restServer

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/dataTypes"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var router *gin.Engine
var mockDBWrapper *mongodb.MongoMock
var mongoDBClient *mongodb.MongoDBClient

func init() {
	// Init structures that are used by all tests like router, routes and mongodbClient.
	gin.SetMode(gin.TestMode)
	mockDBWrapper = new(mongodb.MongoMock)
	mongoDBClient = &mongodb.MongoDBClient{Config: &mongodb.MongoConfig{}, DBWrapper: mockDBWrapper}
	InitRestRoutes(mongoDBClient)
	router = restRouter.NewRouter()
}

func TestSignupHandler(t *testing.T) {

	testUser := `{
		"firstName": "Test",
		"lastName": "User",
		"email": "test@example.com",
		"password": "password"
	}`

	mockDBWrapper.On("InsertData", mock.Anything, mock.Anything, mock.AnythingOfType("dataTypes.User")).Return(&mongo.InsertOneResult{}, nil)

	t.Run("Valid Signup", func(t *testing.T) {
		mockDBWrapper.On("FindOne", mock.Anything, mongoDBClient, bson.M{"email": "test@example.com"}).Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(testUser)))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Contains(t, resp.Body.String(), "User registered successfully.")
	})

	t.Run("User Already Registered", func(t *testing.T) {
		mockDBWrapper.On("FindOne", mock.Anything, mongoDBClient, bson.M{"email": "test@example.com"}).Return(mongo.NewSingleResultFromDocument(bson.D{{Key: "email", Value: "test@example.com"}, {Key: "username", Value: "testuser"}}, nil, nil)).Once()

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(testUser)))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Contains(t, resp.Body.String(), "User already registered.")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(`{invalid json`)))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid payload.")
	})
}

func TestLoginHandler(t *testing.T) {

	testUser := dataTypes.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  getHash([]byte("password"))}

	bsonUser, err := bson.Marshal(testUser)
	if err != nil {
		t.Fatalf("Error converting to BSON: %v", err)
	}

	t.Run("Valid Login", func(t *testing.T) {
		loginRequest := `{"email":"test@example.com", "password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
		resp := httptest.NewRecorder()

		mockDBWrapper.On("FindOne", mock.Anything, mongoDBClient, bson.M{"email": "test@example.com"}).Return(mongo.NewSingleResultFromDocument(bsonUser, nil, nil)).Once()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "authToken")
	})

	t.Run("Invalid Login Payload", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{invalid json`)))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid payload.")
	})

	t.Run("User Not Found", func(t *testing.T) {
		loginRequest := `{"email":"notfound@example.com", "password":"fakepass"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
		resp := httptest.NewRecorder()

		mockDBWrapper.On("FindOne", mock.Anything, mongoDBClient, bson.M{"email": "notfound@example.com"}).Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Incorrect email or password.")
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		loginRequest := `{"email":"test@example.com", "password":"wrongpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
		resp := httptest.NewRecorder()

		mockDBWrapper.On("FindOne", mock.Anything, mongoDBClient, bson.M{"email": "test@example.com"}).Return(mongo.NewSingleResultFromDocument(bsonUser, nil, nil)).Once()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Incorrect email or password.")
	})
}
