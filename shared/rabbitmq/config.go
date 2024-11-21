package rabbitmq

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
)

const (
	DefaultRabbitHost = "rabbitmq"
	DefaultRabbitPort = "5672"
	DefaultRabbitUser = "admin"
	DefaultRabbitPass = "admin_pass" // Proof of concept. This wouldn't be here in production projects.
)

// RabbitConfig holds the configuration details required to connect to a RabbitMQ server.
type RabbitConfig struct {
	Host   string // Address of the RabbitMQ server.
	Port   string // Port number on which the RabbitMQ server is listening.
	User   string // Username used for authenticating to the RabbitMQ server.
	Passwd string // Password used for authenticating to the RabbitMQ server.
}

// Init RabbitConfig structure with environment variables or default variable.
// If an app input param is received while starting the application, this param will be used
func (cfg *RabbitConfig) AddFlagsParams() {
	flag.StringVar(&cfg.Host, "rabbit-host", config.GetEnvironWithDefault("RABBITMQ_HOST", DefaultRabbitHost), "RabbitMQ broker address (RABBITMQ_HOST).")
	flag.StringVar(&cfg.Port, "rabbit-port", config.GetEnvironWithDefault("RABBITMQ__PORT", DefaultRabbitPort), "RabbitMQ broker port (RABBITMQ__PORT).")
	flag.StringVar(&cfg.User, "rabbit-user", config.GetEnvironWithDefault("USERNAME", DefaultRabbitUser), "User to connect to RabbitMQ broker (RABBITMQ_USER).")
	flag.StringVar(&cfg.Passwd, "rabbit-passwd", config.GetEnvironWithDefault("PASSWORD", DefaultRabbitPass), "RabbitMQ password (RABBITMQ_PASSWD).")
}

// Returns a url with the necessary format to connect to RabbitMQ broker.
func (cfg *RabbitConfig) GetURL() string {
	if cfg.User == "" {
		return "amqp://" + cfg.Host + ":" + cfg.Port + "/"
	}
	return "amqp://" + cfg.User + ":" + cfg.Passwd + "@" + cfg.Host + ":" + cfg.Port + "/"
}
