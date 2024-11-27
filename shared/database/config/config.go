package config

import (
	"flag"
	"strings"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/utils"
)

// Default values are meant for development and testing purposes. Generally the project will use MongoDB.
const (
	DefaultMongoDBHost       = "mongodb"
	DefaultMongoDBPort       = "27017"
	DefaultMongoDBUserName   = "admin"
	DefaultMongoDBUserPass   = "admin_pass"
	DefaultMongoDBName       = "K8DEMO"
	DefaultMongoDBCollection = "USERS"
)

type DBConfig struct {
	Host        string   // Database host address.
	Port        string   // Port where database pod is listen.
	DbName      string   // Name of the database.
	Collections []string // Name of the collection inside database.
	User        string   // Database password.
	Passwd      string   // Database password.
}

// AddFlagsParams modifies the MongoConfig struct based on command-line flags and environment variables.
func (cfg *DBConfig) AddFlagsParams() {
	// Host, Port, DbName, User, Passwd flags remain the same
	flag.StringVar(&cfg.Host, "db-host", utils.GetEnvironWithDefault("DB_HOST", DefaultMongoDBHost), "DB server host (DB_HOST).")
	flag.StringVar(&cfg.Port, "db-port", utils.GetEnvironWithDefault("DB_PORT", DefaultMongoDBPort), "MongoDB server port (DB_PORT).")
	flag.StringVar(&cfg.DbName, "db-db", utils.GetEnvironWithDefault("DB_DATABASE", DefaultMongoDBName), "MongoDB database name (DB_DATABASE).")
	flag.StringVar(&cfg.User, "db-user", utils.GetEnvironWithDefault("DB_USER", DefaultMongoDBUserName), "MongoDB username (DB_USER).")
	flag.StringVar(&cfg.Passwd, "db-passwd", utils.GetEnvironWithDefault("DB_PASSWD", DefaultMongoDBUserPass), "MongoDB password (DB_PASSWD).")

	// For collections, we expect a comma-separated list of collection names.
	var collections string
	flag.StringVar(&collections, "db-tables", utils.GetEnvironWithDefault("DB_TABLES",
		DefaultMongoDBCollection), "Comma-separated database tables (collections for Mongodb) names (DB_TABLES).")
	cfg.Collections = strings.Split(collections, ",")
}

// Returns a url with the necessary format to connect to MongoDB.
func (cfg *DBConfig) GetURL() string {

	if cfg.User == "" {
		return "mongodb://" + cfg.Host + ":" + cfg.Port + "/" + cfg.DbName
	}
	return "mongodb://" + cfg.User + ":" + cfg.Passwd + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DbName
}
