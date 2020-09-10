package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fuadajip/kafka-cluster-go/models"
	"github.com/fuadajip/kafka-cluster-go/shared/util"
	"github.com/labstack/echo"
)

// SecretStaticAuthentication represent static secret key middleware for internal communication authentication
func SecretStaticAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac := c.(*util.CustomApplicationContext)
		conf := ac.SharedConf

		tokenHeader := c.Request().Header.Get("Authorization")

		// if token is empty return missing authentication token response
		if tokenHeader == "" {
			return ac.CustomResponse("failed", nil, "Missing Authenticastion Token", "[FAILED][QPLUS-WHITELABEL-QUOTATION][MIDDLEWARE][JWTAuthenctication] Missing Authenticastion Token", http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// Splite header token normally comes with format Bearer xxx
		splittedTokenHeader := strings.Split(tokenHeader, " ")
		if len(splittedTokenHeader) < 2 {
			return ac.CustomResponse("failed", nil, "Malformed authentication token need Bearer token", "[FAILED][QPLUS-WHITELABEL-QUOTATION][MIDDLEWARE][JWTAuthenctication] Malformed authentication token need Bearer token", http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// get token part
		tokenPart := splittedTokenHeader[1]

		if tokenPart != conf.GetTokenSecretStatic() {
			return ac.CustomResponse("failed", nil, "Invalid internal token", fmt.Sprintf("[FAILED][QPLUS-WHITELABEL-QUOTATION][MIDDLEWARE][JWTAuthenctication] invalid internal token"), http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		return next(c)
	}
}
