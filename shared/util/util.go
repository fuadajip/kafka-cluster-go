package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fgrosse/goldi"
	"github.com/fuadajip/kafka-cluster-go/models"
	"github.com/fuadajip/kafka-cluster-go/shared/config"
	"github.com/fuadajip/kafka-cluster-go/shared/constant"
	"github.com/fuadajip/kafka-cluster-go/shared/log"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/nats-io/nats.go"
	"github.com/streadway/amqp"
	"gopkg.in/go-playground/validator.v9"
)

type (
	serviceResponse models.ResponsePattern

	// CustomApplicationContext return service custom application context
	CustomApplicationContext struct {
		echo.Context
		Container     *goldi.Container
		SharedConf    config.ImmutableConfigInterface
		RedisSession  *redis.Client
		MysqlSession  *gorm.DB
		RabbitSession *amqp.Connection
		NatsSession   *nats.EncodedConn
		SqsService    *sqs.SQS
		DynamoService *dynamodb.DynamoDB
		UserJWT       *models.UserJWT
	}

	// CustomValidator return  custom validator
	CustomValidator struct {
		Validator *validator.Validate
	}
)

var (
	logger = log.NewServiceLog(constant.ServiceName)
)

// CustomResponse is a method that returns custom object response
func (c *CustomApplicationContext) CustomResponse(status string, data interface{}, message string, systemMessage string, code int, meta *models.ResponsePatternMeta) error {
	return c.JSON(code, &serviceResponse{
		Status:        status,
		Data:          data,
		Message:       message,
		SystemMessage: systemMessage,
		Code:          code,
		Meta:          meta,
	})
}

// Validate will validate given input with related struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func randomStringEngine(letter []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// RandomString will return random string
func RandomString(n int, kind string) string {
	switch kind {
	case "UPPERCASE":
		var letter = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := randomStringEngine(letter, n)
		return b
	case "LOWERCASE":
		var letter = []rune("abcdefghijklmnopqrstuvwxyz")

		b := randomStringEngine(letter, n)
		return b
	case "UPPERCASE_ALPHANUMERIC":
		var letter = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

		b := randomStringEngine(letter, n)
		return b
	case "LOWERCASE_ALPHANUMERIC":
		var letter = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		b := randomStringEngine(letter, n)
		return b
	default:
		var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		b := randomStringEngine(letter, n)
		return b
	}
}

// CustomHTTPErrorHandler will return custom echo http error handler
func CustomHTTPErrorHandler(err error, e echo.Context) {

	report, ok := err.(*echo.HTTPError)
	var msgError string

	if !ok {
		msgError = "[Generic] Internal server error, error [" + err.Error() + "]"
		report = echo.NewHTTPError(http.StatusInternalServerError, msgError)
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		msgError = "[Validation] Invalid validation, error [ field: %s is %s ]"
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				msgError = fmt.Sprintf(msgError, err.Field(), "is required")
				report = echo.NewHTTPError(http.StatusBadRequest, msgError)
			case "email":
				msgError = fmt.Sprintf(msgError, err.Field(), "is not valid email")
				report = echo.NewHTTPError(http.StatusBadRequest, msgError)
				break
			}
		}

	}

	logger.Error(msgError)
	qr := &serviceResponse{
		Code:    report.Code,
		Data:    nil,
		Message: fmt.Sprintf("%+v", report.Message),
		Meta:    nil,
		Status:  "failed",
	}

	e.JSON(report.Code, qr)
}

// CustomGormPaginationQuery will return method chaining of gorm fetch pagination
func CustomGormPaginationQuery(trx *gorm.DB, limit int, page int, orderBy string, order string) (*gorm.DB, error) {
	pageOffset := limit * (page - 1)

	if limit != 0 || page != 0 {
		trx = trx.Limit(limit).Offset(pageOffset)
	}
	if orderBy != "" && order != "" {
		trx = trx.Order(fmt.Sprintf("%s %s", orderBy, order))
	}

	return trx, nil
}

// PaginationCounter return response meta counter data
func PaginationCounter(query *models.CustomGormPaginationQuery, rows int) (resp *models.ResponsePatternMeta) {

	meta := &models.ResponsePatternMeta{}

	totalPages := math.Ceil(float64(rows) / float64(query.Limit))
	totalPagesInt := int(totalPages)
	meta.Page = &query.Page
	meta.Limit = &query.Limit
	meta.Count = &rows
	meta.Total = &totalPagesInt
	return meta
}

// DebugPrintStruct print struct to console log
func DebugPrintStruct(input ...interface{}) {
	fmt.Println("[DEBUG] =//=//=//")
	result, _ := json.Marshal(input)
	fmt.Println(string(result))
	fmt.Println("[DEBUG END] =//=//=//")
}

// DebugWritePDF write PDF from Base64 string to file
func DebugWritePDF(input *string) {
	dec, err := base64.StdEncoding.DecodeString(*input)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("debug.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

// DebugWriteString write raw string data to file
func DebugWriteString(input *string) {
	output := []byte(*input)
	err := ioutil.WriteFile("debug", output, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("[DEBUG PRINT SUCCESS]")
}

// GetAWSSession return the AWS session with static credentials or role check
func GetAWSSession(config config.ImmutableConfigInterface) (*session.Session, error) {
	var sess *session.Session
	var err error
	if config.GetAWSAccessKey() != "" && config.GetAWSSecretKey() != "" {
		logger.Info("Using static AWS Key: ", config.GetAWSAccessKey())
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(config.GetAWSRegion()),
			Credentials: credentials.NewStaticCredentials(config.GetAWSAccessKey(), config.GetAWSSecretKey(), ""),
		})
	} else {
		logger.Info("Cannot find AWS static credentials. Will init sess with role permission")
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(config.GetAWSRegion()),
		})
	}
	return sess, err
}
