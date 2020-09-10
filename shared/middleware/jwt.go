package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fuadajip/kafka-cluster-go/models"
	"github.com/fuadajip/kafka-cluster-go/shared/constant"
	"github.com/fuadajip/kafka-cluster-go/shared/log"
	"github.com/fuadajip/kafka-cluster-go/shared/util"
	"github.com/labstack/echo"
)

var (
	logger = log.NewLog(constant.ServiceName)
)

// JWTAuthentication represent jwt middleware for restricted resources
func JWTAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac := c.(*util.CustomApplicationContext)
		conf := ac.SharedConf

		tokenHeader := c.Request().Header.Get("Authorization")

		// if token is empty return missing authentication token response
		if tokenHeader == "" {
			return ac.CustomResponse("failed", nil, "Missing Authenticastion Token", "[FAILED][KAFKA-CLUSTER-GO][MIDDLEWARE][JWTAuthenctication] Missing Authenticastion Token", http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// Splite header token normally comes with format Bearer xxx
		splittedTokenHeader := strings.Split(tokenHeader, " ")
		if len(splittedTokenHeader) < 2 {
			return ac.CustomResponse("failed", nil, "Malformed authentication token need Bearer token", "[FAILED][KAFKA-CLUSTER-GO][MIDDLEWARE][JWTAuthenctication] Malformed authentication token need Bearer token", http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// get token part
		tokenPart := splittedTokenHeader[1]

		// Initialize a new instance of `tk`
		tk := &models.UserJWT{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.GetTokenSecret()), nil
		})

		if err != nil {
			return ac.CustomResponse("failed", nil, "Malformed authentication", fmt.Sprintf("[FAILED][KAFKA-CLUSTER-GO][MIDDLEWARE][JWTAuthenctication] Malformed authentication err: %s", err.Error()), http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// check if it's expired
		if time.Now().Unix() > tk.ExpiresAt {
			return ac.CustomResponse("failed", nil, "Session ended", fmt.Sprintf("[FAILED][KAFKA-CLUSTER-GO][MIDDLEWARE][JWTAuthenctication] session ended err: %s", err.Error()), http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// check wheater token is valid or not (not signed in this server/ expired etc)
		if !token.Valid {
			return ac.CustomResponse("failed", nil, "Invalid token", fmt.Sprintf("[FAILED][KAFKA-CLUSTER-GO][MIDDLEWARE][JWTAuthenctication] invalid token err: %s", err.Error()), http.StatusUnauthorized, &models.ResponsePatternMeta{})
		}

		// map the parsed token into models
		ac.UserJWT = tk

		return next(c)

	}
}
