package healthcheck

import "github.com/labstack/echo"

type Usecase interface {
	DoHealthCheck(c echo.Context) (bool, error)
}

type Repository interface {
	MysqlHealthCheck() (bool, error)
	RedisHealthCheck() (bool, error)
	SqsHealthCheck(queueName string) (bool, error)
}
