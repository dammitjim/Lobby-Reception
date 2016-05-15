package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Opts configures various logging options
type Opts struct {
	ServiceName  string
	ServiceGroup string
	Level        string
}

// Fields is an alias for logrus.Fields to keep dem imports tidy like
type Fields map[string]interface{}

// Log is the importable logger
var logger *logrus.Logger
var serviceName string
var serviceGroup string
var host string

// Initialise sets the logger
func Initialise(opts Opts) {
	serviceName = opts.ServiceName
	serviceGroup = opts.ServiceGroup
	host, _ = os.Hostname()

	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{}

	switch opts.Level {
	case "debug":
		l.Level = logrus.DebugLevel
	case "info":
		l.Level = logrus.InfoLevel
	case "warn":
		l.Level = logrus.WarnLevel
	default:
		l.Level = logrus.ErrorLevel
	}

	logger = l
}

// Logger returns the logger object for use
func Logger() *logrus.Logger {
	return logger
}

// Log wraps logrus logging with the service fields
func Log() *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"service": serviceName,
		"group":   serviceGroup,
		"host":    host,
		"pid":     os.Getpid(),
	})
}

// WithFields is an alias for logrus.WithFields to accept our own
func WithFields(f Fields) *logrus.Entry {
	lf := logrus.Fields(f)
	return Log().WithFields(lf)
}
