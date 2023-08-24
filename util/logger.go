package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Info  = logrus.New()
	Error = logrus.New()
)

func InitLogger() {
	infoFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Info.SetOutput(infoFile)
	} else {
		logrus.Warn("Failed to open info.log file. Using default stderr.")
		Info.SetOutput(os.Stderr)
	}

	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Error.SetOutput(errorFile)
	} else {
		logrus.Warn("Failed to open error.log file. Using default stderr.")
		Error.SetOutput(os.Stderr)
	}

	Info.SetFormatter(&logrus.TextFormatter{})
	Error.SetFormatter(&logrus.TextFormatter{})
}
