package api

import (
	"account-service/internal/account"
	"account-service/internal/encryption"
	"account-service/internal/model"
	"time"

	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/jwtManager"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/restRouter"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthTestSuite struct {
	suite.Suite
	router         *gin.Engine
	mockClient     *mongodb.MockClient
	mockCollection *mongodb.MockCollection
	mongoDB        *mongodb.MongoDB
}

func (suite *AuthTestSuite) SetupTest() {
	// Initialize structures used in all tests
	gin.SetMode(gin.TestMode)

	suite.mockClient = &mongodb.MockClient{}
	suite.mockCollection = &mongodb.MockCollection{}
	collections := make(map[string]mongodb.Collection)
	collections["USERS"] = suite.mockCollection

	suite.mongoDB = &mongodb.MongoDB{Client: suite.mockClient, Collections: collections}
	accCtrl := account.NewAccountController(suite.mongoDB)
	jwtConfig := &jwtManager.Config{
		SecretKey:         "test-secret",
		AccessTokenExpiry: 1 * time.Hour,
		Issuer:            "test-issuer",
	}
	jwtMgr, _ := jwtManager.NewManager(jwtConfig)
	suite.router = restRouter.NewRouter(SetAccountRoutes(accCtrl, jwtMgr))
}

func (suite *AuthTestSuite) TestSignupHandler() {
	testUser := `{
		"firstName": "Test",
		"lastName": "User",
		"email": "test@example.com",
		"password": "password"
	}`

	suite.mockCollection.On("InsertOne", mock.Anything, mock.AnythingOfType("*model.User"), mock.Anything).Return(&mongo.InsertOneResult{}, nil)

	suite.Run("Valid Signup", func() {
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(testUser)))
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusCreated, resp.Code)
		assert.Contains(suite.T(), resp.Body.String(), "User registered successfully.")
	})

	suite.Run("User Already Registered", func() {
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).Return(mongo.NewSingleResultFromDocument(bson.D{{Key: "email", Value: "test@example.com"}, {Key: "username", Value: "testuser"}}, nil, nil)).Once()

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(testUser)))
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
		assert.Contains(suite.T(), resp.Body.String(), "account already registered")
	})

	suite.Run("Invalid JSON", func() {
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(`{invalid json`)))
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
		assert.Contains(suite.T(), resp.Body.String(), "Invalid signup request payload.")
	})
}

func (suite *AuthTestSuite) TestLoginHandler() {
	testUser := model.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}
	// Mock should return encrypted password to decrypt it and check hash.
	testUser.Password, _ = encryption.GetHash([]byte("password"))

	bsonUser, err := bson.Marshal(testUser)
	if err != nil {
		suite.T().Fatalf("Error converting to BSON: %v", err)
	}

	suite.Run("Valid Login", func() {
		loginRequest := `{"email":"test@example.com", "password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
		resp := httptest.NewRecorder()

		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(bsonUser, nil, nil)).Once()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusOK, resp.Code)
	})

	suite.Run("Invalid Login Payload", func() {
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{invalid json`)))
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
		assert.Contains(suite.T(), resp.Body.String(), "Invalid login request payload.")
	})

	suite.Run("User Not Found", func() {
		loginRequest := `{"email":"notfound@example.com", "password":"fakepass"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
		resp := httptest.NewRecorder()

		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "notfound@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
		assert.Contains(suite.T(), resp.Body.String(), "user not registered")
	})

	suite.Run("Incorrect Password", func() {
		loginRequest := `{"email":"test@example.com", "password":"wrongpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(loginRequest)))
		resp := httptest.NewRecorder()

		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(bsonUser, nil, nil)).Once()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
		assert.Contains(suite.T(), resp.Body.String(), "incorrect password")
	})
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
