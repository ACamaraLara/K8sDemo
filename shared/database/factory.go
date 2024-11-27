package database

import (
	"context"
	"fmt"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database/config"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
)

func NewDatabase(ctx context.Context, dbType string, config *config.DBConfig) (*DBManager, error) {
	switch dbType {
	case "mongo":
		mongoDb, err := mongodb.NewMongoDBClient(ctx, config)
		if err != nil {
			return nil, err
		}
		return &DBManager{Db: mongoDb}, nil
	default:
		// If no match, an error is returned
		return nil, fmt.Errorf("unsupported config type: %T", config)
	}
}
