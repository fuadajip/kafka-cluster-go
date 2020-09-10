package transporter

import (
	"github.com/fuadajip/kafka-cluster-go/shared/config"
	"github.com/streadway/amqp"
)

type (

	// RabbitmqInterface is an interface that represent rabbitmq methods and its implementation
	RabbitmqInterface interface {
		OpenRabbitConn() (*amqp.Connection, error)
	}

	// rabbit is a struct to map give struct
	rabbit struct {
		SharedConfig config.ImmutableConfigInterface
	}
)

// OpenRabbitConn is a method that hanalde implementation of rabbitmq connection
func (r *rabbit) OpenRabbitConn() (*amqp.Connection, error) {
	logger.Info("Start open rabbitmq connection...")
	conn, err := amqp.Dial(r.SharedConfig.GetRabbitHost())
	if err != nil {
		return nil, err
	}

	return conn, err
}

// NewRabbit is a factory that implement of rabbitmq configuration
func NewRabbit(config config.ImmutableConfigInterface) RabbitmqInterface {
	if config == nil {
		panic("[CONFIG] immutable config is required")
	}

	return &rabbit{config}
}
