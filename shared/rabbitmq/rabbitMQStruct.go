package rabbitmq

import (
	"flag"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultRabbitHost = "rabbitmq"
	DefaultRabbitPort = 5672
	DefaultRabbitUser = "admin"
	DefaultRabbitPass = "admin_pass" // Prove of concept. This wouldn't be here in production projects.
)

// RabbitConfig holds the configuration details required to connect to a RabbitMQ server.
type RabbitConfig struct {
	Host   string // Address of the RabbitMQ server.
	Port   int    // Port number on which the RabbitMQ server is listening.
	User   string // Username used for authenticating to the RabbitMQ server.
	Passwd string // Password used for authenticating to the RabbitMQ server.
}

// Struct that stores a connection to RabbitMQ broker.
type AMQPConn struct {
	Address   string
	Port      int
	User      string
	Passwd    string
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	RbWrapper IRabbitWrapper
	Queues    []*amqp.Queue
}

func NewAMQPConn(cfg RabbitConfig) *AMQPConn {
	return &AMQPConn{
		Address:   cfg.Host,
		Port:      cfg.Port,
		User:      cfg.User,
		Passwd:    cfg.Passwd,
		Conn:      nil,
		Channel:   nil,
		RbWrapper: &RabbitWrapper{},
		Queues:    nil,
	}
}

func (cfg *RabbitConfig) AddFlagsParams() {
	flag.StringVar(&cfg.Host, "rabbit-host", config.GetEnvironWithDefault("RABBITMQ_HOST", DefaultRabbitHost), "RabbitMQ broker address (RABBITMQ_HOST).")
	flag.IntVar(&cfg.Port, "rabbit-port", config.GetEnvironIntWithDefault("RABBITMQ__PORT", DefaultRabbitPort), "RabbitMQ broker port (RABBITMQ__PORT).")
	flag.StringVar(&cfg.User, "rabbit-user", config.GetEnvironWithDefault("USERNAME", DefaultRabbitUser), "User to connect to RabbitMQ broker (RABBITMQ_USER).")
	flag.StringVar(&cfg.Passwd, "rabbit-passwd", config.GetEnvironWithDefault("PASSWORD", DefaultRabbitPass), "RabbitMQ password (RABBITMQ_PASSWD).")
}
