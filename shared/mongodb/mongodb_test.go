package mongodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MongoTestSuite struct {
	suite.Suite
	mockClient     *MockClient
	mockCollection *MockCollection
}

// SetupTest is executed before every test in the suite
func (suite *MongoTestSuite) SetupTest() {
	suite.mockClient = new(MockClient)
	suite.mockCollection = new(MockCollection)
}

// TestCheckConnection tests the checkConnection function
func (suite *MongoTestSuite) TestCheckConnection() {
	suite.mockClient.On("Ping", mock.Anything, mock.Anything).Return(nil).Once()
	suite.mockClient.On("Ping", mock.Anything, mock.Anything).
		Return(fmt.Errorf("TestConnectionError")).Once()

	err := checkConnection(context.TODO(), suite.mockClient)

	assert.NoError(suite.T(), err)

	err = checkConnection(context.TODO(), suite.mockClient)

	assert.Contains(suite.T(), err.Error(), "TestConnectionError")

	suite.mockClient.AssertExpectations(suite.T())
}

// TestSetupCollections tests the setupCollections function
func (suite *MongoTestSuite) TestSetupCollections() {
	suite.mockClient.On("GetDBCollection", "testdb", "collection1").Return(suite.mockCollection)
	suite.mockClient.On("GetDBCollection", "testdb", "collection2").Return(suite.mockCollection)

	conf := &config.DBConfig{
		Host:        "localhost",
		Port:        "27017",
		DbName:      "testdb",
		Collections: []string{"collection1", "collection2"},
	}

	mongoDB := &MongoDB{
		Conf:        conf,
		Client:      suite.mockClient,
		Collections: make(map[string]Collection),
	}

	setupCollections(mongoDB, suite.mockClient, conf)

	suite.mockClient.AssertExpectations(suite.T())

	assert.Len(suite.T(), mongoDB.Collections, 2)
	assert.Contains(suite.T(), mongoDB.Collections, "collection1")
	assert.Contains(suite.T(), mongoDB.Collections, "collection2")
}

// Run the tests
func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
