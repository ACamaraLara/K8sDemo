package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBClient struct {
	client *mongo.Client
}

func NewClient(ctx context.Context, uri string, opts ...*options.ClientOptions) (Client, error) {
	// Use provided options or default options
	clientOpts := options.Client().ApplyURI(uri)

	// ToDo: Abstract this functions to make this testeable.
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &DBClient{client: client}, nil
}

func (m *DBClient) GetDBCollection(dbName, collectionName string) Collection {
	db := m.client.Database(dbName)
	return &DBCollection{collection: db.Collection(collectionName)}
}

func (m *DBClient) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *DBClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.client.Ping(ctx, rp)
}
