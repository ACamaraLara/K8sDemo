package mongodb

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoConfig struct {
	Host       string // Database host address.
	Port       int    // Port where database pod is listen.
	DbName     string // Name of the database.
	Collection string // Name of the collection inside database.
	User       string // Database password.
	Passwd     string // Database password.
}

// Default values are meant for development and testing purposes.
const (
	DefaultMongoDBHost       = "localhost"
	DefaultMongoDBPort       = 27017
	DefaultMongoDBUserName   = "smge"
	DefaultMongoDBUserPass   = "smge"
	DefaultMongoDBName       = "SMGE"
	DefaultMongoDBCollection = "SGORA"
)

// Struct that stores a connection to MongoDataBase.
type MongoDBClient struct {
	Client     *mongo.Client     // Mongo object to manage connexion to database.
	Config     *MongoConfig      // Stores database information.
	Collection *mongo.Collection // Mongo object to manage collection.
	DBWrapper  IMongoWrapper     // Mongo db wrapper to abstract calls to database API.
}

func NewMongoDBClient(cfg MongoConfig) *MongoDBClient {
	return &MongoDBClient{Config: &cfg, DBWrapper: &MongoWrapper{}}
}

func AddFlagsParams(cfg *MongoConfig) {
	flag.StringVar(&cfg.Host, "mongo-host", config.GetEnvironWithDefault("MONGODB_HOST", DefaultMongoDBHost), "MongoDB server host (MONGODB_HOST).")
	flag.IntVar(&cfg.Port, "mongo-port", config.GetEnvironIntWithDefault("MONGODB__PORT", DefaultMongoDBPort), "MongoDB server port (MONGODB__PORT).")
	flag.StringVar(&cfg.DbName, "mongo-db", config.GetEnvironWithDefault("MONGODB_DATABASE", DefaultMongoDBName), "MongoDB database name (MONGODB_DATABASE).")
	flag.StringVar(&cfg.Collection, "mongo-collection", config.GetEnvironWithDefault("MONGODB_COLLECTION", DefaultMongoDBCollection), "MongoDB collection (MONGODB_COLLECTION).")
	flag.StringVar(&cfg.User, "mongo-user", config.GetEnvironWithDefault("MONGODB_USER", DefaultMongoDBUserName), "MongoDB username (MONGODB_USER).")
	flag.StringVar(&cfg.Passwd, "mongo-passwd", config.GetEnvironWithDefault("MONGODB_PASSWD", DefaultMongoDBUserPass), "MongoDB password (MONGODB_PASSWD).")
}
