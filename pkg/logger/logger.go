package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger
func Init(level string) {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func Debug(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Debug(msg)
}

func Info(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Info(msg)
}

func Warn(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Warn(msg)
}

func Error(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Error(msg)
}

// parseFields converts variadic arguments to logrus.Fields
func parseFields(fields ...interface{}) logrus.Fields {
	logFields := logrus.Fields{}
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fields[i].(string)
			value := fields[i+1]
			logFields[key] = value
		}
	}
	return logFields
}
