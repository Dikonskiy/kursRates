// internal/logerr/logger.go
package logerr

import (
	"log"
	"log/slog"
	"os"

	"kursRates/internal/models"
)

func InitLogger() *slog.Logger {
	infoLogFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}
	if models.Config.IsProd {
		infoLog := slog.New(slog.NewJSONHandler(infoLogFile, nil))
		return infoLog
	} else {
		infoLog := slog.New(slog.NewTextHandler(infoLogFile, nil))
		return infoLog
	}
}
