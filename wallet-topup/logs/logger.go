package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

func NewLogger() *Logger {
	log := logrus.New()

	file, err := os.OpenFile("wallet-log/wallet.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("⚠️ Failed to log to file, using default stderr")
	}

	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return log
}
