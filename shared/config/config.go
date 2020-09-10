package config

import (
	"fmt"
	"os"
	"sync"

	Error "github.com/fuadajip/kafka-cluster-go/shared/error"
	"github.com/spf13/viper"
)

type (
	// ImmutableConfigInterface is an interface represent methods in config
	ImmutableConfigInterface interface {
		GetPort() int
		GetDatabaseHost() string
		GetDatabasePort() string
		GetDatabaseName() string
		GetDatabaseUser() string
		GetDatabasePassword() string
		GetAWSAccessKey() string
		GetAWSSecretKey() string
		GetAWSRegion() string
		GetAWSBucket() string
		GetTokenSecret() string
		GetTokenSecretStatic() string
		GetNATSHost() string
		GetRabbitHost() string
		GetRedisHost() string
		GetRedisName() string
		GetRedisPassword() string
	}

	// im is a struct to mapping the structure of related value model
	im struct {
		Port              int    `mapstructure:"PORT"`
		DatabaseHost      string `mapstructure:"DATABASE_HOST"`
		DatabasePort      string `mapstructure:"DATABASE_PORT"`
		DatabaseName      string `mapstructure:"DATABASE_NAME"`
		DatabaseUser      string `mapstructure:"DATABASE_USER"`
		DatabasePassword  string `mapstructure:"DATABASE_PASSWORD"`
		AWSAccessKey      string `mapstructure:"AWS_ACCESS_KEY"`
		AWSSecretKey      string `mapstructure:"AWS_SECRET_KEY"`
		AWSRegion         string `mapstructure:"AWS_REGION"`
		AWSBucket         string `mapstructure:"AWS_BUCKET"`
		TokenSecret       string `mapstructure:"SECRET"`
		TokenSecretStatic string `mapstructure:"SECRET_STATIC"`
		NATSHost          string `mapstructure:"NATS_HOST"`
		RabbitHost        string `mapstructure:"RABBIT_HOST"`
		RedisHost         string `mapstructure:"REDIS_HOST"`
		RedisName         string `mapstructure:"REDIS_NAME"`
		RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	}
)

func (i *im) GetPort() int {
	return i.Port
}

func (i *im) GetDatabaseHost() string {
	return i.DatabaseHost
}

func (i *im) GetDatabasePort() string {
	return i.DatabasePort
}

func (i *im) GetDatabaseName() string {
	return i.DatabaseName
}

func (i *im) GetDatabaseUser() string {
	return i.DatabaseUser
}

func (i *im) GetDatabasePassword() string {
	return i.DatabasePassword
}

func (i *im) GetAWSAccessKey() string {
	return i.AWSAccessKey
}

func (i *im) GetAWSSecretKey() string {
	return i.AWSSecretKey
}

func (i *im) GetAWSRegion() string {
	return i.AWSRegion
}

func (i *im) GetAWSBucket() string {
	return i.AWSBucket
}

func (i *im) GetTokenSecret() string {
	return i.TokenSecret
}
func (i *im) GetTokenSecretStatic() string {
	return i.TokenSecretStatic
}

func (i *im) GetRabbitHost() string {
	return i.RabbitHost
}

func (i *im) GetNATSHost() string {
	return i.NATSHost
}

func (i *im) GetRedisHost() string {
	return i.RedisHost
}

func (i *im) GetRedisName() string {
	return i.RedisName
}

func (i *im) GetRedisPassword() string {
	return i.RedisPassword
}

var (
	imOnce    sync.Once
	myEnv     map[string]string
	immutable im
)

// NewImmutableConfig is a factory that return of its config implementation
func NewImmutableConfig() ImmutableConfigInterface {
	imOnce.Do(func() {
		v := viper.New()
		appEnv, exists := os.LookupEnv("APP_ENV")
		fmt.Println(appEnv)
		if exists {
			if appEnv == "staging" {
				v.SetConfigName("app.config.staging")
			} else if appEnv == "development" {
				v.SetConfigName("app.config.dev")
			} else if appEnv == "production" {
				v.SetConfigName("app.config.prod")
			} else {
				v.SetConfigName("app.config.local")
			}
		} else {
			v.SetConfigName("app.config.local")
		}

		v.AddConfigPath(".")
		v.SetEnvPrefix("QUOTATION")
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			Error.Wrap(500, "[QUOTATION-SYS-001]", err, "[CONFIG][missing] Failed to read app.config.* file", "failed")
		}

		v.Unmarshal(&immutable)
	})

	return &immutable
}
