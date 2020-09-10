package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/fuadajip/kafka-cluster-go/domain/healthcheck"
	"github.com/fuadajip/kafka-cluster-go/shared/constant"
	"github.com/fuadajip/kafka-cluster-go/shared/log"
	"github.com/fuadajip/kafka-cluster-go/shared/util"
)

var (
	logger = log.NewServiceLog(constant.ServiceName)
)

type handlerHealtCheck struct {
	usecase healthcheck.Usecase
}

// AddHealthCheckHandler returns http handler for db session healthcheck
func AddHealthCheckHandler(e *echo.Echo, usecase healthcheck.Usecase) {
	handler := handlerHealtCheck{
		usecase: usecase,
	}

	e.GET("/api/healthz", handler.DoHeathCheck)
}

func (h handlerHealtCheck) DoHeathCheck(c echo.Context) error {
	ac := c.(*util.CustomApplicationContext)

	_, err := h.usecase.DoHealthCheck(c)
	if err != nil {
		msgError := fmt.Sprintf("Internal server error, err : %s", err.Error())
		logger.Error(msgError)
		return ac.CustomResponse("failed", nil, "[FAILED][HealthCheck][DoHeathCheck] failed system unhealthy", msgError, http.StatusInternalServerError, nil)
	}

	return ac.CustomResponse("success", nil, "[SUCCESS][HealthCheck][DoHeathCheck] success system healthy", "System healthy", http.StatusOK, nil)
}
