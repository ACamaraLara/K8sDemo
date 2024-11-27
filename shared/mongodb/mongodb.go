package mongodb

import (
	"context"
	"fmt"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database/config"
	"github.com/rs/zerolog/log"
)

type MongoDB struct {
	Conf        *config.DBConfig
	Client      Client
	Collections map[string]Collection
}

func NewMongoDBClient(ctx context.Context, conf *config.DBConfig) (*MongoDB, error) {
	log.Info().Msg("Connecting to mongodb..." + conf.GetURL())
	client, err := createMongoClient(ctx, conf)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is successful.
	err = checkConnection(ctx, client)
	if err != nil {
		return nil, err
	}

	mongoDB := &MongoDB{
		Conf:        conf,
		Client:      client,
		Collections: make(map[string]Collection),
	}

	setupCollections(mongoDB, client, conf)

	return mongoDB, nil
}

func createMongoClient(ctx context.Context, conf *config.DBConfig) (Client, error) {
	client, err := NewClient(ctx, conf.GetURL())
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize MongoDB client")
		return nil, fmt.Errorf("failed to initialize MongoDB client: %w", err)
	}
	log.Info().Msg("MongoDB client initialized successfully")
	return client, nil
}

func checkConnection(ctx context.Context, client Client) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ping MongoDB")
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	log.Info().Msg("MongoDB is reachable")
	return nil
}

func setupCollections(mongoDB *MongoDB, client Client, conf *config.DBConfig) {
	log.Info().Msg("Setting up collections...")

	for _, collectionName := range conf.Collections {
		collection := client.GetDBCollection(conf.DbName, collectionName)
		mongoDB.Collections[collectionName] = collection
		log.Info().Msgf("Added collection: %s", collectionName)
	}

	log.Info().Msgf("Initialized %d collections", len(mongoDB.Collections))
}

// Implementation of database interface functions.

func (mdb *MongoDB) InsertOne(ctx context.Context, tableName string, document interface{}) error {
	// Don't use the InsertOneResult value because inserted object will already have the _id value set.
	_, err := mdb.Collections[tableName].InsertOne(ctx, document)
	return err
}

func (mdb *MongoDB) FindOne(ctx context.Context, tableName string, document interface{}, filters ...interface{}) error {
	// Process filters before perform the query to Mongodb.
	return mdb.Collections[tableName].FindOne(ctx, processFilters(filters...)).Decode(document)
}
