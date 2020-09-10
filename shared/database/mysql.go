package database

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/fuadajip/kafka-cluster-go/shared/config"
	"github.com/fuadajip/kafka-cluster-go/shared/constant"
	"github.com/fuadajip/kafka-cluster-go/shared/log"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql" // register mysql with dd-trace
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"
)

var (
	once   sync.Once
	logger = log.NewServiceLog(constant.ServiceName)
)

type (

	// MysqlInterface is an interface that represent mysql methods in package database
	MysqlInterface interface {
		OpenMysqlConn(serviceName string) (*gorm.DB, error)
	}

	// database is a struct to map given struct
	database struct {
		SharedConfig config.ImmutableConfigInterface
	}
)

func (d *database) OpenMysqlConn(serviceName string) (*gorm.DB, error) {

	logger.Info("Start open mysql connection...")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		d.SharedConfig.GetDatabaseUser(),
		d.SharedConfig.GetDatabasePassword(),
		d.SharedConfig.GetDatabaseHost(),
		d.SharedConfig.GetDatabasePort(),
		d.SharedConfig.GetDatabaseName(),
	)

	// Register augments the provided driver with tracing, enabling it to be loaded by gormtrace.Open.
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName(serviceName))

	db, err := gormtrace.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(16)
	db.DB().SetMaxOpenConns(32)
	db.DB().SetConnMaxLifetime(1 * time.Hour)

	appEnv, exists := os.LookupEnv("APP_ENV")
	if exists {
		if appEnv == "staging" {
			logger.Info("Sql log mode enabled...")
			db.LogMode(true)
		} else if appEnv == "production" {
			logger.Info("Sql log mode disabled...")
			db.LogMode(false)
		} else if appEnv == "development" {
			logger.Info("Sql log mode enabled...")
			db.LogMode(true)
		}
	} else {
		logger.Info("Sql log mode enabled...")
		db.LogMode(true)
	}

	return db, nil
}

// NewMysql is an factory that implement of mysql database configuration
func NewMysql(config config.ImmutableConfigInterface) MysqlInterface {
	if config == nil {
		panic("[CONFIG] immutable config is required")
	}

	return &database{config}
}
