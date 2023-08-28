package logerr

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Info  = logrus.New()
	Error = logrus.New()
)

type CustomFormatter struct {
	TimestampFormat string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := entry.Level.String()
	msg := entry.Message

	var logLine string
	if entry.Level == logrus.ErrorLevel {
		logLine = timestamp + " [" + level + "] " + msg + "\n"
	} else {
		logLine = timestamp + " [" + level + "] " + msg + "\n"
	}

	return []byte(logLine), nil
}

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

	customFormatter := &CustomFormatter{
		TimestampFormat: "Jan 02 15:04:05.000",
	}
	Info.SetFormatter(customFormatter)
	Error.SetFormatter(customFormatter)
}
