package transporter

import (
	"github.com/fuadajip/kafka-cluster-go/shared/config"
	"github.com/fuadajip/kafka-cluster-go/shared/constant"
	"github.com/fuadajip/kafka-cluster-go/shared/log"
	"github.com/nats-io/nats.go"
)

var (
	logger = log.NewServiceLog(constant.ServiceName)
)

type (

	// NatsInterface is an interface that represent nats methods implementation
	NatsInterface interface {
		OpenNatsConn() (*nats.EncodedConn, error)
	}

	serviceNats struct {
		SharedConfig config.ImmutableConfigInterface
	}
)

// OpenNatsConn will return encoded nats connection
func (r *serviceNats) OpenNatsConn() (*nats.EncodedConn, error) {
	logger.Info("Start open nats connection...")
	nc, err := nats.Connect(r.SharedConfig.GetNATSHost())
	if err != nil {
		return nil, err
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewNats is a factory that implement of nats configuration
func NewNats(config config.ImmutableConfigInterface) NatsInterface {
	if config == nil {
		panic("[CONFIG] immutable config is required")
	}

	return &serviceNats{config}
}
