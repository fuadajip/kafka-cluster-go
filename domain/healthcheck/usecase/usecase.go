package usecase

import (
	"github.com/fuadajip/kafka-cluster-go/domain/healthcheck"
	"github.com/labstack/echo"
)

type usecase struct {
	repository healthcheck.Repository
}

// NewHealthCheckUsecase is a factory that return implementation of methods in healthcheck.Usecase interface
func NewHealthCheckUsecase(repository healthcheck.Repository) healthcheck.Usecase {
	return &usecase{
		repository: repository,
	}
}

func (u usecase) DoHealthCheck(c echo.Context) (bool, error) {
	_, err := u.repository.MysqlHealthCheck()
	if err != nil {
		return false, err
	}

	_, err = u.repository.RedisHealthCheck()
	if err != nil {
		return false, err
	}

	return true, nil
}
