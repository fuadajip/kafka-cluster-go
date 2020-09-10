package container

import (
	"github.com/fgrosse/goldi"
	Config "github.com/fuadajip/kafka-cluster-go/shared/config"
	Database "github.com/fuadajip/kafka-cluster-go/shared/database"
	Transporter "github.com/fuadajip/kafka-cluster-go/shared/transporter"
)

// DefaultContainer returns default given depedency injections
func DefaultContainer() *goldi.Container {

	registry := goldi.NewTypeRegistry()

	config := make(map[string]interface{})
	container := goldi.NewContainer(registry, config)

	container.RegisterType("shared.config", Config.NewImmutableConfig)
	container.RegisterType("shared.redis", Database.NewRedis, "@shared.config")
	container.RegisterType("shared.mysql", Database.NewMysql, "@shared.config")
	container.RegisterType("shared.rabbit", Transporter.NewRabbit, "@shared.config")
	container.RegisterType("shared.nats", Transporter.NewNats, "@shared.config")
	container.RegisterType("shared.sqs", Transporter.NewSQS, "@shared.config")
	container.RegisterType("shared.dynamo", Database.NewDynamo, "@shared.config")

	return container
}
