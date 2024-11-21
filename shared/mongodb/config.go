package mongodb

import (
	"flag"
	"strings"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
)

// Default values are meant for development and testing purposes.
const (
	DefaultMongoDBHost       = "mongodb"
	DefaultMongoDBPort       = "27017"
	DefaultMongoDBUserName   = "admin"
	DefaultMongoDBUserPass   = "admin_pass"
	DefaultMongoDBName       = "K8DEMO"
	DefaultMongoDBCollection = "USERS"
)

type MongoConfig struct {
	Host        string   // Database host address.
	Port        string   // Port where database pod is listen.
	DbName      string   // Name of the database.
	Collections []string // Name of the collection inside database.
	User        string   // Database password.
	Passwd      string   // Database password.
}

// AddFlagsParams modifies the MongoConfig struct based on command-line flags and environment variables.
func (cfg *MongoConfig) AddFlagsParams() {
	// Host, Port, DbName, User, Passwd flags remain the same
	flag.StringVar(&cfg.Host, "mongo-host", config.GetEnvironWithDefault("MONGODB_HOST", DefaultMongoDBHost), "MongoDB server host (MONGODB_HOST).")
	flag.StringVar(&cfg.Port, "mongo-port", config.GetEnvironWithDefault("MONGODB_PORT", DefaultMongoDBPort), "MongoDB server port (MONGODB_PORT).")
	flag.StringVar(&cfg.DbName, "mongo-db", config.GetEnvironWithDefault("MONGODB_DATABASE", DefaultMongoDBName), "MongoDB database name (MONGODB_DATABASE).")
	flag.StringVar(&cfg.User, "mongo-user", config.GetEnvironWithDefault("MONGODB_USER", DefaultMongoDBUserName), "MongoDB username (MONGODB_USER).")
	flag.StringVar(&cfg.Passwd, "mongo-passwd", config.GetEnvironWithDefault("MONGODB_PASSWD", DefaultMongoDBUserPass), "MongoDB password (MONGODB_PASSWD).")

	// For collections, we expect a comma-separated list of collection names.
	var collections string
	flag.StringVar(&collections, "mongo-collections", config.GetEnvironWithDefault("MONGODB_COLLECTIONS",
		DefaultMongoDBCollection), "Comma-separated MongoDB collection names (MONGODB_COLLECTIONS).")
	cfg.Collections = strings.Split(collections, ",")
}

// Returns a url with the necessary format to connect to MongoDB.
func (cfg *MongoConfig) GetURL() string {

	if cfg.User == "" {
		return "mongodb://" + cfg.Host + ":" + cfg.Port + "/" + cfg.DbName
	}
	return "mongodb://" + cfg.User + ":" + cfg.Passwd + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DbName
}
