package database

import (
	"context"
)

type dataBase interface {
	InsertOne(ctx context.Context, tableName string, document interface{}) error
	FindOne(ctx context.Context, tablename string, document interface{}, filters ...interface{}) error
}

type DBManager struct {
	Db dataBase
}

func (dbm *DBManager) InsertOne(ctx context.Context, tableName string, document interface{}) error {
	return dbm.Db.InsertOne(ctx, tableName, document)
}

func (dbm *DBManager) FindOne(ctx context.Context, tableName string, document interface{}, filters ...interface{}) error {
	return dbm.Db.FindOne(ctx, tableName, document, filters...)
}
