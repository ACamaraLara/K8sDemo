package mongodb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Returns a url with the necessary format to connect to MongoDB.
func (mb *MongoConfig) GetURL() string {
	return "mongodb://" + mb.User + ":" + mb.Passwd + "@" + mb.Host + ":" + strconv.Itoa(mb.Port) + "/" + mb.DbName
}

// ConnectMongoClient connects the initialized object to the given database.
// @throws an error if something went wrong during the connection.
func (mb *MongoDBClient) ConnectMongoClient(ctx context.Context) error {

	clientOptions := options.Client().ApplyURI(mb.Config.GetURL())

	client, err := mb.DBWrapper.Connect(ctx, clientOptions)

	if err != nil {
		return fmt.Errorf("error connecting to MongoDB url: %v", err)
	}

	mb.Client = client

	if !mb.checkConnection(ctx) {
		return fmt.Errorf("mongodb connection check failed")
	}

	coll := mb.DBWrapper.GetDBCollection(mb)

	mb.Collection = coll

	log.Info().Msg("Connected to MongoDB on url " + mb.Config.GetURL())

	return nil
}

// DisconnectMongoClient disconnects this object from the database.
func (mb *MongoDBClient) DisconnectMongoClient(ctx context.Context) error {

	if err := mb.DBWrapper.Disconnect(mb, ctx); err != nil {
		return fmt.Errorf("error disconnecting from database %v", err)
	}

	// Free pointers to avoid undefined behavior.
	mb.Client = nil
	mb.Collection = nil

	return nil
}

// Function to check if the client is successfully connected to database.
func (mb *MongoDBClient) checkConnection(ctx context.Context) bool {

	if mb.Client == nil {
		return false
	}

	err := mb.DBWrapper.PingToDB(mb, ctx)
	if err != nil {
		log.Error().Msg("Error pinging mongoDB:" + err.Error())
	}

	return err == nil
}
