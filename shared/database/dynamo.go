package database

import (
	"github.com/fuadajip/kafka-cluster-go/shared/config"
	"github.com/fuadajip/kafka-cluster-go/shared/util"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (

	// DynamoInterface is an interface that represent dynamodb in package database
	DynamoInterface interface {
		OpenDynamoDBConn() (*dynamodb.DynamoDB, error)
	}
)

func (d *database) OpenDynamoDBConn() (*dynamodb.DynamoDB, error) {
	logger.Info("Start open dynamo connection...")

	sess, err := util.GetAWSSession(d.SharedConfig)

	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	return svc, nil
}

// NewDynamo is an factory that implemeent of dynamo database configuration
func NewDynamo(config config.ImmutableConfigInterface) DynamoInterface {
	if config == nil {
		panic("[CONFIG] immutable config is rerquired")
	}

	return &database{config}
}
