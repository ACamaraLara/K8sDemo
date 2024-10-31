package mongodb

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMock struct {
	mock.Mock
}

func (mgW *MongoMock) Connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	args := mgW.Called(ctx, opts)
	return args.Get(0).(*mongo.Client), args.Error(1)
}

func (mgW *MongoMock) GetDBCollection(mb *MongoDBClient) *mongo.Collection {
	args := mgW.Called(mb)
	return args.Get(0).(*mongo.Collection)
}

func (mgW *MongoMock) Disconnect(mb *MongoDBClient, ctx context.Context) error {
	args := mgW.Called(mb, ctx)
	return args.Error(0)
}

func (mgW *MongoMock) PingToDB(mb *MongoDBClient, ctx context.Context) error {
	args := mgW.Called(mb, ctx)
	return args.Error(0)
}

func (mgW *MongoMock) InsertData(mb *MongoDBClient, ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := mgW.Called(mb, ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (mgW *MongoMock) FindOne(ctx context.Context, mb *MongoDBClient, filter interface{}) *mongo.SingleResult {
	args := mgW.Called(ctx, mb, filter)
	return args.Get(0).(*mongo.SingleResult)
}
