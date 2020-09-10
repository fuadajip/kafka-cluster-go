package log

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type (

	// ServiceLoggerInterface is an interface that represent methods in log package
	ServiceLoggerInterface interface {
		Info(args ...interface{})
		Infof(format string, args ...interface{})
		Error(args ...interface{})
		Errorf(format string, args ...interface{})
	}

	// ServiceLogger return object of interface in Service log package
	ServiceLogger struct {
		ServiceLoggerInterface
		prefix string
	}
)

var (
	logger   *logrus.Logger
	once     sync.Once
	instance *ServiceLogger
)

// Info is a logrus log message at level info on the standard logger
func (q *ServiceLogger) Info(args ...interface{}) {
	q.decorateLog().Info(args...)
}

// Infof is a logrus log message at level infof on the standard logger
func (q *ServiceLogger) Infof(format string, args ...interface{}) {
	q.decorateLog().Infof(format, args...)
}

// Error is a logrus log message at level error on the standard logger
func (q *ServiceLogger) Error(args ...interface{}) {
	q.decorateLog().Error(args...)
}

// Errorf is a logrus log message at level errorf on the standard logger
func (q *ServiceLogger) Errorf(format string, args ...interface{}) {
	q.decorateLog().Errorf(format, args...)
}

// NewServiceLog is a factory that return  interface of log pakcage
func NewServiceLog(prefix string) *ServiceLogger {
	once.Do(func() {
		logger = logrus.New()
		logger.Formatter = &prefixed.TextFormatter{
			FullTimestamp: true,
		}
		instance = &ServiceLogger{
			prefix: prefix,
		}
	})

	return instance
}

func (q *ServiceLogger) decorateLog() *logrus.Entry {
	var source string
	if pc, file, line, ok := runtime.Caller(2); ok {
		var funcName string
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
			if i := strings.LastIndex(funcName, "."); i != -1 {
				funcName = funcName[i+1:]
			}
		}

		source = fmt.Sprintf("%s:%v:%s()", path.Base(file), line, path.Base(funcName))
	}
	return logger.WithFields(logrus.Fields{
		"prefix": q.prefix,
		"source": source,
	})
}
