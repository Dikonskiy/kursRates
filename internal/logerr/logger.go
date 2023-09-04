// internal/logerr/logger.go
package logerr

import (
	"log"
	"log/slog"
	"os"
)

func InitLogger() slog.Logger {
	infoLogFile, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}
	infoLog := slog.New(slog.NewJSONHandler(infoLogFile, nil))
	return *infoLog
}
