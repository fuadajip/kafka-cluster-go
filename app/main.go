package main

import (
	"fmt"

	"github.com/fuadajip/kafka-cluster-go/shared/constant"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo"

	"gopkg.in/go-playground/validator.v9"

	Config "github.com/fuadajip/kafka-cluster-go/shared/config"
	Container "github.com/fuadajip/kafka-cluster-go/shared/container"
	Database "github.com/fuadajip/kafka-cluster-go/shared/database"
	Logger "github.com/fuadajip/kafka-cluster-go/shared/log"
	Transporter "github.com/fuadajip/kafka-cluster-go/shared/transporter"
	CommonUtil "github.com/fuadajip/kafka-cluster-go/shared/util"
	Util "github.com/fuadajip/kafka-cluster-go/shared/util"

	//# --- domain import
	healthCheckHandler "github.com/fuadajip/kafka-cluster-go/domain/healthcheck/delivery/http"
	healthCheckRepository "github.com/fuadajip/kafka-cluster-go/domain/healthcheck/repository"
	healthCheckUsecase "github.com/fuadajip/kafka-cluster-go/domain/healthcheck/usecase"
	//# --- end domain import
)

var (
	logger = Logger.NewServiceLog(constant.ServiceName)
)

func main() {
	e := echo.New()
	container := Container.DefaultContainer()

	conf := container.MustGet("shared.config").(Config.ImmutableConfigInterface)
	mysql := container.MustGet("shared.mysql").(Database.MysqlInterface)

	mysqlSess, err := mysql.OpenMysqlConn(constant.ServiceName)
	if err != nil {
		msgError := fmt.Sprintf("Failed to open mysql connection: %s", err.Error())
		logger.Errorf(msgError)
	}

	redis := container.MustGet("shared.redis").(Database.RedisInterface)
	redisSess, err := redis.OpenRedisConn()
	if err != nil {
		msgError := fmt.Sprintf("Failed to open redis connection: %s", err.Error())
		logger.Errorf(msgError)
	}

	sqs := container.MustGet("shared.sqs").(Transporter.SQSInterface)
	sqsSess, err := sqs.CreateSQSService()
	if err != nil {
		msgError := fmt.Sprintf("Failed to create sqs service: %s", err.Error())
		logger.Errorf(msgError)
	}

	// rabbit := container.MustGet("shared.rabbit").(Transporter.RabbitmqInterface)
	// rabbitSess, err := rabbit.OpenRabbitConn()
	// if err != nil {
	// 	msgError := fmt.Sprintf("Failed to open rabbitmq connection: %s", err.Error())
	// 	logger.Errorf(msgError)
	// 	panic(msgError)
	// }

	// nats := container.MustGet("shared.nats").(Transporter.NatsInterface)
	// natsSess, err := nats.OpenNatsConn()
	// if err != nil {
	// 	msgError := fmt.Sprintf("Failed to open nats connection: %s", err.Error())
	// 	logger.Errorf(msgError)
	// }

	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	// provides protection against cross-site scripting (XSS) attack, content type sniffing,
	// clickjacking, insecure connection and other code injection attacks.
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())

	e.Use(echotrace.Middleware(echotrace.WithServiceName(constant.ServiceName)))

	e.Validator = &CommonUtil.CustomValidator{Validator: validator.New()}
	e.HTTPErrorHandler = func(err error, e echo.Context) {
		CommonUtil.CustomHTTPErrorHandler(err, e)
	}

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ac := &Util.CustomApplicationContext{
				Context:      c,
				Container:    container,
				SharedConf:   conf,
				MysqlSession: mysqlSess,
				RedisSession: redisSess,
				SqsService:   sqsSess,
				// NatsSession:  natsSess,
				// RabbitSession: rabbitSess,
			}

			return h(ac)
		}
	})

	// assign global context
	_ = &Util.CustomApplicationContext{
		SharedConf:   conf,
		RedisSession: redisSess,
		MysqlSession: mysqlSess,
		SqsService:   sqsSess,
		// RabbitSession: rabbitSess,
		// NatsSession: natsSess,
	}

	//# --- domain dependency injection
	healthCheckRepo := healthCheckRepository.NewHealthCheckRepository(redisSess, mysqlSess, sqsSess)
	healthCheckUcase := healthCheckUsecase.NewHealthCheckUsecase(healthCheckRepo)

	//# --- end domain dependency injection

	//# --- domain delivery handler injection
	healthCheckHandler.AddHealthCheckHandler(e, healthCheckUcase)

	//# --- end domain delivery handler injection

	e.Logger.Info(e.Start(fmt.Sprintf(":%d", conf.GetPort())))
}
