package transporter

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fuadajip/kafka-cluster-go/shared/config"
	"github.com/fuadajip/kafka-cluster-go/shared/util"
)

type (
	// SQSInterface is an interface that represent sqs methods implementation
	SQSInterface interface {
		CreateSQSService() (svc *sqs.SQS, err error)
	}

	serviceSQS struct {
		SharedConfig config.ImmutableConfigInterface
	}
)

// NewSQS is a factory that implement of sqs configuration
func NewSQS(config config.ImmutableConfigInterface) SQSInterface {
	if config == nil {
		panic("[CONFIG] immutable config aws sqs required")
	}

	return &serviceSQS{config}
}

// CreateSQSService will return sqs service session after created
func (r *serviceSQS) CreateSQSService() (svc *sqs.SQS, err error) {
	logger.Info("Start open sqs session service...")
	sess, err := util.GetAWSSession(r.SharedConfig)

	if err != nil {
		return nil, err
	}
	service := sqs.New(sess)
	return service, nil
}
