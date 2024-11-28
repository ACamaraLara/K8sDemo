package account_test

import (
	"account-service/internal/account"
	"account-service/internal/encryption"
	"account-service/internal/model"
	"context"
	"errors"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthTestSuite struct {
	suite.Suite
	mockClient     *mongodb.MockClient
	mockCollection *mongodb.MockCollection
	mongoDB        *mongodb.MongoDB
	accountCtrl    *account.AccountController
}

func (suite *AuthTestSuite) SetupTest() {
	// Mock MongoDB client and collection setup
	suite.mockClient = &mongodb.MockClient{}
	suite.mockCollection = &mongodb.MockCollection{}
	collections := make(map[string]mongodb.Collection)
	collections["USERS"] = suite.mockCollection

	suite.mongoDB = &mongodb.MongoDB{Client: suite.mockClient, Collections: collections}
	suite.accountCtrl = account.NewAccountController(suite.mongoDB)
}

func (suite *AuthTestSuite) TestSignupHandler() {
	// Define a mock user
	user := &model.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "password123",
	}

	suite.Run("Valid Signup", func() {
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"},
			mock.Anything).Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()
		suite.mockCollection.On("InsertOne", mock.Anything, mock.AnythingOfType("*model.User"),
			mock.Anything).Return(&mongo.InsertOneResult{}, nil).Once()

		err := suite.accountCtrl.Signup(context.Background(), user)
		suite.NoError(err)
		suite.mockCollection.AssertExpectations(suite.T())
	})

	suite.Run("User already exists", func() {
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(bson.D{{Key: "email", Value: "test@example.com"}, {Key: "username", Value: "testuser"}}, nil, nil)).Once()

		err := suite.accountCtrl.Signup(context.Background(), user)

		suite.Error(err)
		suite.Contains(err.Error(), "account already registered")
		suite.mockCollection.AssertExpectations(suite.T())
	})

	suite.Run("Error inserting user", func() {
		insertError := errors.New("error inserting user")
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"},
			mock.Anything).Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()
		suite.mockCollection.On("InsertOne", mock.Anything, mock.AnythingOfType("*model.User"),
			mock.Anything).Return(&mongo.InsertOneResult{}, insertError).Once()

		err := suite.accountCtrl.Signup(context.Background(), user)

		suite.Error(err)
		suite.Equal(err, insertError)
		suite.mockCollection.AssertExpectations(suite.T())
	})
}

func (suite *AuthTestSuite) TestLoginHandler() {
	// Define a mock user
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
		// Mock the behavior of findUserByEmail (returns a user)
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(bsonUser, nil, nil)).Once()
		user, err := suite.accountCtrl.Login(context.Background(), testUser.Email, "password")

		suite.NoError(err)
		suite.Equal(user, &testUser)
		suite.mockCollection.AssertExpectations(suite.T())
	})

	suite.Run("User Not Found", func() {
		// Mock the behavior of findUserByEmail (returns a user)
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(nil, nil, nil)).Once()
		_, err := suite.accountCtrl.Login(context.Background(), testUser.Email, "password")

		suite.Error(err)
		suite.Equal(err.Error(), "user not registered")
		suite.mockCollection.AssertExpectations(suite.T())
	})

	suite.Run("Incorrect password", func() {
		suite.mockCollection.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}, mock.Anything).
			Return(mongo.NewSingleResultFromDocument(bsonUser, nil, nil)).Once()
		_, err := suite.accountCtrl.Login(context.Background(), testUser.Email, "wrongPassword") // Set wrong password.

		suite.Error(err)
		suite.Contains(err.Error(), "incorrect password")
		suite.mockCollection.AssertExpectations(suite.T())
	})

}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
