package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type AppLogger struct {
	log *logrus.Logger
}

func NewLogger() *AppLogger {
	log := logrus.New()

	file, err := os.OpenFile("wallet-log/wallet.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("⚠️ Failed to log to file, using default stderr")
	}

	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	return &AppLogger{log: log}
}

func (l *AppLogger) Info(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *AppLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *AppLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *AppLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *AppLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
