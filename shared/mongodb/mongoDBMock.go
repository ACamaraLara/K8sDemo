package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMock struct{}

func (mgW *MongoMock) Connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	return &mongo.Client{}, nil
}

func (mgW *MongoMock) GetDBCollection(mb *MongoDBClient) *mongo.Collection {
	return &mongo.Collection{}
}

func (mgW *MongoMock) Disconnect(mb *MongoDBClient, ctx context.Context) error {
	return nil
}

func (mgW *MongoMock) PingToDB(mb *MongoDBClient, ctx context.Context) error {
	return nil
}

func (mgW *MongoMock) InsertData(mb *MongoDBClient, ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return &mongo.InsertOneResult{}, nil
}

func (mgW *MongoMock) FindOne(ctx context.Context, mb *MongoDBClient, filter interface{}) *mongo.SingleResult {
	return &mongo.SingleResult{}
}
