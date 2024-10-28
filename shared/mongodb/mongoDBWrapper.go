package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoWrapper interface {
	Connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error)
	GetDBCollection(mb *MongoDBClient) *mongo.Collection
	Disconnect(mb *MongoDBClient, ctx context.Context) error
	PingToDB(mb *MongoDBClient, ctx context.Context) error
	InsertData(mb *MongoDBClient, ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, mb *MongoDBClient, filter interface{}) *mongo.SingleResult
}

type MongoWrapper struct{}

func (mgW *MongoWrapper) Connect(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(context.Background(), opts)
}

func (mgW *MongoWrapper) GetDBCollection(mb *MongoDBClient) *mongo.Collection {
	return mb.Client.Database(mb.Config.DbName).Collection(mb.Config.Collection)
}

func (mgW *MongoWrapper) Disconnect(mb *MongoDBClient, ctx context.Context) error {
	return mb.Client.Disconnect(ctx)
}

func (mgW *MongoWrapper) PingToDB(mb *MongoDBClient, ctx context.Context) error {
	return mb.Client.Ping(ctx, nil)
}

func (mgW *MongoWrapper) InsertData(mb *MongoDBClient, ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return mb.Collection.InsertOne(ctx, document)
}

func (mgW *MongoWrapper) FindOne(ctx context.Context, mb *MongoDBClient, filter interface{}) *mongo.SingleResult {
	return mb.Collection.FindOne(ctx, filter)
}
