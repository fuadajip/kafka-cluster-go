package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fuadajip/kafka-cluster-go/domain/healthcheck"
	errors "github.com/fuadajip/kafka-cluster-go/shared/error"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type repoHandler struct {
	redisSess *redis.Client
	mysqlSess *gorm.DB
	sqsSess   *sqs.SQS
}

// NewHealthCheckRepository returns implementation of methods in auth.Repository
func NewHealthCheckRepository(redisSess *redis.Client, mysqlSess *gorm.DB, sqsSess *sqs.SQS) healthcheck.Repository {
	return &repoHandler{
		redisSess: redisSess,
		mysqlSess: mysqlSess,
		sqsSess:   sqsSess,
	}
}

// MysqlHealthCheck is method that implement healthcheck.Repository
func (r repoHandler) MysqlHealthCheck() (bool, error) {

	if r.mysqlSess == nil {
		return false, errors.New("INVALID_MYSQL_SESSION")
	}

	err := r.mysqlSess.DB().Ping()
	if err != nil {
		return false, err
	}

	return true, nil

}

// RedisHealthCheck is method that implement healthcheck.Repository
func (r repoHandler) RedisHealthCheck() (bool, error) {

	if r.redisSess == nil {
		return false, errors.New("INVALID_REDIS_SESSION")
	}
	_, err := r.redisSess.Ping().Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r repoHandler) SqsHealthCheck(queueName string) (bool, error) {

	_, err := r.sqsSess.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
