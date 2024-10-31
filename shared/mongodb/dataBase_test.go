package mongodb

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// Init function executed before start tests to avoid verbose logs.
func init() {
	// Disables logger for unit testing.
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// Tests that connection func has the expected behavior.
func TestConnectionNotFail(t *testing.T) {
	mockDBWrapper := new(MongoMock)
	mongoDB := &MongoDBClient{Config: &MongoConfig{}, DBWrapper: mockDBWrapper}

	// Define expected behavior for Connect method in the mock
	mockDBWrapper.On("Connect", mock.Anything, mock.Anything).Return(&mongo.Client{}, nil)
	mockDBWrapper.On("GetDBCollection", mongoDB).Return(&mongo.Collection{})
	mockDBWrapper.On("PingToDB", mongoDB, mock.Anything).Return(nil)

	if err := mongoDB.ConnectMongoClient(context.TODO()); err != nil {
		t.Error("Expected none error but one given", err)
	}

	if mongoDB.Client == nil {
		t.Error("Connection shouldn't be a null pointer.")
	}

	if mongoDB.Collection == nil {
		t.Error("Channel shouldn't be a null pointer.")
	}
}

// Tests that close connection func has the expected behavior.
func TestCloseConnectionNotFail(t *testing.T) {
	mockDBWrapper := new(MongoMock)
	mongoDB := &MongoDBClient{Config: &MongoConfig{}, DBWrapper: mockDBWrapper}

	mockDBWrapper.On("Disconnect", mongoDB, mock.Anything).Return(nil)

	if err := mongoDB.DisconnectMongoClient(context.TODO()); err != nil {
		t.Error("Expected none error but one given", err)
	}

	if mongoDB.Client != nil {
		t.Error("Client should be a null pointer after Disconnect client.")
	}

	if mongoDB.Collection != nil {
		t.Error("Collection should be a null pointer.")
	}
}
